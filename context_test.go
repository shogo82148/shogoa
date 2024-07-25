package shogoa_test

import (
	"context"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "shogoa context value " + k.name }

var _ = Describe("ResponseData", func() {
	var data *shogoa.ResponseData
	var rw http.ResponseWriter
	var req *http.Request
	var params url.Values

	BeforeEach(func() {
		var err error
		req, err = http.NewRequest("GET", "google.com", nil)
		Ω(err).ShouldNot(HaveOccurred())
		rw = &TestResponseWriter{Status: 42}
		params = url.Values{"query": []string{"value"}}
		ctx := shogoa.NewContext(context.Background(), rw, req, params)
		data = shogoa.ContextResponse(ctx)
	})

	Context("SwitchWriter", func() {
		var rwo http.ResponseWriter

		It("sets the response writer and returns the previous one", func() {
			Ω(rwo).Should(BeNil())
			rwo = data.SwitchWriter(&TestResponseWriter{Status: 43})
			Ω(rwo).ShouldNot(BeNil())
			Ω(rwo).Should(BeAssignableToTypeOf(&TestResponseWriter{}))
			trw := rwo.(*TestResponseWriter)
			Ω(trw.Status).Should(Equal(42))
		})
	})

	Context("Write", func() {
		It("should call WriteHeader(http.StatusOK) if WriteHeader has not yet been called", func() {
			_, err := data.Write(nil)
			Ω(err).Should(BeNil())
			Ω(data.Status).Should(Equal(http.StatusOK))
		})

		It("should not affect Status if WriteHeader has been called", func() {
			status := http.StatusBadRequest
			data.WriteHeader(status)
			_, err := data.Write(nil)
			Ω(err).Should(BeNil())
			Ω(data.Status).Should(Equal(status))
		})
	})

	Context("Context", func() {
		mergeContext := func(parent, child context.Context) context.Context {
			req, err := http.NewRequestWithContext(child, "GET", "google.com", nil)
			Ω(err).Should(BeNil())
			return shogoa.NewContext(parent, &TestResponseWriter{Status: 42}, req, url.Values{})
		}
		Context("Deadline", func() {
			It("should be empty if the parent and the child have no deadline", func() {
				ctx := mergeContext(context.Background(), context.Background())
				_, ok := ctx.Deadline()
				Ω(ok).Should(Equal(false))
			})
			It("should return the parent's deadline if the child have no deadline", func() {
				deadline := time.Now().Add(time.Second)
				parent, cancel := context.WithDeadline(context.Background(), deadline)
				defer cancel()
				ctx := mergeContext(parent, context.Background())
				got, ok := ctx.Deadline()
				Ω(ok).Should(Equal(true))
				Ω(got).Should(BeTemporally("~", deadline, 100*time.Millisecond))
			})
			It("should return the child's deadline if the parent have no deadline", func() {
				deadline := time.Now().Add(time.Second)
				child, cancel := context.WithDeadline(context.Background(), deadline)
				defer cancel()
				ctx := mergeContext(context.Background(), child)
				got, ok := ctx.Deadline()
				Ω(ok).Should(Equal(true))
				Ω(got).Should(BeTemporally("~", deadline, 100*time.Millisecond))
			})
			It("should return the child's deadline if it is earlier than the parent's one", func() {
				deadline1 := time.Now().Add(time.Second)
				child, cancel := context.WithDeadline(context.Background(), deadline1)
				defer cancel()
				deadline2 := time.Now().Add(2 * time.Second)
				parent, cancel := context.WithDeadline(context.Background(), deadline2)
				defer cancel()

				ctx := mergeContext(parent, child)
				got, ok := ctx.Deadline()
				Ω(ok).Should(Equal(true))
				Ω(got).Should(BeTemporally("~", deadline1, 100*time.Millisecond))
			})
			It("should return the parent's deadline if it is earlier than the child's one", func() {
				deadline1 := time.Now().Add(2 * time.Second)
				child, cancel := context.WithDeadline(context.Background(), deadline1)
				defer cancel()
				deadline2 := time.Now().Add(time.Second)
				parent, cancel := context.WithDeadline(context.Background(), deadline2)
				defer cancel()

				ctx := mergeContext(parent, child)
				got, ok := ctx.Deadline()
				Ω(ok).Should(Equal(true))
				Ω(got).Should(BeTemporally("~", deadline2, 100*time.Millisecond))
			})
		})
		Context("Done", func() {
			It("should be canceled when the parent is canceled", func() {
				deadline := time.Now().Add(time.Second)
				parent, cancel := context.WithDeadline(context.Background(), deadline)
				defer cancel()

				ctx := mergeContext(parent, context.Background())
				select {
				case <-ctx.Done():
				case <-time.After(5 * time.Second):
					Fail("timeout")
				}
				Ω(ctx.Err()).ShouldNot(BeNil())
				Ω(time.Now()).Should(BeTemporally("~", deadline, 500*time.Millisecond))
			})
			It("should be canceled when the child is canceled", func() {
				deadline := time.Now().Add(time.Second)
				child, cancel := context.WithDeadline(context.Background(), deadline)
				defer cancel()

				ctx := mergeContext(context.Background(), child)
				select {
				case <-ctx.Done():
				case <-time.After(5 * time.Second):
					Fail("timeout")
				}
				Ω(ctx.Err()).ShouldNot(BeNil())
				Ω(time.Now()).Should(BeTemporally("~", deadline, 500*time.Millisecond))
			})
		})
		Context("Value", func() {
			key := &contextKey{"key"}
			otherKey := &contextKey{"other-key"}
			It("should return the value associated with the child if it exists", func() {
				parent := context.WithValue(context.Background(), key, "parent value")
				child := context.WithValue(context.Background(), key, "child value")
				ctx := mergeContext(parent, child)
				Ω(ctx.Value(key)).Should(Equal("child value"))
			})
			It("should return the value associated with the parent if the child associates nothing", func() {
				parent := context.WithValue(context.Background(), key, "parent value")
				child := context.WithValue(context.Background(), otherKey, "child value")
				ctx := mergeContext(parent, child)
				Ω(ctx.Value(key)).Should(Equal("parent value"))
			})
		})
	})
})
