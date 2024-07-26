package shogoatest

import (
	"io"
	"log/slog"
	"testing"

	"github.com/shogo82148/shogoa"
	"github.com/shogo82148/shogoa/middleware"
)

// TInterface is an interface for Go's testing.T and testing.B.
//
// Deprecated: use testing.TB instead.
type TInterface = testing.TB

// ResponseSetterFunc func
type ResponseSetterFunc func(resp interface{})

// Encode implements a dummy encoder that returns the value being encoded
func (r ResponseSetterFunc) Encode(v interface{}) error {
	r(v)
	return nil
}

// Service provide a general shogoa.Service used for testing purposes
func Service(logBuf io.Writer, respSetter ResponseSetterFunc) *shogoa.Service {
	s := shogoa.New("test")
	logHandler := slog.NewJSONHandler(logBuf, nil)
	s.WithLogger(shogoa.NewLogger(logHandler))
	s.Use(middleware.LogRequest(true))
	s.Use(middleware.LogResponse())
	newEncoder := func(io.Writer) shogoa.Encoder {
		return respSetter
	}
	s.Decoder.Register(shogoa.NewJSONDecoder, "*/*")
	s.Encoder.Register(newEncoder, "*/*")
	return s
}
