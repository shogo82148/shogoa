package shogoa

import "testing"

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		val    string
		valid  bool
	}{
		{
			name:   "invalid data format",
			format: FormatDate,
			val:    "201510-26",
			valid:  false,
		},
		{
			name:   "valid data format",
			format: FormatDate,
			val:    "2015-10-26",
			valid:  true,
		},
		{
			name:   "invalid datetime format",
			format: FormatDateTime,
			val:    "201510-26T08:31:23Z",
			valid:  false,
		},
		{
			name:   "valid datetime format",
			format: FormatDateTime,
			val:    "2015-10-26T08:31:23Z",
			valid:  true,
		},
		{
			name:   "invalid UUID format",
			format: FormatUUID,
			val:    "invalid-uuid",
			valid:  false,
		},
		{
			name:   "valid UUID format",
			format: FormatUUID,
			val:    "550e8400-e29b-41d4-a716-446655440000",
			valid:  true,
		},
		{
			name:   "invalid email format",
			format: FormatEmail,
			val:    "invalid-email",
			valid:  false,
		},
		{
			name:   "valid email format",
			format: FormatEmail,
			val:    "john.doe@example.com",
			valid:  true,
		},
		{
			name:   "invalid hostname format",
			format: FormatHostname,
			val:    "_hi_",
			valid:  false,
		},
		{
			name:   "valid hostname format",
			format: FormatHostname,
			val:    "example.com",
			valid:  true,
		},
		{
			name:   "invalid ipv4 format",
			format: FormatIPv4,
			val:    "invalid-ipv4",
			valid:  false,
		},
		{
			name:   "valid ipv6 format but we are checking for ipv4",
			format: FormatIPv4,
			val:    "2001:db8::1",
			valid:  false,
		},
		{
			name:   "valid ipv4 format",
			format: FormatIPv4,
			val:    "192.0.2.1",
			valid:  true,
		},
		{
			name:   "invalid ipv6 format",
			format: FormatIPv6,
			val:    "invalid-ipv6",
			valid:  false,
		},
		{
			name:   "valid ipv4 format but we are checking for ipv6",
			format: FormatIPv6,
			val:    "192.0.2.1",
			valid:  false,
		},
		{
			name:   "valid ipv6 format",
			format: FormatIPv6,
			val:    "2001:db8::1",
			valid:  true,
		},
		{
			name:   "invalid ip format",
			format: FormatIP,
			val:    "invalid-ip",
			valid:  false,
		},
		{
			name:   "valid ip format (ipv4)",
			format: FormatIP,
			val:    "192.0.2.1",
			valid:  true,
		},
		{
			name:   "valid ip format (ipv6)",
			format: FormatIP,
			val:    "2001:db8::1",
			valid:  true,
		},
		{
			name:   "invalid URI format",
			format: FormatURI,
			val:    "invalid-uri",
			valid:  false,
		},
		{
			name:   "valid URI format",
			format: FormatURI,
			val:    "https://example.com",
			valid:  true,
		},
		{
			name:   "invalid MAC format",
			format: FormatMAC,
			val:    "invalid-mac",
			valid:  false,
		},
		{
			name:   "valid MAC format",
			format: FormatMAC,
			val:    "00:00:5e:00:53:01",
			valid:  true,
		},
		{
			name:   "invalid CIDR format",
			format: FormatCIDR,
			val:    "invalid-cidr",
			valid:  false,
		},
		{
			name:   "valid CIDR format",
			format: FormatCIDR,
			val:    "192.0.2.0/24",
			valid:  true,
		},
		{
			name:   "invalid regexp format",
			format: FormatRegexp,
			val:    "[",
			valid:  false,
		},
		{
			name:   "valid regexp format",
			format: FormatRegexp,
			val:    "[a-z]+",
			valid:  true,
		},
		{
			name:   "invalid RFC1123 format",
			format: FormatRFC1123,
			val:    "invalid-rfc1123",
			valid:  false,
		},
		{
			name:   "valid RFC1123 format",
			format: FormatRFC1123,
			val:    "Mon, 02 Jan 2006 15:04:05 MST",
			valid:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFormat(tc.format, tc.val)
			if tc.valid && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Error("expected an error")
			}
		})
	}
}

func TestValidatePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		val     string
		valid   bool
	}{
		{
			name:    "invalid pattern",
			pattern: "^[a-z]+$",
			val:     "123",
			valid:   false,
		},
		{
			name:    "valid pattern",
			pattern: "^[a-z]+$",
			val:     "abc",
			valid:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid := ValidatePattern(tc.pattern, tc.val)
			if valid != tc.valid {
				t.Errorf("unexpected result: %t", valid)
			}
		})
	}
}
