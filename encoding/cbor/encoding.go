package cbor

import (
	"io"

	"github.com/shogo82148/shogoa"
	"github.com/ugorji/go/codec"
)

var (
	// Handle used by encoder and decoder.
	Handle codec.CborHandle

	// Enforce that codec.Decoder satisfies shogoa.ResettableDecoder at compile time
	_ shogoa.ResettableDecoder = (*codec.Decoder)(nil)
	_ shogoa.ResettableEncoder = (*codec.Encoder)(nil)
)

// NewDecoder returns a cbor decoder.
func NewDecoder(r io.Reader) shogoa.Decoder {
	return codec.NewDecoder(r, &Handle)
}

// NewEncoder returns a cbor encoder.
func NewEncoder(w io.Writer) shogoa.Encoder {
	return codec.NewEncoder(w, &Handle)
}
