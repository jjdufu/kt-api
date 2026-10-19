package system_setting

import (
	"strings"
	"time"

	"github.com/QuantumNous/new-api/setting/config"
)

// Upstream callback injection modes for the optional fast path.
const (
	TaskCallbackUpstreamInjectNever                    = "never"
	TaskCallbackUpstreamInjectOnlyWhenCustomerCallback = "only_when_customer_callback"
	TaskCallbackUpstreamInjectAlways                   = "always"
)

// TaskCallbackSetting configures delivery of kt-api's own task callbacks to
// customers (Seedance callback_url) and the optional upstream fast path.
// Stored as options "task_callback_setting.<json tag>".
type TaskCallbackSetting struct {
	// Enabled accepts callback_url on submit and runs the delivery worker.
	Enabled bool `json:"enabled"`
	// ExtraHeadersEnabled adds X-KtApi-Event-Id / -Delivery-Attempt / -Task-Id.
	ExtraHeadersEnabled bool `json:"extra_headers_enabled"`

	// Customer callback_url validation (SSRF protection).
	HTTPSOnly           bool     `json:"https_only"`
	AllowPrivateIp      bool     `json:"allow_private_ip"` // testing only
	DomainAllowlistMode bool     `json:"domain_allowlist_mode"`
	DomainList          []string `json:"domain_list"`   // example.com, *.example.com
	AllowedPorts        []string `json:"allowed_ports"` // 80, 443, 8000-9000
	MaxURLLength        int      `json:"max_url_length"`
	// RejectUnresolvable rejects a callback host that does not resolve at
	// submit time; when false the dial-time check alone decides.
	RejectUnresolvable bool `json:"reject_unresolvable"`

	// Delivery.
	TimeoutSeconds       int  `json:"timeout_seconds"`
	RetryIntervalSeconds int  `json:"retry_interval_seconds"`
	TerminalMaxAttempts  int  `json:"terminal_max_attempts"`
	ExtendedRetry        bool `json:"extended_retry"`
	Workers              int  `json:"workers"`
	PerHostConcurrency   int  `json:"per_host_concurrency"`
	RetentionDays        int  `json:"retention_days"`

	// Optional HMAC signing (kt-api extension, off by default).
	SigningEnabled bool   `json:"signing_enabled"`
	SigningSecret  string `json:"signing_secret"`

	// Optional upstream callback fast path (off by default).
	UpstreamPublicBase string `json:"upstream_public_base"`
	UpstreamInjectMode string `json:"upstream_inject_mode"`
}

var defaultTaskCallbackSetting = TaskCallbackSetting{
	Enabled:              true,
	ExtraHeadersEnabled:  true,
	HTTPSOnly:            false,
	AllowPrivateIp:       false,
	DomainAllowlistMode:  false,
	DomainList:           []string{},
	AllowedPorts:         []string{"80", "443", "8080", "8443"},
	MaxURLLength:         2048,
	RejectUnresolvable:   true,
	TimeoutSeconds:       5,
	RetryIntervalSeconds: 5,
	TerminalMaxAttempts:  4,
	ExtendedRetry:        false,
	Workers:              16,
	PerHostConcurrency:   4,
	RetentionDays:        7,
	SigningEnabled:       false,
	SigningSecret:        "",
	UpstreamPublicBase:   "",
	UpstreamInjectMode:   TaskCallbackUpstreamInjectOnlyWhenCustomerCallback,
}

func init() {
	config.GlobalConfig.Register("task_callback_setting", &defaultTaskCallbackSetting)
}

// GetTaskCallbackSetting returns the live task callback setting.
func GetTaskCallbackSetting() *TaskCallbackSetting {
	return &defaultTaskCallbackSetting
}

func clampInt(value, minimum, maximum, fallback int) int {
	if value <= 0 {
		return fallback
	}
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

// Timeout is the per-attempt delivery timeout (1-30s, default 5s).
func (s *TaskCallbackSetting) Timeout() time.Duration {
	return time.Duration(clampInt(s.TimeoutSeconds, 1, 30, 5)) * time.Second
}

// RetryInterval is the wait between documented retries (1-60s, default 5s).
func (s *TaskCallbackSetting) RetryInterval() time.Duration {
	return time.Duration(clampInt(s.RetryIntervalSeconds, 1, 60, 5)) * time.Second
}

// TerminalAttempts is the attempt budget of succeeded/failed events
// (1-10, default 4 = first attempt + 3 retries).
func (s *TaskCallbackSetting) TerminalAttempts() int {
	return clampInt(s.TerminalMaxAttempts, 1, 10, 4)
}

// WorkerCount bounds concurrent deliveries per instance (1-128, default 16).
func (s *TaskCallbackSetting) WorkerCount() int {
	return clampInt(s.Workers, 1, 128, 16)
}

// HostConcurrency bounds concurrent deliveries per target host (1-64, default 4).
func (s *TaskCallbackSetting) HostConcurrency() int {
	return clampInt(s.PerHostConcurrency, 1, 64, 4)
}

// Retention is how long finished events are kept (1-90 days, default 7).
func (s *TaskCallbackSetting) Retention() time.Duration {
	return time.Duration(clampInt(s.RetentionDays, 1, 90, 7)) * 24 * time.Hour
}

// URLLengthLimit bounds customer callback URLs (64-8192, default 2048).
func (s *TaskCallbackSetting) URLLengthLimit() int {
	return clampInt(s.MaxURLLength, 64, 8192, 2048)
}

// UpstreamCallbackBase returns the normalized public base URL for upstream
// callbacks, or "" when the fast path is disabled.
func (s *TaskCallbackSetting) UpstreamCallbackBase() string {
	if s.UpstreamInjectMode == TaskCallbackUpstreamInjectNever {
		return ""
	}
	return strings.TrimRight(strings.TrimSpace(s.UpstreamPublicBase), "/")
}

// InjectUpstreamCallback reports whether a task should carry kt-api's own
// upstream callback URL, given whether the customer supplied callback_url.
func (s *TaskCallbackSetting) InjectUpstreamCallback(customerCallback bool) bool {
	if s.UpstreamCallbackBase() == "" {
		return false
	}
	switch s.UpstreamInjectMode {
	case TaskCallbackUpstreamInjectAlways:
		return true
	case TaskCallbackUpstreamInjectNever:
		return false
	default:
		return customerCallback
	}
}
