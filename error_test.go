package shogoa

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	gerr := &ErrorResponse{
		ID:     "foo",
		Code:   "invalid",
		Status: http.StatusBadRequest,
		Detail: "error",
		Meta:   map[string]any{"what": 42},
	}
	b, err := json.Marshal(gerr)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"id":"foo","code":"invalid","status":400,"detail":"error","meta":{"what":42}}` {
		t.Fatalf("unexpected response: %s", b)
	}
}

func TestInvalidParamTypeError(t *testing.T) {
	valErr := InvalidParamTypeError("param", 42, "43")
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "param") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "42") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "43") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestMissingParamError(t *testing.T) {
	valErr := MissingParamError("param")
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "param") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestInvalidAttributeTypeError(t *testing.T) {
	valErr := InvalidAttributeTypeError("attr", 42, "43")
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "attr") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "42") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "43") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestMissingAttributeError(t *testing.T) {
	valErr := MissingAttributeError("ctx", "attr")
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "ctx") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "attr") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestMissingHeaderError(t *testing.T) {
	valErr := MissingHeaderError("header")
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "header") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestMethodNotAllowedError(t *testing.T) {
	t.Run("multiple allowed methods", func(t *testing.T) {
		valErr := MethodNotAllowedError("POST", []string{"OPTIONS", "GET"})
		err := valErr.(*ErrorResponse)
		if !strings.Contains(err.Detail, "POST") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "one of OPTIONS, GET") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
	})

	t.Run("single allowed method", func(t *testing.T) {
		valErr := MethodNotAllowedError("POST", []string{"GET"})
		err := valErr.(*ErrorResponse)
		if !strings.Contains(err.Detail, "POST") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if strings.Contains(err.Detail, "one of") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
	})
}

func TestInvalidEnumValueError(t *testing.T) {
	valErr := InvalidEnumValueError("ctx", 42, []any{"43", "44"})
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "ctx") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "42") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, `"43", "44"`) {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestInvalidFormatError(t *testing.T) {
	valErr := InvalidFormatError("ctx", "target", FormatDateTime, errors.New("boo"))
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "ctx") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "target") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "date-time") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "boo") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestInvalidPatternError(t *testing.T) {
	valErr := InvalidPatternError("ctx", "target", "pattern")
	err := valErr.(*ErrorResponse)
	if !strings.Contains(err.Detail, "ctx") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "target") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
	if !strings.Contains(err.Detail, "pattern") {
		t.Fatalf("unexpected response: %s", err.Detail)
	}
}

func TestInvalidRangeError(t *testing.T) {
	t.Run("with an int value", func(t *testing.T) {
		valErr := InvalidRangeError("ctx", "target", 42, true)
		err := valErr.(*ErrorResponse)
		if !strings.Contains(err.Detail, "ctx") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "target") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "42") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "greater than or equal to") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
	})

	t.Run("with an float64 value", func(t *testing.T) {
		valErr := InvalidRangeError("ctx", "target", 42.42, true)
		err := valErr.(*ErrorResponse)
		if !strings.Contains(err.Detail, "ctx") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "target") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "42.42") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "greater than or equal to") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
	})
}

func TestInvalidLengthError(t *testing.T) {
	t.Run("on strings", func(t *testing.T) {
		valErr := InvalidLengthError("ctx", "target", len("target"), 42, true)
		err := valErr.(*ErrorResponse)
		if !strings.Contains(err.Detail, "ctx") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "target") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "42") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "greater than or equal to") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
	})

	t.Run("on slices", func(t *testing.T) {
		valErr := InvalidLengthError("ctx", []string{"target1", "target2"}, 2, 42, true)
		err := valErr.(*ErrorResponse)
		if !strings.Contains(err.Detail, "ctx") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, `[]string{"target1", "target2"}`) {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "42") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
		if !strings.Contains(err.Detail, "greater than or equal to") {
			t.Fatalf("unexpected response: %s", err.Detail)
		}
	})
}

// MergeableErrorResponse contains the details of a error response.
// It implements ServiceMergeableError.
type MergeableErrorResponse struct {
	*ErrorResponse
	MergeCalled int
}

// Merge will set that merge was called and return the underlying ErrorResponse.
func (e *MergeableErrorResponse) Merge(other error) error {
	e.MergeCalled++
	return e
}

func TestMergeErrors(t *testing.T) {
	t.Run("with two nil errors returns a nil error", func(t *testing.T) {
		got := MergeErrors(nil, nil)
		if got != nil {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("a nil and an Error", func(t *testing.T) {
		err := &ErrorResponse{Detail: "foo"}
		got := MergeErrors(nil, err)
		if got != err {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("a nil and a MergeableError", func(t *testing.T) {
		err := &MergeableErrorResponse{ErrorResponse: &ErrorResponse{Detail: "foo"}}
		got := MergeErrors(nil, err)
		if got != err {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("a nil and non-Error", func(t *testing.T) {
		err := errors.New("foo")
		got := MergeErrors(nil, err)
		if got.(*ErrorResponse).Detail != "foo" {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("an Error and a nil", func(t *testing.T) {
		err := &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}
		got := MergeErrors(err, nil)
		if got != err {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("two Errors", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common"}
		got := MergeErrors(err1, err2)
		if got.(*ErrorResponse).Detail != "foo; bar" {
			t.Fatalf("unexpected response: %v", got)
		}
		if got.(*ErrorResponse).Status != 42 {
			t.Fatalf("unexpected status: %v", got.(*ErrorResponse).Status)
		}
		if got.(*ErrorResponse).Code != "common" {
			t.Fatalf("unexpected code: %v", got.(*ErrorResponse).Code)
		}
	})

	t.Run("two Errors with different code", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "another"}
		got := MergeErrors(err1, err2)
		if got.(*ErrorResponse).Detail != "foo; bar" {
			t.Fatalf("unexpected response: %v", got)
		}
		if got.(*ErrorResponse).Status != http.StatusBadRequest {
			t.Fatalf("unexpected status: %v", got.(*ErrorResponse).Status)
		}
		if got.(*ErrorResponse).Code != "bad_request" {
			t.Fatalf("unexpected code: %v", got.(*ErrorResponse).Code)
		}
	})

	t.Run("two Errors with different status", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}
		err2 := &ErrorResponse{Detail: "bar", Status: 43, Code: "common"}
		got := MergeErrors(err1, err2)
		if got.(*ErrorResponse).Detail != "foo; bar" {
			t.Fatalf("unexpected response: %v", got)
		}
		if got.(*ErrorResponse).Status != http.StatusBadRequest {
			t.Fatalf("unexpected status: %v", got.(*ErrorResponse).Status)
		}
		if got.(*ErrorResponse).Code != "bad_request" {
			t.Fatalf("unexpected code: %v", got.(*ErrorResponse).Code)
		}
	})

	t.Run("two Errors with nil target metadata", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common"}
		got := MergeErrors(err1, err2)
		if got.(*ErrorResponse).Meta != nil {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("nil metadata and non-nil metadata", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common", Meta: map[string]any{"foo": 1, "bar": 2}}
		got := MergeErrors(err1, err2)
		meta := got.(*ErrorResponse).Meta
		if len(meta) != 2 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["foo"] != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["bar"] != 2 {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("non-nil metadata and nil metadata", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common", Meta: map[string]any{"baz": 1, "qux": 2}}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common"}
		got := MergeErrors(err1, err2)
		meta := got.(*ErrorResponse).Meta
		if len(meta) != 2 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["baz"] != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["qux"] != 2 {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("non-nil metadata and non-nil metadata", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common", Meta: map[string]any{"baz": 1, "qux": 2}}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common", Meta: map[string]any{"foo": 1, "bar": 2}}
		got := MergeErrors(err1, err2)

		// Check that the metadata is merged.
		meta := got.(*ErrorResponse).Meta
		if len(meta) != 4 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["foo"] != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["bar"] != 2 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["baz"] != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["qux"] != 2 {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("an Error and a MergeableError", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common", Meta: map[string]any{"baz": 1, "qux": 2}}
		err2 := &MergeableErrorResponse{ErrorResponse: &ErrorResponse{Detail: "detail"}}
		got := MergeErrors(err1, err2)
		merr := got.(*MergeableErrorResponse)
		if merr.MergeCalled != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("metadata with a common key", func(t *testing.T) {
		err1 := &ErrorResponse{Detail: "foo", Status: 42, Code: "common", Meta: map[string]any{"foo": "bar", "qux": 44}}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common", Meta: map[string]any{"foo": 43, "baz": 42}}
		got := MergeErrors(err1, err2)

		// Check that the metadata is merged.
		meta := got.(*ErrorResponse).Meta
		if len(meta) != 3 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["foo"] != 43 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["baz"] != 42 {
			t.Fatalf("unexpected response: %v", got)
		}
		if meta["qux"] != 44 {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("a MergeableError and a nil", func(t *testing.T) {
		err := &MergeableErrorResponse{ErrorResponse: &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}}
		got := MergeErrors(err, nil)
		if got != err {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("a MergeableError and an Error", func(t *testing.T) {
		err1 := &MergeableErrorResponse{ErrorResponse: &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}}
		err2 := &ErrorResponse{Detail: "bar", Status: 42, Code: "common"}
		got := MergeErrors(err1, err2)
		if got != err1 {
			t.Fatalf("unexpected response: %v", got)
		}
		if err1.MergeCalled != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
	})

	t.Run("two MergeableErrors", func(t *testing.T) {
		err1 := &MergeableErrorResponse{ErrorResponse: &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}}
		err2 := &MergeableErrorResponse{ErrorResponse: &ErrorResponse{Detail: "foo", Status: 42, Code: "common"}}
		got := MergeErrors(err1, err2)
		if got != err1 {
			t.Fatalf("unexpected response: %v", got)
		}
		if err1.MergeCalled != 1 {
			t.Fatalf("unexpected response: %v", got)
		}
	})
}
