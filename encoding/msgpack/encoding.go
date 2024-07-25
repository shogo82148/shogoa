package msgpack

import (
	"io"

	"github.com/shogo82148/shogoa"
	"github.com/ugorji/go/codec"
)

// Enforce that codec.Decoder satisfies shogoa.ResettableDecoder at compile time
var (
	_ shogoa.ResettableDecoder = (*codec.Decoder)(nil)
	_ shogoa.ResettableEncoder = (*codec.Encoder)(nil)

	Handle codec.MsgpackHandle
)

// NewDecoder returns a msgpack decoder.
func NewDecoder(r io.Reader) shogoa.Decoder {
	return codec.NewDecoder(r, &Handle)
}

// NewEncoder returns a msgpack encoder.
func NewEncoder(w io.Writer) shogoa.Encoder {
	return codec.NewEncoder(w, &Handle)
}
