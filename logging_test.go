package shogoa_test

import (
	"bytes"
	"context"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
)

var _ = Describe("Info", func() {
	Context("with a nil Log", func() {
		It("doesn't log and doesn't crash", func() {
			Ω(func() { shogoa.LogInfo(context.Background(), "foo", "bar") }).ShouldNot(Panic())
		})
	})
})

var _ = Describe("Warn", func() {
	Context("with a nil Log", func() {
		It("doesn't log and doesn't crash", func() {
			Ω(func() { shogoa.LogWarn(context.Background(), "foo", "bar") }).ShouldNot(Panic())
		})
	})
})

var _ = Describe("Error", func() {
	Context("with a nil Log", func() {
		It("doesn't log and doesn't crash", func() {
			Ω(func() { shogoa.LogError(context.Background(), "foo", "bar") }).ShouldNot(Panic())
		})
	})
})

var _ = Describe("LogAdapter", func() {
	Context("with a valid Log", func() {
		var logger shogoa.LogAdapter
		const msg = "message"
		data := []interface{}{"data", "foo"}

		var out bytes.Buffer

		BeforeEach(func() {
			stdlogger := log.New(&out, "", log.LstdFlags)
			logger = shogoa.NewLogger(stdlogger)
		})

		It("Info logs", func() {
			logger.Info(msg, data...)
			Ω(out.String()).Should(ContainSubstring(msg + " data=foo"))
		})

		It("Warn logs", func() {
			logger := logger.(shogoa.WarningLogAdapter)
			logger.Warn(msg, data...)
			Ω(out.String()).Should(ContainSubstring(msg + " data=foo"))
		})

		It("Error logs", func() {
			logger.Error(msg, data...)
			Ω(out.String()).Should(ContainSubstring(msg + " data=foo"))
		})
	})
})
