package uuid

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strings"
	"testing"
)

func TestNewV4(t *testing.T) {
	const count int = 10000
	seen := make(map[UUID]struct{})
	for i := 0; i < count; i++ {
		u := NewV4()
		seen[u] = struct{}{}
		if (u[6] & 0xf0) != 0x40 {
			t.Errorf("expected version 4 UUID, got %x", u[6]&0xf0)
		}
		if (u[8] & 0xc0) != 0x80 {
			t.Errorf("expected variant bits UUID, got %x", u[8]&0xc0)
		}
	}
	if len(seen) != count {
		t.Errorf("expected %d unique UUIDs, got %d", count, len(seen))
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		input string
		want  UUID
	}{
		{
			input: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
		{
			input: "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
		{
			input: "URN:UUID:6BA7B810-9DAD-11D1-80B4-00C04FD430C8",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
		{
			input: "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := FromString(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		input string
		want  UUID
	}{
		{
			input: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
		{
			input: "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
		{
			input: "URN:UUID:6BA7B810-9DAD-11D1-80B4-00C04FD430C8",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
		{
			input: "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
			want:  UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var got UUID
			err := got.UnmarshalText([]byte(tt.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestMarshalText(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	if got, err := u.MarshalText(); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if string(got) != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("expected 6ba7b810-9dad-11d1-80b4-00c04fd430c8, got %s", got)
	}
}

func TestString(t *testing.T) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	if got := u.String(); got != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Errorf("expected 6ba7b810-9dad-11d1-80b4-00c04fd430c8, got %s", got)
	}
}

func TestMarshalJSON(t *testing.T) {
	u, err := FromString("c0586f01-87b5-462b-a673-3b2dcf619091")
	if err != nil {
		t.Fatal(err)
	}

	type Payload struct {
		ID   UUID
		Name string
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(Payload{ID: u, Name: "Test"}); err != nil {
		t.Fatal(err)
	}
	want := `{"ID":"c0586f01-87b5-462b-a673-3b2dcf619091","Name":"Test"}` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	encoded := `{"ID":"c0586f01-87b5-462b-a673-3b2dcf619091","Name":"Test"}`

	var payload struct {
		ID   UUID
		Name string
	}
	decoder := json.NewDecoder(strings.NewReader(encoded))
	if err := decoder.Decode(&payload); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(payload.ID[:], []byte{0xc0, 0x58, 0x6f, 0x01, 0x87, 0xb5, 0x46, 0x2b, 0xa6, 0x73, 0x3b, 0x2d, 0xcf, 0x61, 0x90, 0x91}) {
		t.Errorf("expected c0586f01-87b5-462b-a673-3b2dcf619091, got %s", payload.ID)
	}
}

func FuzzFromString(f *testing.F) {
	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("{6ba7b810-9dad-11d1-80b4-00c04fd430c8}")
	f.Fuzz(func(t *testing.T, input string) {
		u1, err := FromString(input)
		if err != nil {
			t.Skip()
		}
		var u2 UUID
		if err := u2.UnmarshalText([]byte(input)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u1 != u2 {
			t.Errorf("expected %v, got %v", u1, u2)
		}
	})
}

func FuzzString(f *testing.F) {
	f.Add([]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8})
	f.Fuzz(func(t *testing.T, input []byte) {
		var u1 UUID
		if err := u1.UnmarshalBinary(input); err != nil {
			t.Skip()
		}
		s := u1.String()
		u2, err := FromString(s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u1 != u2 {
			t.Errorf("expected %v, got %v", u1, u2)
		}
	})
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(NewV4())
	}
}

func BenchmarkFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u, err := FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		if err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(u)
	}
}

func BenchmarkUnmarshalText(b *testing.B) {
	input := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	for i := 0; i < b.N; i++ {
		var u UUID
		err := u.UnmarshalText(input)
		if err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(u)
	}
}

func BenchmarkMarshalText(b *testing.B) {
	u := NewV4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, err := u.MarshalText()
		if err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(v)
	}
}

func BenchmarkMarshalBinary(b *testing.B) {
	u := NewV4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, err := u.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(v)
	}
}

func BenchmarkUnmarshalBinary(b *testing.B) {
	input, err := NewV4().MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var u UUID
		err := u.UnmarshalBinary(input)
		if err != nil {
			b.Fatal(err)
		}
		runtime.KeepAlive(u)
	}
}

func BenchmarkString(b *testing.B) {
	u := NewV4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(u.String())
	}
}
