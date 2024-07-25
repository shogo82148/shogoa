package goakit_test

import (
	"bytes"

	"github.com/go-kit/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
	goakit "github.com/shogo82148/shogoa/logging/kit"
)

var _ = Describe("New", func() {
	var buf bytes.Buffer
	var logger log.Logger
	var adapter shogoa.LogAdapter

	BeforeEach(func() {
		buf.Reset()
		logger = log.NewLogfmtLogger(&buf)
		adapter = goakit.New(logger)
	})

	It("creates an adapter that logs", func() {
		msg := "msg"
		adapter.Info(msg)
		Ω(buf.String()).Should(Equal("lvl=info msg=" + msg + "\n"))
	})

	It("creates an adapter that logs", func() {
		adapter := adapter.(shogoa.WarningLogAdapter)
		msg := "msg"
		adapter.Warn(msg)
		Ω(buf.String()).Should(Equal("lvl=warn msg=" + msg + "\n"))
	})

	It("creates an adapter that logs", func() {
		msg := "msg"
		adapter.Error(msg)
		Ω(buf.String()).Should(Equal("lvl=error msg=" + msg + "\n"))
	})
})
