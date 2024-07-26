package shogoa_test

import (
	"bytes"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
)

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
