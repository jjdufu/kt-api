package service

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/stretchr/testify/require"
)

func testTaskCallbackSetting() *system_setting.TaskCallbackSetting {
	return &system_setting.TaskCallbackSetting{
		Enabled:            true,
		AllowedPorts:       []string{"80", "443", "8080", "8443"},
		MaxURLLength:       2048,
		RejectUnresolvable: true,
	}
}

func publicResolver() staticSSRFResolver {
	return staticSSRFResolver{
		"hooks.example.com":  {{IP: net.ParseIP("93.184.216.34")}},
		"api.customer.io":    {{IP: net.ParseIP("2606:4700::6810:85e5")}},
		"internal.evil.com":  {{IP: net.ParseIP("10.0.0.8")}},
		"mixed.evil.com":     {{IP: net.ParseIP("93.184.216.34")}, {IP: net.ParseIP("127.0.0.1")}},
		"metadata.evil.com":  {{IP: net.ParseIP("169.254.169.254")}},
		"localhost":          {{IP: net.ParseIP("127.0.0.1")}},
		"a.allowed.example":  {{IP: net.ParseIP("93.184.216.35")}},
		"other.example.org":  {{IP: net.ParseIP("93.184.216.36")}},
		"rebind.example.com": {{IP: net.ParseIP("93.184.216.37")}},
	}
}

func TestValidateTaskCallbackURLAcceptsPublicTargets(t *testing.T) {
	setting := testTaskCallbackSetting()
	for _, raw := range []string{
		"https://hooks.example.com/kt/callback?token=abc",
		"http://hooks.example.com/cb",
		"https://hooks.example.com:8443/cb",
		"https://api.customer.io/v1/hook#frag",
		"https://93.184.216.34/cb",
		"HTTPS://Hooks.Example.com./cb",
	} {
		target, err := validateTaskCallbackURLWithResolver(context.Background(), raw, setting, publicResolver())
		require.NoError(t, err, raw)
		require.Equal(t, raw, target.URL)
		require.Len(t, target.URLHash, 64)
		require.NotContains(t, target.LogRef(), "/cb")
		require.NotContains(t, target.LogRef(), "token")
	}
}

func TestValidateTaskCallbackURLRejectsUnsafeTargets(t *testing.T) {
	setting := testTaskCallbackSetting()
	cases := map[string]string{
		"ftp://hooks.example.com/cb":                             "http or https",
		"file:///etc/passwd":                                     "http or https",
		"gopher://hooks.example.com/":                            "http or https",
		"javascript:alert(1)":                                    "http or https",
		"https:hooks.example.com":                                "absolute URL",
		"https://user:pass@hooks.example.com/":                   "credentials",
		"https://hooks.example.com/ a":                           "whitespace",
		"https://hooks.example.com/\x00":                         "control",
		"https://hooks.example.com:22/":                          "not allowed",
		"https://hooks.example.com:99999/":                       "invalid port",
		"http://127.0.0.1/":                                      "not allowed",
		"http://0.0.0.0/":                                        "not allowed",
		"http://10.1.2.3/":                                       "not allowed",
		"http://172.16.0.1/":                                     "not allowed",
		"http://192.168.50.51/":                                  "not allowed",
		"http://100.64.0.1/":                                     "not allowed",
		"http://100.100.100.200/":                                "not allowed",
		"http://169.254.169.254/latest/":                         "not allowed",
		"http://[::1]/":                                          "not allowed",
		"http://[fe80::1]/":                                      "not allowed",
		"http://[fe80::1%25eth0]/":                               "zone",
		"http://[fd00:ec2::254]/":                                "not allowed",
		"http://[::ffff:127.0.0.1]/":                             "not allowed",
		"http://0x7f000001/":                                     "canonical",
		"http://2130706433/":                                     "canonical",
		"http://127.1/":                                          "canonical",
		"http://0177.0.0.1/":                                     "canonical",
		"http://localhost/":                                      "disallowed address",
		"http://internal.evil.com/":                              "disallowed address",
		"http://mixed.evil.com/":                                 "disallowed address",
		"http://metadata.evil.com/":                              "disallowed address",
		"http://unknown.example.net/":                            "could not be resolved",
		"https://hooks.example.com/" + strings.Repeat("a", 2100): "exceed",
	}
	for raw, want := range cases {
		_, err := validateTaskCallbackURLWithResolver(context.Background(), raw, setting, publicResolver())
		require.Error(t, err, raw)
		require.Contains(t, err.Error(), want, raw)
		require.True(t, strings.HasPrefix(err.Error(), "callback_url"), raw)
	}
}

func TestValidateTaskCallbackURLSwitches(t *testing.T) {
	ctx := context.Background()

	httpsOnly := testTaskCallbackSetting()
	httpsOnly.HTTPSOnly = true
	_, err := validateTaskCallbackURLWithResolver(ctx, "http://hooks.example.com/", httpsOnly, publicResolver())
	require.ErrorContains(t, err, "must use https")
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://hooks.example.com/", httpsOnly, publicResolver())
	require.NoError(t, err)

	allowPrivate := testTaskCallbackSetting()
	allowPrivate.AllowPrivateIp = true
	_, err = validateTaskCallbackURLWithResolver(ctx, "http://192.168.50.51:8080/cb", allowPrivate, publicResolver())
	require.NoError(t, err)

	allowlist := testTaskCallbackSetting()
	allowlist.DomainAllowlistMode = true
	allowlist.DomainList = []string{"*.allowed.example"}
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://a.allowed.example/", allowlist, publicResolver())
	require.NoError(t, err)
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://other.example.org/", allowlist, publicResolver())
	require.ErrorContains(t, err, "not in whitelist")

	blocklist := testTaskCallbackSetting()
	blocklist.DomainList = []string{"other.example.org"}
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://other.example.org/", blocklist, publicResolver())
	require.ErrorContains(t, err, "blacklist")

	lenient := testTaskCallbackSetting()
	lenient.RejectUnresolvable = false
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://unknown.example.net/", lenient, publicResolver())
	require.NoError(t, err)

	ports := testTaskCallbackSetting()
	ports.AllowedPorts = []string{"9000-9001"}
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://hooks.example.com:9001/", ports, publicResolver())
	require.NoError(t, err)
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://hooks.example.com/", ports, publicResolver())
	require.ErrorContains(t, err, "port 443 is not allowed")

	short := testTaskCallbackSetting()
	short.MaxURLLength = 64
	_, err = validateTaskCallbackURLWithResolver(ctx, "https://hooks.example.com/"+strings.Repeat("a", 64), short, publicResolver())
	require.ErrorContains(t, err, "64 bytes")
}

// sequenceResolver returns a different answer per lookup to simulate DNS
// rebinding between submit-time validation and delivery.
type sequenceResolver struct {
	answers [][]net.IPAddr
	calls   atomic.Int32
}

func (r *sequenceResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	i := int(r.calls.Add(1)) - 1
	if i >= len(r.answers) {
		i = len(r.answers) - 1
	}
	return r.answers[i], nil
}

func TestTaskCallbackClientBlocksDNSRebinding(t *testing.T) {
	setting := testTaskCallbackSetting()
	resolver := &sequenceResolver{answers: [][]net.IPAddr{
		{{IP: net.ParseIP("93.184.216.37")}},
		{{IP: net.ParseIP("127.0.0.1")}},
	}}
	_, err := validateTaskCallbackURLWithResolver(context.Background(), "https://rebind.example.com/cb", setting, resolver)
	require.NoError(t, err)

	dialed := false
	client := newTaskCallbackHTTPClient(resolver, func(ctx context.Context, network, address string) (net.Conn, error) {
		dialed = true
		return nil, &net.OpError{Op: "dial"}
	}, func() *system_setting.TaskCallbackSetting { return setting })
	_, err = client.Post("https://rebind.example.com/cb", "application/json", strings.NewReader("{}"))
	require.Error(t, err)
	require.True(t, IsTaskCallbackBlocked(err), err.Error())
	require.False(t, dialed)
}

func TestTaskCallbackClientRejectsAnyPrivateAnswer(t *testing.T) {
	setting := testTaskCallbackSetting()
	client := newTaskCallbackHTTPClient(publicResolver(), func(ctx context.Context, network, address string) (net.Conn, error) {
		t.Fatalf("must not dial %s", address)
		return nil, nil
	}, func() *system_setting.TaskCallbackSetting { return setting })
	_, err := client.Post("http://mixed.evil.com/cb", "application/json", strings.NewReader("{}"))
	require.True(t, IsTaskCallbackBlocked(err))

	setting.HTTPSOnly = true
	_, err = client.Post("http://hooks.example.com/cb", "application/json", strings.NewReader("{}"))
	require.True(t, IsTaskCallbackBlocked(err))
}

func TestTaskCallbackClientDialsValidatedIPAndDoesNotFollowRedirects(t *testing.T) {
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		http.Redirect(w, r, "http://169.254.169.254/latest/meta-data", http.StatusFound)
	}))
	defer server.Close()
	_, port, err := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	require.NoError(t, err)

	setting := testTaskCallbackSetting()
	setting.AllowedPorts = []string{port}
	var dialedAddr string
	client := newTaskCallbackHTTPClient(publicResolver(), func(ctx context.Context, network, address string) (net.Conn, error) {
		dialedAddr = address
		// The resolved public IP is validated; the test then routes the
		// connection to the local server.
		return (&net.Dialer{}).DialContext(ctx, network, strings.TrimPrefix(server.URL, "http://"))
	}, func() *system_setting.TaskCallbackSetting { return setting })
	resp, err := client.Post("http://hooks.example.com:"+port+"/cb", "application/json", strings.NewReader("{}"))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusFound, resp.StatusCode)
	require.Equal(t, int32(1), hits.Load())
	require.Equal(t, net.JoinHostPort("93.184.216.34", port), dialedAddr)
}
