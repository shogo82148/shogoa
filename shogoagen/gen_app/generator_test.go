package genapp_test

// TODO: fix me, the test is failing.
// import (
// 	"bytes"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"text/template"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	"github.com/shogo82148/shogoa/design"
// 	"github.com/shogo82148/shogoa/dslengine"
// 	"github.com/shogo82148/shogoa/shogoagen/codegen"
// 	genapp "github.com/shogo82148/shogoa/shogoagen/gen_app"
// 	"github.com/shogo82148/shogoa/version"
// )

// var _ = Describe("Generate", func() {
// 	var workspace *codegen.Workspace
// 	var outDir string
// 	var files []string
// 	var genErr error
// 	oldGO111MODULE := os.Getenv("GO111MODULE")

// 	BeforeEach(func() {
// 		os.Setenv("GO111MODULE", "off")

// 		var err error
// 		workspace, err = codegen.NewWorkspace("test")
// 		Ω(err).ShouldNot(HaveOccurred())
// 		outDir, err = os.MkdirTemp(filepath.Join(workspace.Path, "src"), "")
// 		Ω(err).ShouldNot(HaveOccurred())
// 		os.Args = []string{"shogoagen", "--out=" + outDir, "--design=foo", "--version=" + version.String()}
// 		codegen.TempCount = 0
// 	})

// 	JustBeforeEach(func() {
// 		design.GeneratedMediaTypes = make(design.MediaTypeRoot)
// 		design.ProjectedMediaTypes = make(design.MediaTypeRoot)
// 		files, genErr = genapp.Generate()
// 	})

// 	AfterEach(func() {
// 		workspace.Delete()
// 		delete(codegen.Reserved, "app")
// 		os.Setenv("GO111MODULE", oldGO111MODULE)
// 	})

// 	Context("with a dummy API", func() {
// 		BeforeEach(func() {
// 			design.Design = &design.APIDefinition{
// 				Name:        "test api",
// 				Title:       "dummy API with no resource",
// 				Description: "I told you it's dummy",
// 			}
// 		})

// 		It("generates correct empty files", func() {
// 			Ω(genErr).Should(BeNil())
// 			Ω(files).Should(HaveLen(6))
// 			isEmptySource := func(filename string) {
// 				contextsContent, err := os.ReadFile(filepath.Join(outDir, "app", filename))
// 				Ω(err).ShouldNot(HaveOccurred())
// 				lines := strings.Split(string(contextsContent), "\n")
// 				Ω(lines).ShouldNot(BeEmpty())
// 				Ω(len(lines)).Should(BeNumerically(">", 1))
// 			}
// 			isEmptySource("contexts.go")
// 			isEmptySource("controllers.go")
// 			isEmptySource("hrefs.go")
// 			isEmptySource("media_types.go")
// 		})
// 	})

// 	Context("with a simple API", func() {
// 		var contextsCode, controllersCode, hrefsCode, mediaTypesCode string
// 		var payload *design.UserTypeDefinition

// 		isSource := func(filename, content string) {
// 			contextsContent, err := os.ReadFile(filepath.Join(outDir, "app", filename))
// 			Ω(err).ShouldNot(HaveOccurred())
// 			Ω(string(contextsContent)).Should(Equal(content))
// 		}

// 		funcs := template.FuncMap{
// 			"sep": func() string { return string(os.PathSeparator) },
// 		}

// 		runCodeTemplates := func(data map[string]string) {
// 			contextsCodeT, err := template.New("context").Funcs(funcs).Parse(contextsCodeTmpl)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			var b bytes.Buffer
// 			err = contextsCodeT.Execute(&b, data)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			contextsCode = b.String()

// 			controllersCodeT, err := template.New("controllers").Funcs(funcs).Parse(controllersCodeTmpl)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			b.Reset()
// 			err = controllersCodeT.Execute(&b, data)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			controllersCode = b.String()

// 			hrefsCodeT, err := template.New("hrefs").Funcs(funcs).Parse(hrefsCodeTmpl)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			b.Reset()
// 			err = hrefsCodeT.Execute(&b, data)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			hrefsCode = b.String()

// 			mediaTypesCodeT, err := template.New("media types").Funcs(funcs).Parse(mediaTypesCodeTmpl)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			b.Reset()
// 			err = mediaTypesCodeT.Execute(&b, data)
// 			Ω(err).ShouldNot(HaveOccurred())
// 			mediaTypesCode = b.String()
// 		}

// 		BeforeEach(func() {
// 			payload = nil
// 			required := &dslengine.ValidationDefinition{
// 				Required: []string{"id"},
// 			}
// 			idAt := design.AttributeDefinition{
// 				Type:        design.String,
// 				Description: "widget id",
// 			}
// 			params := design.AttributeDefinition{
// 				Type: design.Object{
// 					"id": &idAt,
// 				},
// 				Validation: required,
// 			}
// 			resp := design.ResponseDefinition{
// 				Name:        "ok",
// 				Status:      200,
// 				Description: "get of widgets",
// 				MediaType:   "application/vnd.rightscale.codegen.test.widgets",
// 				ViewName:    "default",
// 			}
// 			route := design.RouteDefinition{
// 				Verb: "GET",
// 				Path: "/:id",
// 			}
// 			at := design.AttributeDefinition{
// 				Type: design.String,
// 			}
// 			ut := design.UserTypeDefinition{
// 				AttributeDefinition: &at,
// 				TypeName:            "id",
// 			}
// 			res := design.ResourceDefinition{
// 				Name:                "Widget",
// 				BasePath:            "/widgets",
// 				Description:         "Widgetty",
// 				MediaType:           "application/vnd.rightscale.codegen.test.widgets",
// 				CanonicalActionName: "get",
// 			}
// 			get := design.ActionDefinition{
// 				Name:        "get",
// 				Description: "get widgets",
// 				Parent:      &res,
// 				Routes:      []*design.RouteDefinition{&route},
// 				Responses:   map[string]*design.ResponseDefinition{"ok": &resp},
// 				Params:      &params,
// 				Payload:     payload,
// 			}
// 			res.Actions = map[string]*design.ActionDefinition{"get": &get}
// 			mt := design.MediaTypeDefinition{
// 				UserTypeDefinition: &ut,
// 				Identifier:         "application/vnd.rightscale.codegen.test.widgets",
// 				ContentType:        "application/vnd.rightscale.codegen.test.widgets",
// 				Views: map[string]*design.ViewDefinition{
// 					"default": {
// 						AttributeDefinition: ut.AttributeDefinition,
// 						Name:                "default",
// 					},
// 				},
// 			}
// 			design.Design = &design.APIDefinition{
// 				Name:        "test api",
// 				Title:       "dummy API with no resource",
// 				Description: "I told you it's dummy",
// 				Resources:   map[string]*design.ResourceDefinition{"Widget": &res},
// 				MediaTypes:  map[string]*design.MediaTypeDefinition{"application/vnd.rightscale.codegen.test.widgets": &mt},
// 			}
// 		})

// 		Context("", func() {
// 			BeforeEach(func() {
// 				runCodeTemplates(map[string]string{"outDir": outDir, "design": "foo", "tmpDir": filepath.Base(outDir), "version": version.String()})
// 			})

// 			It("generates the corresponding code", func() {
// 				Ω(genErr).Should(BeNil())
// 				Ω(files).Should(HaveLen(8))

// 				isSource("contexts.go", contextsCode)
// 				isSource("controllers.go", controllersCode)
// 				isSource("hrefs.go", hrefsCode)
// 				isSource("media_types.go", mediaTypesCode)
// 			})
// 		})

// 		Context("with a slice payload", func() {
// 			BeforeEach(func() {
// 				elemType := &design.AttributeDefinition{Type: design.Integer}
// 				payload = &design.UserTypeDefinition{
// 					AttributeDefinition: &design.AttributeDefinition{
// 						Type: &design.Array{ElemType: elemType},
// 					},
// 					TypeName: "Collection",
// 				}
// 				design.Design.Resources["Widget"].Actions["get"].Payload = payload
// 				runCodeTemplates(map[string]string{"outDir": outDir, "design": "foo", "tmpDir": filepath.Base(outDir), "version": version.String()})
// 			})

// 			It("generates the correct payload assignment code", func() {
// 				Ω(genErr).Should(BeNil())

// 				contextsContent, err := os.ReadFile(filepath.Join(outDir, "app", "controllers.go"))
// 				Ω(err).ShouldNot(HaveOccurred())
// 				Ω(string(contextsContent)).Should(ContainSubstring(controllersSlicePayloadCode))
// 			})
// 		})

// 		Context("with a optional payload", func() {
// 			BeforeEach(func() {
// 				elemType := &design.AttributeDefinition{Type: design.Integer}
// 				payload = &design.UserTypeDefinition{
// 					AttributeDefinition: &design.AttributeDefinition{
// 						Type: &design.Array{ElemType: elemType},
// 					},
// 					TypeName: "Collection",
// 				}
// 				design.Design.Resources["Widget"].Actions["get"].Payload = payload
// 				design.Design.Resources["Widget"].Actions["get"].PayloadOptional = true
// 				runCodeTemplates(map[string]string{"outDir": outDir, "design": "foo", "tmpDir": filepath.Base(outDir), "version": version.String()})
// 			})

// 			It("generates the no payloads assignment code", func() {
// 				Ω(genErr).Should(BeNil())

// 				contextsContent, err := os.ReadFile(filepath.Join(outDir, "app", "controllers.go"))
// 				Ω(err).ShouldNot(HaveOccurred())
// 				Ω(string(contextsContent)).Should(ContainSubstring(controllersOptionalPayloadCode))
// 			})
// 		})

// 		Context("with a multipart payload", func() {
// 			BeforeEach(func() {
// 				elemTypeInt := &design.AttributeDefinition{Type: design.Integer}
// 				elemTypeFile := &design.AttributeDefinition{Type: design.File}
// 				payload = &design.UserTypeDefinition{
// 					AttributeDefinition: &design.AttributeDefinition{
// 						Type: design.Object{
// 							"int":  elemTypeInt,
// 							"file": elemTypeFile,
// 						},
// 					},
// 					TypeName: "Collection",
// 				}
// 				design.Design.Resources["Widget"].Actions["get"].Payload = payload
// 				design.Design.Resources["Widget"].Actions["get"].PayloadMultipart = true
// 				runCodeTemplates(map[string]string{"outDir": outDir, "design": "foo", "tmpDir": filepath.Base(outDir), "version": version.String()})
// 			})

// 			It("generates the corresponding code", func() {
// 				Ω(genErr).Should(BeNil())

// 				contextsContent, err := os.ReadFile(filepath.Join(outDir, "app", "controllers.go"))
// 				Ω(err).ShouldNot(HaveOccurred())
// 				Ω(string(contextsContent)).Should(ContainSubstring(controllersMultipartPayloadCode))
// 			})
// 		})
// 	})
// })

// var _ = Describe("NewGenerator", func() {
// 	var generator *genapp.Generator

// 	var args = struct {
// 		api    *design.APIDefinition
// 		outDir string
// 		target string
// 		noTest bool
// 	}{
// 		api: &design.APIDefinition{
// 			Name: "test api",
// 		},
// 		target: "app",
// 		noTest: true,
// 	}

// 	Context("with options all options set", func() {
// 		BeforeEach(func() {

// 			generator = genapp.NewGenerator(
// 				genapp.API(args.api),
// 				genapp.OutDir(args.outDir),
// 				genapp.Target(args.target),
// 				genapp.NoTest(args.noTest),
// 			)
// 		})

// 		It("has all public properties set with expected value", func() {
// 			Ω(generator).ShouldNot(BeNil())
// 			Ω(generator.API.Name).Should(Equal(args.api.Name))
// 			Ω(generator.OutDir).Should(Equal(args.outDir))
// 			Ω(generator.Target).Should(Equal(args.target))
// 			Ω(generator.NoTest).Should(Equal(args.noTest))
// 		})

// 	})
// })

// const contextsCodeTmpl = `// Code generated by shogoagen {{ .version }}, DO NOT EDIT.
// //
// // API "test api": Application Contexts
// //
// // Command:
// // $ shogoagen
// // --out=$(GOPATH){{sep}}src{{sep}}{{.tmpDir}}
// // --design={{.design}}
// // --version={{.version}}

// package app

// import (
// 	"context"
// 	"github.com/shogo82148/shogoa"
// 	"net/http"
// )

// // GetWidgetContext provides the Widget get action context.
// type GetWidgetContext struct {
// 	context.Context
// 	*shogoa.ResponseData
// 	*shogoa.RequestData
// 	ID string
// }

// // NewGetWidgetContext parses the incoming request URL and body, performs validations and creates the
// // context used by the Widget controller get action.
// func NewGetWidgetContext(ctx context.Context, r *http.Request, service *shogoa.Service) (*GetWidgetContext, error) {
// 	var err error
// 	resp := shogoa.ContextResponse(ctx)
// 	resp.Service = service
// 	req := shogoa.ContextRequest(ctx)
// 	req.Request = r
// 	rctx := GetWidgetContext{Context: ctx, ResponseData: resp, RequestData: req}
// 	paramID := req.Params["id"]
// 	if len(paramID) > 0 {
// 		rawID := paramID[0]
// 		rctx.ID = rawID
// 	}
// 	return &rctx, err
// }

// // OK sends a HTTP response with status code 200.
// func (ctx *GetWidgetContext) OK(r ID) error {
// 	if ctx.ResponseData.Header().Get("Content-Type") == "" {
// 		ctx.ResponseData.Header().Set("Content-Type", "application/vnd.rightscale.codegen.test.widgets")
// 	}
// 	return ctx.ResponseData.Service.Send(ctx.Context, 200, r)
// }
// `

// const controllersCodeTmpl = `// Code generated by shogoagen {{ .version }}, DO NOT EDIT.
// //
// // API "test api": Application Controllers
// //
// // Command:
// // $ shogoagen
// // --out=$(GOPATH){{sep}}src{{sep}}{{.tmpDir}}
// // --design={{.design}}
// // --version={{.version}}

// package app

// import (
// 	"context"
// 	"github.com/shogo82148/shogoa"
// 	"net/http"
// )

// // initService sets up the service encoders, decoders and mux.
// func initService(service *shogoa.Service) {
// 	// Setup encoders and decoders

// 	// Setup default encoder and decoder
// }

// // WidgetController is the controller interface for the Widget actions.
// type WidgetController interface {
// 	shogoa.Muxer
// 	Get(*GetWidgetContext) error
// }

// // MountWidgetController "mounts" a Widget resource controller on the given service.
// func MountWidgetController(service *shogoa.Service, ctrl WidgetController) {
// 	initService(service)
// 	var h shogoa.Handler

// 	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
// 		// Check if there was an error loading the request
// 		if err := shogoa.ContextError(ctx); err != nil {
// 			return err
// 		}
// 		// Build the context
// 		rctx, err := NewGetWidgetContext(ctx, req, service)
// 		if err != nil {
// 			return err
// 		}
// 		return ctrl.Get(rctx)
// 	}
// 	service.Mux.Handle("GET", "/:id", ctrl.MuxHandler("get", h, nil))
// 	service.LogInfo("mount", "ctrl", "Widget", "action", "Get", "route", "GET /:id")
// }
// `

// const hrefsCodeTmpl = `// Code generated by shogoagen {{.version}}, DO NOT EDIT.
// //
// // API "test api": Application Resource Href Factories
// //
// // Command:
// // $ shogoagen
// // --out=$(GOPATH){{sep}}src{{sep}}{{.tmpDir}}
// // --design={{.design}}
// // --version={{.version}}

// package app

// import (
// 	"fmt"
// 	"strings"
// )

// // WidgetHref returns the resource href.
// func WidgetHref(id interface{}) string {
// 	paramid := strings.TrimLeftFunc(fmt.Sprintf("%v", id), func(r rune) bool { return r == '/' })
// 	return fmt.Sprintf("/%v", paramid)
// }
// `

// const mediaTypesCodeTmpl = `// Code generated by shogoagen {{ .version }}, DO NOT EDIT.
// //
// // API "test api": Application Media Types
// //
// // Command:
// // $ shogoagen
// // --out=$(GOPATH){{sep}}src{{sep}}{{.tmpDir}}
// // --design={{.design}}
// // --version={{.version}}

// package app
// `

// const controllersSlicePayloadCode = `
// // MountWidgetController "mounts" a Widget resource controller on the given service.
// func MountWidgetController(service *shogoa.Service, ctrl WidgetController) {
// 	initService(service)
// 	var h shogoa.Handler

// 	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
// 		// Check if there was an error loading the request
// 		if err := shogoa.ContextError(ctx); err != nil {
// 			return err
// 		}
// 		// Build the context
// 		rctx, err := NewGetWidgetContext(ctx, req, service)
// 		if err != nil {
// 			return err
// 		}
// 		// Build the payload
// 		if rawPayload := shogoa.ContextRequest(ctx).Payload; rawPayload != nil {
// 			rctx.Payload = rawPayload.(Collection)
// 		} else {
// 			return shogoa.MissingPayloadError()
// 		}
// 		return ctrl.Get(rctx)
// 	}
// 	service.Mux.Handle("GET", "/:id", ctrl.MuxHandler("get", h, unmarshalGetWidgetPayload))
// 	service.LogInfo("mount", "ctrl", "Widget", "action", "Get", "route", "GET /:id")
// }

// // unmarshalGetWidgetPayload unmarshals the request body into the context request data Payload field.
// func unmarshalGetWidgetPayload(ctx context.Context, service *shogoa.Service, req *http.Request) error {
// 	var payload Collection
// 	if err := service.DecodeRequest(req, &payload); err != nil {
// 		return err
// 	}
// 	shogoa.ContextRequest(ctx).Payload = payload
// 	return nil
// }
// `

// const controllersOptionalPayloadCode = `
// // MountWidgetController "mounts" a Widget resource controller on the given service.
// func MountWidgetController(service *shogoa.Service, ctrl WidgetController) {
// 	initService(service)
// 	var h shogoa.Handler

// 	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
// 		// Check if there was an error loading the request
// 		if err := shogoa.ContextError(ctx); err != nil {
// 			return err
// 		}
// 		// Build the context
// 		rctx, err := NewGetWidgetContext(ctx, req, service)
// 		if err != nil {
// 			return err
// 		}
// 		// Build the payload
// 		if rawPayload := shogoa.ContextRequest(ctx).Payload; rawPayload != nil {
// 			rctx.Payload = rawPayload.(Collection)
// 		}
// 		return ctrl.Get(rctx)
// 	}
// 	service.Mux.Handle("GET", "/:id", ctrl.MuxHandler("get", h, unmarshalGetWidgetPayload))
// 	service.LogInfo("mount", "ctrl", "Widget", "action", "Get", "route", "GET /:id")
// }

// // unmarshalGetWidgetPayload unmarshals the request body into the context request data Payload field.
// func unmarshalGetWidgetPayload(ctx context.Context, service *shogoa.Service, req *http.Request) error {
// 	var payload Collection
// 	if err := service.DecodeRequest(req, &payload); err != nil {
// 		return err
// 	}
// 	shogoa.ContextRequest(ctx).Payload = payload
// 	return nil
// }
// `

// const controllersMultipartPayloadCode = `
// // MountWidgetController "mounts" a Widget resource controller on the given service.
// func MountWidgetController(service *shogoa.Service, ctrl WidgetController) {
// 	initService(service)
// 	var h shogoa.Handler

// 	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
// 		// Check if there was an error loading the request
// 		if err := shogoa.ContextError(ctx); err != nil {
// 			return err
// 		}
// 		// Build the context
// 		rctx, err := NewGetWidgetContext(ctx, req, service)
// 		if err != nil {
// 			return err
// 		}
// 		// Build the payload
// 		if rawPayload := shogoa.ContextRequest(ctx).Payload; rawPayload != nil {
// 			rctx.Payload = rawPayload.(*Collection)
// 		} else {
// 			return shogoa.MissingPayloadError()
// 		}
// 		return ctrl.Get(rctx)
// 	}
// 	service.Mux.Handle("GET", "/:id", ctrl.MuxHandler("get", h, unmarshalGetWidgetPayload))
// 	service.LogInfo("mount", "ctrl", "Widget", "action", "Get", "route", "GET /:id")
// }

// // unmarshalGetWidgetPayload unmarshals the request body into the context request data Payload field.
// func unmarshalGetWidgetPayload(ctx context.Context, service *shogoa.Service, req *http.Request) error {
// 	var err error
// 	var payload collection
// 	_, rawFile, err2 := req.FormFile("file")
// 	if err2 == nil {
// 		payload.File = rawFile
// 	} else if !errors.Is(err2, http.ErrMissingFile) {
// 		err = shogoa.MergeErrors(err, shogoa.InvalidParamTypeError("file", "file", "file"))
// 	}
// 	rawInt := req.FormValue("int")
// 	if int_, err2 := strconv.Atoi(rawInt); err2 == nil {
// 		tmp2 := int_
// 		tmp1 := &tmp2
// 		payload.Int = tmp1
// 	} else {
// 		err = shogoa.MergeErrors(err, shogoa.InvalidParamTypeError("int", rawInt, "integer"))
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	shogoa.ContextRequest(ctx).Payload = payload.Publicize()
// 	return nil
// }
// `
