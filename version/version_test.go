package version

import (
	"fmt"
	"regexp"
	"testing"
)

func TestString(t *testing.T) {
	re := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	ver := String()
	if !re.MatchString(ver) {
		t.Errorf("unexpected version string: %s", ver)
	}
}

func TestCompatible(t *testing.T) {
	tests := []struct {
		name     string
		ver      string
		expected bool
	}{
		{
			name:     "a version with identical major should be compatible",
			ver:      fmt.Sprintf("v%d.12.13", Major),
			expected: true,
		},
		{
			name:     "a version with different major should not be compatible",
			ver:      "v99999121299999.1.0",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			compatible, err := Compatible(tc.ver)
			if err != nil {
				t.Fatal(err)
			}
			if compatible != tc.expected {
				t.Errorf("unexpected result: %t", compatible)
			}
		})
	}
}

func TestCompatible_error(t *testing.T) {
	tests := []struct {
		name string
		ver  string
	}{
		{
			name: "version string too short",
			ver:  "v1.2",
		},
		{
			name: "version string must start with 'v'",
			ver:  "w1.2.3",
		},
		{
			name: "version not of the form Major.Minor.Build",
			ver:  "v99999121299999.2",
		},
		{
			name: "major version must be a number",
			ver:  "vX.2.3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.ver, func(t *testing.T) {
			_, err := Compatible(tc.ver)
			if err == nil {
				t.Errorf("%s: expected an error", tc.name)
			}
		})
	}
}
