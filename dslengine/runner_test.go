package dslengine

import (
	"testing"
)

func TestReportError(t *testing.T) {
	Reset()

	const lineNumber = 14

	// NOTE: moving the line below requires updating the
	// constant above to match its number.
	ReportError("err")

	if len(Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(Errors))
	}
	err := Errors[0]
	if err.File != "runner_test.go" {
		t.Errorf("expected error file to be 'runner_test.go', got %s", err.File)
	}
	if err.Line != lineNumber {
		t.Errorf("expected error line to be %d, got %d", lineNumber, err.Line)
	}
}
