package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/system_setting"
)

// ContextKeySeedanceCallbackTarget carries the validated customer
// callback_url (*TaskCallbackTarget) from the submit middleware to the
// submit handler.
const ContextKeySeedanceCallbackTarget = "seedance_callback_target"

// taskCallbackResolveTimeout bounds the submit-time DNS check of a customer
// callback_url host.
const taskCallbackResolveTimeout = 2 * time.Second

// TaskCallbackTarget is a validated customer callback_url.
type TaskCallbackTarget struct {
	URL     string
	Scheme  string
	Host    string
	URLHash string
}

// LogRef identifies the target in logs without exposing path or query, which
// may carry customer secrets: scheme://host#<sha256 prefix>.
func (t *TaskCallbackTarget) LogRef() string {
	if t == nil {
		return ""
	}
	return t.Scheme + "://" + t.Host + "#" + TaskCallbackURLHashPrefix(t.URLHash)
}

// TaskCallbackURLHash returns the hex SHA-256 of a callback URL.
func TaskCallbackURLHash(rawURL string) string {
	sum := sha256.Sum256([]byte(rawURL))
	return hex.EncodeToString(sum[:])
}

// TaskCallbackURLHashPrefix returns the 12-char prefix used in logs.
func TaskCallbackURLHashPrefix(hash string) string {
	if len(hash) > 12 {
		return hash[:12]
	}
	return hash
}

// TaskCallbackBlockedError marks a callback target rejected by SSRF policy.
// Deliveries failing with it are classified ssrf_blocked and not retried.
type TaskCallbackBlockedError struct {
	Reason string
}

func (e *TaskCallbackBlockedError) Error() string {
	return "callback target blocked: " + e.Reason
}

func taskCallbackBlocked(format string, args ...any) error {
	return &TaskCallbackBlockedError{Reason: fmt.Sprintf(format, args...)}
}

// IsTaskCallbackBlocked reports whether err is an SSRF policy rejection.
func IsTaskCallbackBlocked(err error) bool {
	var blocked *TaskCallbackBlockedError
	return errors.As(err, &blocked)
}

// taskCallbackProtection builds the SSRF policy from task_callback_setting.
// It is independent from fetch_setting: customer callbacks are always
// protected and domain hosts always have their resolved IPs checked.
func taskCallbackProtection(setting *system_setting.TaskCallbackSetting) (*common.SSRFProtection, error) {
	protection, err := common.NewSSRFProtectionFromFetchSetting(
		setting.AllowPrivateIp,
		setting.DomainAllowlistMode,
		false, // IP blacklist mode with an empty list: only the private-range check applies
		setting.DomainList,
		nil,
		setting.AllowedPorts,
		true,
	)
	if err != nil {
		return nil, err
	}
	return protection, nil
}

type parsedTaskCallbackURL struct {
	url    *url.URL
	scheme string
	host   string
	port   int
}

// parseTaskCallbackURL performs every check that does not need DNS.
func parseTaskCallbackURL(raw string, setting *system_setting.TaskCallbackSetting, protection *common.SSRFProtection) (*parsedTaskCallbackURL, error) {
	if len(raw) > setting.URLLengthLimit() {
		return nil, fmt.Errorf("callback_url must not exceed %d bytes", setting.URLLengthLimit())
	}
	for _, r := range raw {
		if r <= 0x20 || r == 0x7f {
			return nil, errors.New("callback_url must not contain whitespace or control characters")
		}
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, errors.New("callback_url is not a valid URL")
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "https":
	case "http":
		if setting.HTTPSOnly {
			return nil, errors.New("callback_url must use https")
		}
	default:
		return nil, errors.New("callback_url must use http or https")
	}
	if u.Opaque != "" || u.Host == "" {
		return nil, errors.New("callback_url must be an absolute URL with a host")
	}
	if u.User != nil {
		return nil, errors.New("callback_url must not contain user credentials")
	}
	host := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	if host == "" {
		return nil, errors.New("callback_url must be an absolute URL with a host")
	}
	if strings.Contains(host, "%") {
		return nil, errors.New("callback_url host must not contain an IPv6 zone")
	}
	port := 443
	if scheme == "http" {
		port = 80
	}
	if portText := u.Port(); portText != "" {
		port, err = strconv.Atoi(portText)
		if err != nil || port < 1 || port > 65535 {
			return nil, errors.New("callback_url has an invalid port")
		}
	}
	if ip := net.ParseIP(host); ip == nil && isNonCanonicalIPv4Host(host) {
		return nil, errors.New("callback_url host is not a canonical IP address or domain name")
	}
	if err := protection.ValidateNetworkTarget(host, port); err != nil {
		return nil, fmt.Errorf("callback_url is not allowed: %v", err)
	}
	return &parsedTaskCallbackURL{url: u, scheme: scheme, host: host, port: port}, nil
}

// isNonCanonicalIPv4Host rejects numeric host spellings such as 2130706433,
// 0x7f000001, 0177.0.0.1 or 127.1 that some resolvers turn into IPs. Real
// top-level domains are never numeric.
func isNonCanonicalIPv4Host(host string) bool {
	labels := strings.Split(host, ".")
	last := labels[len(labels)-1]
	if last == "" {
		return true
	}
	if strings.HasPrefix(last, "0x") {
		return true
	}
	for _, r := range last {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ValidateTaskCallbackURL validates a customer callback_url at submit time.
// Error messages are safe to return to the customer.
func ValidateTaskCallbackURL(ctx context.Context, raw string) (*TaskCallbackTarget, error) {
	return validateTaskCallbackURLWithResolver(ctx, raw, system_setting.GetTaskCallbackSetting(), net.DefaultResolver)
}

func validateTaskCallbackURLWithResolver(ctx context.Context, raw string, setting *system_setting.TaskCallbackSetting, resolver ssrfResolver) (*TaskCallbackTarget, error) {
	protection, err := taskCallbackProtection(setting)
	if err != nil {
		return nil, errors.New("callback_url validation is misconfigured")
	}
	parsed, err := parseTaskCallbackURL(raw, setting, protection)
	if err != nil {
		return nil, err
	}
	if net.ParseIP(parsed.host) == nil {
		lookupCtx, cancel := context.WithTimeout(ctx, taskCallbackResolveTimeout)
		defer cancel()
		addrs, lookupErr := resolver.LookupIPAddr(lookupCtx, parsed.host)
		if lookupErr != nil || len(addrs) == 0 {
			if setting.RejectUnresolvable {
				return nil, errors.New("callback_url host could not be resolved")
			}
		}
		for _, addr := range addrs {
			if err := protection.ValidateResolvedIP(parsed.host, addr.IP); err != nil {
				return nil, errors.New("callback_url resolves to a disallowed address")
			}
		}
	}
	return &TaskCallbackTarget{
		URL:     raw,
		Scheme:  parsed.scheme,
		Host:    parsed.host,
		URLHash: TaskCallbackURLHash(raw),
	}, nil
}

// taskCallbackDialer re-validates every resolved IP immediately before
// dialing and connects to the IP itself, defeating DNS rebinding.
type taskCallbackDialer struct {
	resolver      ssrfResolver
	dialContext   func(ctx context.Context, network, address string) (net.Conn, error)
	getProtection func() (*common.SSRFProtection, error)
}

func (d *taskCallbackDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	protection, err := d.getProtection()
	if err != nil {
		return nil, err
	}
	host, portText, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid dial address %s: %w", addr, err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %s", portText)
	}
	if err := protection.ValidateNetworkTarget(host, port); err != nil {
		return nil, taskCallbackBlocked("%v", err)
	}
	if ip := net.ParseIP(host); ip != nil {
		return d.dialContext(ctx, network, net.JoinHostPort(ip.String(), portText))
	}
	if isNonCanonicalIPv4Host(strings.ToLower(host)) {
		return nil, taskCallbackBlocked("non-canonical host %s", host)
	}
	resolved, err := d.resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, &net.DNSError{Err: err.Error(), Name: host}
	}
	var candidates []net.IP
	for _, ipAddr := range resolved {
		ip := ipAddr.IP
		if ip == nil || !networkAllowsIP(network, ip) {
			continue
		}
		if err := protection.ValidateResolvedIP(host, ip); err != nil {
			return nil, taskCallbackBlocked("%v", err)
		}
		candidates = append(candidates, ip)
	}
	var lastErr error
	for _, ip := range candidates {
		conn, err := d.dialContext(ctx, network, net.JoinHostPort(ip.String(), portText))
		if err == nil {
			return conn, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &net.DNSError{Err: "no usable IP addresses", Name: host}
}

// taskCallbackRoundTripper re-checks the URL policy (scheme, https-only,
// port, domain list) before each request so setting changes apply to queued
// events too.
type taskCallbackRoundTripper struct {
	base       http.RoundTripper
	getSetting func() *system_setting.TaskCallbackSetting
}

func (t *taskCallbackRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("invalid request")
	}
	setting := t.getSetting()
	protection, err := taskCallbackProtection(setting)
	if err != nil {
		return nil, err
	}
	if _, err := parseTaskCallbackURL(req.URL.String(), setting, protection); err != nil {
		return nil, taskCallbackBlocked("%v", err)
	}
	return t.base.RoundTrip(req)
}

func (t *taskCallbackRoundTripper) CloseIdleConnections() {
	if closer, ok := t.base.(interface{ CloseIdleConnections() }); ok {
		closer.CloseIdleConnections()
	}
}

// newTaskCallbackHTTPClient builds the delivery client: no proxy (a proxy
// would bypass the dial-time IP check), no redirects, normal TLS
// verification, dial-time SSRF validation. The per-attempt deadline is
// applied by the caller's context.
func newTaskCallbackHTTPClient(resolver ssrfResolver, dialContext func(ctx context.Context, network, address string) (net.Conn, error), getSetting func() *system_setting.TaskCallbackSetting) *http.Client {
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	if dialContext == nil {
		dialContext = (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext
	}
	if getSetting == nil {
		getSetting = system_setting.GetTaskCallbackSetting
	}
	dialer := &taskCallbackDialer{
		resolver:    resolver,
		dialContext: dialContext,
		getProtection: func() (*common.SSRFProtection, error) {
			return taskCallbackProtection(getSetting())
		},
	}
	transport := &http.Transport{
		Proxy:                  nil,
		DialContext:            dialer.DialContext,
		ForceAttemptHTTP2:      true,
		MaxIdleConns:           256,
		MaxIdleConnsPerHost:    8,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    5 * time.Second,
		ResponseHeaderTimeout:  30 * time.Second,
		MaxResponseHeaderBytes: 64 << 10,
	}
	return &http.Client{
		Transport: &taskCallbackRoundTripper{base: transport, getSetting: getSetting},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
