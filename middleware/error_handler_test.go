package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
	"github.com/shogo82148/shogoa/middleware"
)

// errorResponse contains the details of a error response. It implements ServiceError.
type errorResponse struct {
	// ID is the unique error instance identifier.
	ID string `json:"id" yaml:"id" xml:"id" form:"id"`
	// Code identifies the class of errors.
	Code string `json:"code" yaml:"code" xml:"code" form:"code"`
	// Status is the HTTP status code used by responses that cary the error.
	Status int `json:"status" yaml:"status" xml:"status" form:"status"`
	// Detail describes the specific error occurrence.
	Detail string `json:"detail" yaml:"detail" xml:"detail" form:"detail"`
	// Meta contains additional key/value pairs useful to clients.
	Meta map[string]interface{} `json:"meta,omitempty" yaml:"meta,omitempty" xml:"meta,omitempty" form:"meta,omitempty"`
}

// Error returns the error occurrence details.
func (e *errorResponse) Error() string {
	msg := fmt.Sprintf("[%s] %d %s: %s", e.ID, e.Status, e.Code, e.Detail)
	for k, v := range e.Meta {
		msg += ", " + fmt.Sprintf("%s: %v", k, v)
	}
	return msg
}

var _ = Describe("ErrorHandler", func() {
	var service *shogoa.Service
	var h shogoa.Handler
	var verbose bool

	var rw *testResponseWriter

	BeforeEach(func() {
		service = nil
		h = nil
		verbose = true
		rw = nil
	})

	JustBeforeEach(func() {
		rw = newTestResponseWriter()
		eh := middleware.ErrorHandler(service, verbose)(h)
		req, err := http.NewRequest("GET", "/foo", nil)
		Ω(err).ShouldNot(HaveOccurred())
		ctx := newContext(service, rw, req, nil)
		err = eh(ctx, rw, req)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("with a handler returning a Go error", func() {

		BeforeEach(func() {
			service = newService(nil)
			h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				return errors.New("boom")
			}
		})

		It("turns Go errors into HTTP 500 responses", func() {
			Ω(rw.Status).Should(Equal(500))
			Ω(rw.ParentHeader["Content-Type"]).Should(Equal([]string{"text/plain"}))
			Ω(string(rw.Body)).Should(Equal(`"boom"` + "\n"))
		})

		Context("not verbose", func() {
			BeforeEach(func() {
				verbose = false
			})

			It("hides the error details", func() {
				var decoded errorResponse
				Ω(rw.Status).Should(Equal(500))
				Ω(rw.ParentHeader["Content-Type"]).Should(Equal([]string{shogoa.ErrorMediaIdentifier}))
				err := service.Decoder.Decode(&decoded, bytes.NewBuffer(rw.Body), "application/json")
				Ω(err).ShouldNot(HaveOccurred())
				msg := shogoa.ErrInternal(`Internal Server Error [zzz]`).Error()
				msg = regexp.QuoteMeta(msg)
				msg = strings.Replace(msg, "zzz", ".+", 1)
				endIDidx := strings.Index(msg, "]")
				msg = `\[.*\]` + msg[endIDidx+1:]
				Ω(fmt.Sprintf("%v", decoded.Error())).Should(MatchRegexp(msg))
			})

			Context("and shogoa 500 error", func() {
				var origID string

				BeforeEach(func() {
					h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
						e := shogoa.ErrInternal("shogoa-500-boom")
						origID = e.(shogoa.ServiceError).Token()
						return e
					}
				})

				It("preserves the error ID from the original error", func() {
					var decoded errorResponse
					Ω(origID).ShouldNot(Equal(""))
					Ω(rw.Status).Should(Equal(500))
					Ω(rw.ParentHeader["Content-Type"]).Should(Equal([]string{shogoa.ErrorMediaIdentifier}))
					err := service.Decoder.Decode(&decoded, bytes.NewBuffer(rw.Body), "application/json")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(decoded.ID).Should(Equal(origID))
				})
			})

			Context("and shogoa 504 error", func() {
				BeforeEach(func() {
					meaningful := shogoa.NewErrorClass("shogoa-504-with-info", http.StatusGatewayTimeout)
					h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
						return meaningful("gatekeeper says no")
					}
				})

				It("passes the response", func() {
					var decoded errorResponse
					Ω(rw.Status).Should(Equal(http.StatusGatewayTimeout))
					Ω(rw.ParentHeader["Content-Type"]).Should(Equal([]string{shogoa.ErrorMediaIdentifier}))
					err := service.Decoder.Decode(&decoded, bytes.NewBuffer(rw.Body), "application/json")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(decoded.Code).Should(Equal("shogoa-504-with-info"))
					Ω(decoded.Detail).Should(Equal("gatekeeper says no"))
				})
			})
		})
	})

	Context("with a handler returning a shogoa error", func() {
		var gerr error

		BeforeEach(func() {
			service = newService(nil)
			gerr = shogoa.NewErrorClass("code", 418)("teapot", "foobar", 42)
			h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				return gerr
			}
		})

		It("maps shogoa errors to HTTP responses", func() {
			var decoded errorResponse
			Ω(rw.Status).Should(Equal(gerr.(shogoa.ServiceError).ResponseStatus()))
			Ω(rw.ParentHeader["Content-Type"]).Should(Equal([]string{shogoa.ErrorMediaIdentifier}))
			err := service.Decoder.Decode(&decoded, bytes.NewBuffer(rw.Body), "application/json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(decoded.Error()).Should(Equal(gerr.Error()))
		})
	})

	// TODO: FIXME
	// Context("with a handler returning a pkg errors wrapped error", func() {
	// 	var wrappedError error
	// 	var logger *testLogger
	// 	verbose = true
	// 	BeforeEach(func() {
	// 		logger = new(testLogger)
	// 		service = newService(logger)
	// 		wrappedError = pErrors.Wrap(shogoa.ErrInternal("something crazy happened"), "an error")
	// 		h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	// 			return wrappedError
	// 		}
	// 	})

	// 	It("maps pkg errors to HTTP responses", func() {
	// 		var decoded errorResponse
	// 		cause := pErrors.Cause(wrappedError)
	// 		Ω(rw.Status).Should(Equal(cause.(shogoa.ServiceError).ResponseStatus()))
	// 		Ω(rw.ParentHeader["Content-Type"]).Should(Equal([]string{shogoa.ErrorMediaIdentifier}))
	// 		err := service.Decoder.Decode(&decoded, bytes.NewBuffer(rw.Body), "application/json")
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 		Ω(decoded.Error()).Should(Equal(cause.Error()))
	// 	})
	// 	It("logs pkg errors stacktrace", func() {
	// 		var decoded errorResponse
	// 		err := service.Decoder.Decode(&decoded, bytes.NewBuffer(rw.Body), "application/json")
	// 		Ω(err).ShouldNot(HaveOccurred())
	// 		Ω(logger.ErrorEntries).Should(HaveLen(1))
	// 		data := logger.ErrorEntries[0].Data[1]
	// 		Ω(data).Should(ContainSubstring("error_handler_test.go"))
	// 	})
	// })
})
