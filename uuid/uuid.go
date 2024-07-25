// The uuid package is minimum implement of UUID(Universally Unique Identifier).
package uuid

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// UUID is a type that represents a UUID.
type UUID [16]byte

// String parse helpers.
var (
	urnPrefix      = "urn:uuid:"
	urnPrefixBytes = []byte(urnPrefix)
	bytePos        = [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	}
)

// FromString returns UUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(text string) (UUID, error) {
	switch len(text) {
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36:

	// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + len(urnPrefixBytes):
		prefix := text[:len(urnPrefixBytes)]
		if !strings.EqualFold(prefix, urnPrefix) {
			return UUID{}, fmt.Errorf("uuid: invalid UUID string: %q", text)
		}
		text = text[len(urnPrefixBytes):]

	// {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	case 36 + 2:
		if text[0] != '{' || text[len(text)-1] != '}' {
			return UUID{}, fmt.Errorf("uuid: invalid UUID string: %q", text)
		}
		text = text[1 : len(text)-1]

	default:
		return UUID{}, fmt.Errorf("uuid: invalid UUID string: %q", text)
	}
	if text[8] != '-' || text[13] != '-' || text[18] != '-' || text[23] != '-' {
		return UUID{}, fmt.Errorf("uuid: invalid UUID string: %q", text)
	}

	var u UUID
	for i, x := range bytePos {
		v, ok := xtob(text[x], text[x+1])
		if !ok {
			return UUID{}, fmt.Errorf("uuid: invalid UUID string: %q", text)
		}
		u[i] = v
	}
	return u, nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// Following formats are supported:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36:

	// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + len(urnPrefixBytes):
		prefix := text[:len(urnPrefixBytes)]
		if !bytes.EqualFold(prefix, urnPrefixBytes) {
			return fmt.Errorf("uuid: invalid UUID string: %q", text)
		}
		text = text[len(urnPrefixBytes):]

	// {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	case 36 + 2:
		if text[0] != '{' || text[len(text)-1] != '}' {
			return fmt.Errorf("uuid: invalid UUID string: %q", text)
		}
		text = text[1 : len(text)-1]

	default:
		return fmt.Errorf("uuid: invalid UUID string: %q", text)
	}
	if text[8] != '-' || text[13] != '-' || text[18] != '-' || text[23] != '-' {
		return fmt.Errorf("uuid: invalid UUID string: %q", text)
	}

	var tmp UUID
	for i, x := range bytePos {
		v, ok := xtob(text[x], text[x+1])
		if !ok {
			return fmt.Errorf("uuid: invalid UUID string: %q", text)
		}
		tmp[i] = v
	}
	*u = tmp

	return nil
}

// NewV4 Wrapper over the real NewV4 method
func NewV4() UUID {
	var u UUID
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	u[6] = (u[6] & 0x0f) | 0x40 // version 4
	u[8] = (u[8] & 0x3f) | 0x80 // variant bits; RFC4122
	return u
}

// MarshalText Wrapper over the real MarshalText method
func (u UUID) MarshalText() ([]byte, error) {
	var buf [36]byte

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return buf[:], nil
}

// MarshalBinary Wrapper over the real MarshalBinary method
func (u UUID) MarshalBinary() ([]byte, error) {
	return u[:], nil
}

// UnmarshalBinary Wrapper over the real UnmarshalBinary method
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
	}
	copy(u[:], data)
	return nil
}

// String returns the canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	var buf [36]byte

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf[:])
}

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob converts hex characters x1 and x2 into a byte.
func xtob(x1, x2 byte) (byte, bool) {
	b1 := xvalues[x1]
	b2 := xvalues[x2]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}
