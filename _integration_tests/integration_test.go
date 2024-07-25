package tests

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestBootstrapReadme(t *testing.T) {
	cleanup("./readme/*")
	if err := shogoagen("./readme", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/readme/design"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./readme"); err != nil {
		t.Error(err.Error())
	}
}

func TestDefaultMedia(t *testing.T) {
	cleanup("./media/*")
	if err := shogoagen("./media", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/media/design"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./media"); err != nil {
		t.Error(err.Error())
	}
	b, err := os.ReadFile("./media/app/contexts.go")
	if err != nil {
		t.Fatal("failed to load contexts.go")
	}
	expected := `// CreateGreetingPayload is the Greeting create action payload.
type CreateGreetingPayload struct {
	// A required string field in the parent type.
	Message string ` + "`" + `form:"message" json:"message" yaml:"message" xml:"message"` + "`" + `
	// An optional boolean field in the parent type.
	ParentOptional *bool ` + "`" + `form:"parent_optional,omitempty" json:"parent_optional,omitempty" yaml:"parent_optional,omitempty" xml:"parent_optional,omitempty"` + "`" + `
}
`
	if !strings.Contains(string(b), expected) {
		t.Errorf("DefaultMedia attribute definitions reference failed. Generated context:\n%s", string(b))
	}
}

func TestDefaultTime(t *testing.T) {
	defer cleanup("./default-value/*")
	if err := shogoagen("./default-value", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/default-value/design"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./default-value"); err != nil {
		t.Error(err.Error())
	}
	b, err := os.ReadFile("./default-value/app/contexts.go")
	if err != nil {
		t.Fatal("failed to load contexts.go")
	}
	expected := `func NewCheckTimetestContext(ctx context.Context, r *http.Request, service *shogoa.Service) (*CheckTimetestContext, error) {
	var err error
	resp := shogoa.ContextResponse(ctx)
	resp.Service = service
	req := shogoa.ContextRequest(ctx)
	req.Request = r
	rctx := CheckTimetestContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramTimes := req.Params["times"]
	if len(paramTimes) == 0 {
		rctx.Times, err = time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	} else {
		rawTimes := paramTimes[0]
		if times, err2 := time.Parse(time.RFC3339, rawTimes); err2 == nil {
			rctx.Times = times
		} else {
			err = shogoa.MergeErrors(err, shogoa.InvalidParamTypeError("times", rawTimes, "datetime"))
		}
	}
	return &rctx, err
}`
	if !strings.Contains(string(b), expected) {
		t.Errorf("Default time attribute definitions reference failed. Generated context:\n%s", string(b))
	}
}

func TestCellar(t *testing.T) {
	cleanup("./shogoa-cellar/*")
	if err := shogoagen("./shogoa-cellar", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/shogoa-cellar/design"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./shogoa-cellar"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./shogoa-cellar/tool/cellar-cli"); err != nil {
		t.Error(err.Error())
	}
}

func TestCustomFieldName(t *testing.T) {
	cleanup("./field/*")
	if err := shogoagen("./field", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/field/design"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./field"); err != nil {
		t.Error(err.Error())
	}
	b, err := os.ReadFile("./field/app/user_types.go")
	if err != nil {
		t.Fatal("failed to load user_types.go")
	}

	expected := `// UploadPayload user type.
type UploadPayload struct {
	// A required file field in the parent type.
	FilePrimary *multipart.FileHeader ` + "`" + `form:"file1" json:"file1" yaml:"file1" xml:"file1"` + "`" + `
	// An optional file field in the parent type.
	FileSecondary *multipart.FileHeader ` + "`" + `form:"file2,omitempty" json:"file2,omitempty" yaml:"file2,omitempty" xml:"file2,omitempty"` + "`" + `
	// A required int field in the parent type.
	ID int ` + "`" + `form:"id" json:"id" yaml:"id" xml:"id"` + "`" + `
}
`
	if !strings.Contains(string(b), expected) {
		t.Errorf("UploadPayload attribute definitions reference failed. Generated user_types:\n%s", string(b))
	}

	b, err = os.ReadFile("./field/app/media_types.go")
	if err != nil {
		t.Fatal("failed to load media_types.go")
	}

	expected = `// Multimedia media type (default view)
//
// Identifier: application/vnd.multimedia+json; view=default
type Multimedia struct {
	// Media ID
	MediaID int ` + "`" + `form:"id" json:"id" yaml:"id" xml:"id"` + "`" + `
	// An optional string field in the Multimedia
	Note *string ` + "`" + `form:"optional_note,omitempty" json:"optional_note,omitempty" yaml:"optional_note,omitempty" xml:"optional_note,omitempty"` + "`" + `
	// Media URL
	MediaURL string ` + "`" + `form:"url" json:"url" yaml:"url" xml:"url"` + "`" + `
}
`
	if !strings.Contains(string(b), expected) {
		t.Errorf("Multimedia attribute definitions reference failed. Generated media_types:\n%s", string(b))
	}

	expected = `// multimedia list (default view)
//
// Identifier: application/vnd.multimedialist+json; view=default
type Multimedialist struct {
	// A required array field in the parent media type
	MediaList []*Multimedia ` + "`" + `form:"media" json:"media" yaml:"media" xml:"media"` + "`" + `
}
`
	if !strings.Contains(string(b), expected) {
		t.Errorf("Multimedialist attribute definitions reference failed. Generated media_types:\n%s", string(b))
	}
}

func TestIssue161(t *testing.T) {
	cleanup("./issue161/*")
	if err := shogoagen("./issue161", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/issue161/design"); err != nil {
		t.Error(err.Error())
	}
	if err := gobuild("./issue161"); err != nil {
		t.Error(err.Error())
	}
}

// TODO: fix me
// func TestIssue301(t *testing.T) {
// 	cleanup("./issue301/*")
// 	if err := goagen("./issue301", "bootstrap", "-d", "github.com/shogo82148/shogoa/_integration_tests/issue301/design"); err != nil {
// 		t.Error(err.Error())
// 	}
// 	if err := gobuild("./issue301"); err != nil {
// 		t.Error(err.Error())
// 	}
// 	b, err := os.ReadFile("./issue301/app/user_types.go")
// 	if err != nil {
// 		t.Fatal("failed to load user_types.go")
// 	}

// 	// include user definition type: "github.com/shogo82148/shogoa/design"
// 	expectedImport := `import (
// 	"github.com/shogo82148/shogoa"
// 	"github.com/shogo82148/shogoa/design"
// 	"time"
// )`

// 	expectedFinalize := `func (ut *issue301Type) Finalize() {
// 	var defaultPrimitiveTypeNumber float64 = 3.140000
// 	if ut.PrimitiveTypeNumber == nil {
// 		ut.PrimitiveTypeNumber = &defaultPrimitiveTypeNumber
// 	}
// 	var defaultPrimitiveTypeTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
// 	if ut.PrimitiveTypeTime == nil {
// 		ut.PrimitiveTypeTime = &defaultPrimitiveTypeTime
// 	}
// 	var defaultUserDefinitionType design.SecuritySchemeKind = 10
// 	if ut.UserDefinitionType == nil {
// 		ut.UserDefinitionType = &defaultUserDefinitionType
// 	}
// }`
// 	if !strings.Contains(string(b), expectedImport) {
// 		t.Errorf("Failed to generate 'import' block code that sets default values for user-defined types. Generated context:\n%s", string(b))
// 	}

// 	if !strings.Contains(string(b), expectedFinalize) {
// 		t.Errorf("Failed to generate 'Finalize' function code that sets default values for user-defined types. Generated context:\n%s", string(b))
// 	}
// }

func shogoagen(dir, command string, args ...string) error {
	pkg, err := build.Import("github.com/shogo82148/shogoa/goagen", "", 0)
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "run")
	for _, f := range pkg.GoFiles {
		cmd.Args = append(cmd.Args, path.Join(pkg.Dir, f))
	}
	cmd.Dir = dir
	cmd.Args = append(cmd.Args, command)
	cmd.Args = append(cmd.Args, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s\n%s", err.Error(), out)
	}
	return nil
}

func gobuild(dir string) error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s\n%s", err.Error(), out)
	}
	return nil
}

func cleanup(dir string) {
	files, err := filepath.Glob(dir)
	if err != nil {
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f, "design") {
			continue
		}
		os.RemoveAll(f)
	}
}
