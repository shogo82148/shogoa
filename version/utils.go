package version

import (
	"fmt"
	"strconv"
	"strings"
)

// String returns the complete version number.
func String() string {
	return "v" + strconv.Itoa(Major) + "." + strconv.Itoa(Minor) + "." + strconv.Itoa(Build)
}

// Compatible returns true if Major matches the major version of the given version string.
// It returns an error if the given string is not a valid version string.
func Compatible(v string) (bool, error) {
	if len(v) < 6 {
		return false, fmt.Errorf("invalid version string format %#v", v)
	}
	if v[0] != 'v' {
		return false, fmt.Errorf("version string must start with 'v' %#v", v)
	}
	v = v[1:]
	elems := strings.Split(v, ".")
	if len(elems) != 3 {
		return false, fmt.Errorf("version not of the form Major.Minor.Build %#v", v)
	}
	mj, err := strconv.Atoi(elems[0])
	if err != nil {
		return false, fmt.Errorf("invalid major version number %#v, must be number", elems[0])
	}
	return mj == Major, nil
}
