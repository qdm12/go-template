package log

import (
	"net/http"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractClientIP(t *testing.T) {
	t.Parallel()

	makeHeader := func(keyValues map[string][]string) http.Header {
		header := http.Header{}
		for key, values := range keyValues {
			for _, value := range values {
				header.Add(key, value)
			}
		}
		return header
	}

	testCases := map[string]struct {
		r  *http.Request
		ip netip.Addr
	}{
		"nil request": {},
		"empty request": {
			r: &http.Request{},
		},
		"request with remote address": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
			},
			ip: netip.AddrFrom4([4]byte{99, 99, 99, 99}),
		},
		"request with xRealIP header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP": {"88.88.88.88"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with xRealIP header and public XForwardedFor IP": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP":       {"77.77.77.77"},
					"X-Forwarded-For": {"88.88.88.88"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with xRealIP header and private XForwardedFor IP": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP":       {"88.88.88.88"},
					"X-Forwarded-For": {"10.0.0.5"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with single public IP in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"88.88.88.88"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with two public IPs in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"88.88.88.88", "77.77.77.77"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with private and public IPs in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5", "88.88.88.88", "10.0.0.1", "77.77.77.77"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with single private IP in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{192, 168, 1, 5}),
		},
		"request with private IPs in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5", "10.0.0.17"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{192, 168, 1, 5}),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ip := extractClientIP(testCase.r)
			assert.Equal(t, testCase.ip, ip)
		})
	}
}

func Test_splitHostPort(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		address    string
		ip         string
		port       string
		errMessage string
	}{
		"empty_address": {
			errMessage: "missing port in address",
		},
		"invalid_address_with_brackets": {
			address:    "[abc]",
			errMessage: "address [abc]: missing port in address",
		},
		"address_with_brackets_without_port": {
			address:    "[::1]",
			errMessage: "address [::1]: missing port in address",
		},
		"address_with_brackets": {
			address: "[::1]:8000",
			ip:      "::1",
			port:    "8000",
		},
		"malformed_ipv6_address_port": {
			address:    "::x:",
			errMessage: "address ::x:: too many colons in address",
		},
		"ipv6_address": {
			address: "::1:8000",
			ip:      "::1",
			port:    "8000",
		},
		"ipv4_address": {
			address: "1.2.3.4:8000",
			ip:      "1.2.3.4",
			port:    "8000",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ip, port, err := splitHostPort(testCase.address)
			assert.Equal(t, testCase.ip, ip)
			assert.Equal(t, testCase.port, port)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
