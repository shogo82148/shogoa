//go:build go1.21
// +build go1.21

package goaslog_test

import (
	"bytes"
	"log/slog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
	goaslog "github.com/shogo82148/shogoa/logging/slog"
)

var _ = Describe("goaslog", func() {
	var handler slog.Handler
	var adapter shogoa.LogAdapter
	var buf bytes.Buffer

	BeforeEach(func() {
		handler = slog.NewJSONHandler(&buf, nil)
		adapter = goaslog.New(handler)
	})

	It("adapts info messages", func() {
		msg := "msg"
		adapter.Info(msg)
		Ω(buf.String()).Should(ContainSubstring(msg))
	})

	It("adapts warn messages", func() {
		adapter := adapter.(shogoa.WarningLogAdapter)
		msg := "msg"
		adapter.Warn(msg)
		Ω(buf.String()).Should(ContainSubstring(msg))
	})

	It("adapts error messages", func() {
		msg := "msg"
		adapter.Error(msg)
		Ω(buf.String()).Should(ContainSubstring(msg))
	})
})
