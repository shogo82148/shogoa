/*
Package goalogrus contains an adapter that makes it possible to configure shogoa so it uses logrus
as logger backend.
Usage:

	logger := logrus.New()
	// Initialize logger handler using logrus package
	service.WithLogger(goalogrus.New(logger))
	// ... Proceed with configuring and starting the shogoa service

	// In handlers:
	goalogrus.Entry(ctx).Info("foo", "bar")
*/
package goalogrus

import (
	"context"
	"fmt"

	"github.com/shogo82148/shogoa"
	"github.com/sirupsen/logrus"
)

// adapter is the logrus shogoa logger adapter.
type adapter struct {
	*logrus.Entry
}

// New wraps a logrus logger into a shogoa logger.
func New(logger *logrus.Logger) shogoa.LogAdapter {
	return FromEntry(logrus.NewEntry(logger))
}

// FromEntry wraps a logrus log entry into a shogoa logger.
func FromEntry(entry *logrus.Entry) shogoa.LogAdapter {
	return &adapter{Entry: entry}
}

// Entry returns the logrus log entry stored in the given context if any, nil otherwise.
func Entry(ctx context.Context) *logrus.Entry {
	logger := shogoa.ContextLogger(ctx)
	if a, ok := logger.(*adapter); ok {
		return a.Entry
	}
	return nil
}

// Info logs messages using logrus.
func (a *adapter) Info(msg string, data ...interface{}) {
	a.Entry.WithFields(data2rus(data)).Info(msg)
}

// Warn logs message using logrus.
func (a *adapter) Warn(msg string, data ...interface{}) {
	a.Entry.WithFields(data2rus(data)).Warn(msg)
}

// Error logs errors using logrus.
func (a *adapter) Error(msg string, data ...interface{}) {
	a.Entry.WithFields(data2rus(data)).Error(msg)
}

// New creates a new logger given a context.
func (a *adapter) New(data ...interface{}) shogoa.LogAdapter {
	return &adapter{Entry: a.Entry.WithFields(data2rus(data))}
}

func data2rus(keyvals []interface{}) logrus.Fields {
	n := (len(keyvals) + 1) / 2
	res := make(logrus.Fields, n)
	for i := 0; i < len(keyvals); i += 2 {
		k := keyvals[i]
		var v interface{} = shogoa.ErrMissingLogValue
		if i+1 < len(keyvals) {
			v = keyvals[i+1]
		}
		res[fmt.Sprintf("%v", k)] = v
	}
	return res
}
