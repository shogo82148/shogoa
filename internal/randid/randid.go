package randid

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// New generates a random string of length n.
// The string contains only alphabets and numbers.
func New(n int) string {
	m := (n + 3) / 4 * 4
	buf := make([]byte, m + m/4*3)
	buf1 := buf[:m]
	buf2 := buf[m:]
	var builder strings.Builder
	builder.Grow(n)

LOOP:
	for {
		if _, err := rand.Read(buf2); err != nil {
			panic(err)
		}
		base64.StdEncoding.Encode(buf1, buf2)
		for _, b := range buf1 {
			if b == '+' || b == '/' {
				continue
			}
			builder.WriteByte(b)
			if builder.Len() == n {
				break LOOP
			}
		}
	}
	return builder.String()
}
