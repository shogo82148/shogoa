package goalogrus_test

import (
	"bytes"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
	goalogrus "github.com/shogo82148/shogoa/logging/logrus"
	"github.com/sirupsen/logrus"
)

var _ = Describe("goalogrus", func() {
	var logger *logrus.Logger
	var adapter shogoa.LogAdapter
	var buf bytes.Buffer

	BeforeEach(func() {
		logger = logrus.New()
		logger.Out = &buf
		adapter = goalogrus.New(logger)
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

var _ = Describe("FromEntry", func() {
	var entry *logrus.Entry
	var adapter shogoa.LogAdapter
	var buf bytes.Buffer

	BeforeEach(func() {
		logger := logrus.New()
		logger.Out = &buf
		entry = logrus.NewEntry(logger)
		adapter = goalogrus.FromEntry(entry)
	})

	It("creates an adapter that logs", func() {
		msg := "msg"
		adapter.Info(msg)
		Ω(buf.String()).Should(ContainSubstring(msg))
	})

	Context("Entry", func() {
		var ctx context.Context

		BeforeEach(func() {
			ctx = shogoa.WithLogger(context.Background(), adapter)
		})

		It("extracts the log entry", func() {
			Ω(goalogrus.Entry(ctx)).Should(Equal(entry))
		})
	})
})
