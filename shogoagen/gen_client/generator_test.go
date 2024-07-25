package genclient_test

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa/design"
	"github.com/shogo82148/shogoa/dslengine"
	"github.com/shogo82148/shogoa/shogoagen/codegen"
	genclient "github.com/shogo82148/shogoa/shogoagen/gen_client"
	"github.com/shogo82148/shogoa/version"
)

var _ = Describe("Generate", func() {
	var workspace *codegen.Workspace
	var outDir string
	var files []string
	var genErr error

	oldGO111MODULE := os.Getenv("GO111MODULE")
	BeforeEach(func() {
		os.Setenv("GO111MODULE", "off")
	})
	AfterEach(func() {
		os.Setenv("GO111MODULE", oldGO111MODULE)
	})

	BeforeEach(func() {
		var err error
		workspace, err = codegen.NewWorkspace("test")
		Ω(err).ShouldNot(HaveOccurred())
		outDir, err = os.MkdirTemp(filepath.Join(workspace.Path, "src"), "")
		Ω(err).ShouldNot(HaveOccurred())
		os.Args = []string{"shogoagen", "--out=" + outDir, "--design=foo", "--version=" + version.String()}
	})

	JustBeforeEach(func() {
		files, genErr = genclient.Generate()
	})

	AfterEach(func() {
		workspace.Delete()
		delete(codegen.Reserved, "client")
	})

	Context("with a basic design", func() {
		var (
			resourceHeader   string
			clientHeader     string
			userTypesHeader  string
			mediaTypesHeader string
		)

		funcs := template.FuncMap{
			"sep": func() string { return string(os.PathSeparator) },
		}

		genHeader := func(data map[string]string) string {
			clientHeaderT, err := template.New("context").Funcs(funcs).Parse(clientHeaderTmpl)
			Ω(err).ShouldNot(HaveOccurred())
			var b bytes.Buffer
			err = clientHeaderT.Execute(&b, data)
			Ω(err).ShouldNot(HaveOccurred())
			return b.String()
		}

		BeforeEach(func() {
			codegen.TempCount = 0
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct

			data := map[string]string{
				"outDir":  outDir,
				"design":  "foo",
				"tmpDir":  filepath.Base(outDir),
				"version": version.String(),
			}

			// Generate Headers
			data["title"] = "Client"
			clientHeader = genHeader(data)

			data["title"] = "foo Resource Client"
			resourceHeader = genHeader(data)

			data["title"] = "Application Media Types"
			mediaTypesHeader = genHeader(data)

			data["title"] = "Application User Types"
			userTypesHeader = genHeader(data)
		})

		It("generates code generated header", func() {
			Ω(genErr).Should(BeNil())
			content, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(content)).Should(HavePrefix(resourceHeader))

			content, err = os.ReadFile(filepath.Join(outDir, "client", "client.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(content)).Should(HavePrefix(clientHeader))

			content, err = os.ReadFile(filepath.Join(outDir, "client", "media_types.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(content)).Should(HavePrefix(mediaTypesHeader))

			content, err = os.ReadFile(filepath.Join(outDir, "client", "user_types.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(content)).Should(HavePrefix(userTypesHeader))
		})
	})

	Context("with a required UUID header", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			o := design.Object{
				"header_name": &design.AttributeDefinition{Type: design.UUID},
			}
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{Verb: "GET", Path: ""}},
								Headers: &design.AttributeDefinition{
									Type: o,
									Validation: &dslengine.ValidationDefinition{
										Required: []string{"header_name"},
									},
								}}},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates header initialization code that compiles", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			c, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			content := string(c)
			Ω(content).Should(ContainSubstring("header.Set(\"header_name\", tmp3)\n"))
		})
	})

	Context("with querystring params in path", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			o := design.Object{
				"foo": &design.AttributeDefinition{Type: design.String},
				"bar": &design.AttributeDefinition{Type: &design.Array{ElemType: &design.AttributeDefinition{Type: design.Integer}}},
				"baz": &design.AttributeDefinition{Type: design.DateTime},
				"bat": &design.AttributeDefinition{Type: design.UUID},
			}
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name:   "show",
								Params: &design.AttributeDefinition{Type: o},
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "/foo/:foo/bar/:bar/baz/:baz/bat/:bat",
									},
								},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates path initialization code that uses all defined URL params in proper format", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			c, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			content := string(c)
			Ω(content).Should(ContainSubstring("func ShowFooPath("))
			Ω(content).Should(ContainSubstring(`param0 := foo`))
			Ω(content).Should(ContainSubstring(`tmp2 := make([]string, len(bar))
	for i, e := range bar {
		tmp3 := strconv.Itoa(e)
		tmp2[i] = tmp3
	}
	param1 := strings.Join(tmp2, ",")`))
			Ω(content).Should(ContainSubstring(`param2 := baz.Format(time.RFC3339)`))
			Ω(content).Should(ContainSubstring(`param3 := bat.String()`))
			Ω(content).Should(ContainSubstring(`fmt.Sprintf("/foo/%s/bar/%s/baz/%s/bat/%s", param0, param1, param2, param3)`))
		})
	})

	Context("with jsonapi like querystring params", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			o := design.Object{
				"fields[foo]": &design.AttributeDefinition{Type: design.String},
				"fields[bar]": &design.AttributeDefinition{Type: &design.Array{ElemType: &design.AttributeDefinition{Type: design.String}}},
				"fields[baz]": &design.AttributeDefinition{Type: &design.Array{ElemType: &design.AttributeDefinition{Type: design.Integer}}},
				"fields[bat]": &design.AttributeDefinition{Type: design.DateTime},
			}
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
								QueryParams: &design.AttributeDefinition{Type: o},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates param initialization code that uses the param name given in the design", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			c, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			content := string(c)
			Ω(content).Should(ContainSubstring("func ShowFooPath("))
			Ω(content).Should(ContainSubstring(`values.Set("fields[foo]", *fieldsFoo)`))
			Ω(content).Should(ContainSubstring(`	for _, p := range fieldsBar {
		tmp3 := p
		values.Add("fields[bar]", tmp3)
	}
`))
			Ω(content).Should(ContainSubstring(`	for _, p := range fieldsBaz {
		tmp5 := strconv.Itoa(p)
		values.Add("fields[baz]", tmp5)
	}
`))
			Ω(content).Should(ContainSubstring(`		tmp4 := fieldsBat.Format(time.RFC3339)
		values.Set("fields[bat]", tmp4)`))
		})

		Context("with --notool", func() {
			BeforeEach(func() {
				os.Args = append(os.Args, "--notool")
			})

			It("should not return an error", func() {
				Ω(genErr).Should(BeNil())
				Ω(files).Should(HaveLen(5)) // 9, minus 4 entries for tool paths
			})
		})
	})

	Context("with an action using websocket", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			o := design.Object{
				"fields[foo]": &design.AttributeDefinition{Type: design.String},
				"fields[bar]": &design.AttributeDefinition{Type: &design.Array{ElemType: &design.AttributeDefinition{Type: design.String}}},
				"fields[baz]": &design.AttributeDefinition{Type: &design.Array{ElemType: &design.AttributeDefinition{Type: design.Integer}}},
				"fields[bat]": &design.AttributeDefinition{Type: design.DateTime},
			}
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name:    "show",
								Schemes: []string{"ws"},
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
								QueryParams: &design.AttributeDefinition{Type: o},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates param initialization code that uses the param name given in the design", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			c, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			content := string(c)
			Ω(content).Should(ContainSubstring("func ShowFooPath("))
			Ω(content).Should(ContainSubstring(`values.Set("fields[foo]", *fieldsFoo)`))
			Ω(content).Should(ContainSubstring(`	if fieldsBar != nil {
		for _, p := range fieldsBar {
			tmp3 := p
			values.Add("fields[bar]", tmp3)
		}
	}
`))
			Ω(content).Should(ContainSubstring(`	if fieldsBaz != nil {
		for _, p := range fieldsBaz {
			tmp5 := strconv.Itoa(p)
			values.Add("fields[baz]", tmp5)
		}
	}
`))
			Ω(content).Should(ContainSubstring(`		tmp4 := fieldsBat.Format(time.RFC3339)
		values.Set("fields[bat]", tmp4)`))
		})

		Context("with --notool", func() {
			BeforeEach(func() {
				os.Args = append(os.Args, "--notool")
			})

			It("should not return an error", func() {
				Ω(genErr).Should(BeNil())
				Ω(files).Should(HaveLen(5)) // 9, minus 4 entries for tool paths
			})
		})
	})

	Context("with an action with multiple routes", func() {
		BeforeEach(func() {
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
									{
										Verb: "GET",
										Path: "/foo",
									},
								},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
			showAct.Routes[1].Parent = showAct
		})

		It("generates Path function with unique names", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			content, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring("func ShowFooPath("))
			Ω(strings.Count(string(content), "func ShowFooPath(")).Should(Equal(1))
			Ω(content).Should(ContainSubstring("func ShowFooPath2("))
			Ω(strings.Count(string(content), "func ShowFooPath2(")).Should(Equal(1))
		})

		Context("with a file server", func() {
			BeforeEach(func() {
				res := design.Design.Resources["foo"]
				res.FileServers = []*design.FileServerDefinition{
					{
						Parent:      res,
						FilePath:    "/swagger/swagger.json",
						RequestPath: "/swagger.json",
					},
				}
			})

			It("generates a Download function", func() {
				Ω(genErr).Should(BeNil())
				Ω(files).Should(HaveLen(9))
				content, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
				Ω(err).ShouldNot(HaveOccurred())
				Ω(content).Should(ContainSubstring("func (c *Client) DownloadSwaggerJSON("))
			})

		})
	})

	Context("with an action with security configured", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			securitySchemeDef := &design.SecuritySchemeDefinition{
				SchemeName: "jwt-1",
				Kind:       design.JWTSecurityKind,
			}
			design.Design = &design.APIDefinition{
				Name:        "testapi",
				Title:       "dummy API with no resource",
				Description: "I told you it's dummy",
				Consumes:    design.DefaultEncoders,
				SecuritySchemes: []*design.SecuritySchemeDefinition{
					securitySchemeDef,
				},
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								QueryParams: &design.AttributeDefinition{
									Type: design.Object{
										"param": &design.AttributeDefinition{Type: design.Integer},
										"time":  &design.AttributeDefinition{Type: design.DateTime},
										"uuid":  &design.AttributeDefinition{Type: design.UUID},
									},
								},
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
								Security: &design.SecurityDefinition{
									Scheme: securitySchemeDef,
								},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates the correct client Fields", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			content, err := os.ReadFile(filepath.Join(outDir, "client", "client.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring("JWT1Signer goaclient.Signer"))
			Ω(content).Should(ContainSubstring("func (c *Client) SetJWT1Signer(signer goaclient.Signer) {\n	c.JWT1Signer = signer\n}"))
		})

		It("generates the Signer.Sign call from Action", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			content, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring(`		if err := c.JWT1Signer.Sign(req); err != nil {
			return nil, err
		}`))
		})
	})

	Context("with an action with a user type payload", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			testType := &design.UserTypeDefinition{
				AttributeDefinition: &design.AttributeDefinition{
					Type: design.Object{
						"param": &design.AttributeDefinition{Type: design.Integer},
						"time":  &design.AttributeDefinition{Type: design.DateTime},
						"uuid":  &design.AttributeDefinition{Type: design.UUID},
					},
				},
				TypeName: "TestType",
			}
			design.Design = &design.APIDefinition{
				Types: map[string]*design.UserTypeDefinition{
					"TestType": testType,
				},
				Name:        "testapi",
				Title:       "dummy API with no resource",
				Description: "I told you it's dummy",
				Consumes:    design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
								Payload: testType,
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates the user type imports", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			content, err := os.ReadFile(filepath.Join(outDir, "client", "user_types.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring("uuid \"github.com/shogo82148/shogoa/uuid\""))
		})
	})

	Context("with a multipartform action with a user type payload", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			elemType := &design.AttributeDefinition{Type: design.String}
			testType := &design.UserTypeDefinition{
				AttributeDefinition: &design.AttributeDefinition{
					Type: design.Object{
						"param": &design.AttributeDefinition{Type: design.Integer},
						"time":  &design.AttributeDefinition{Type: design.DateTime},
						"uuid":  &design.AttributeDefinition{Type: design.UUID},
						"file":  &design.AttributeDefinition{Type: design.File},
						"image": &design.AttributeDefinition{Type: design.File},
						"name":  &design.AttributeDefinition{Type: &design.Array{ElemType: elemType}},
						"addr":  &design.AttributeDefinition{Type: &design.Array{ElemType: elemType}},
					},
					Validation: &dslengine.ValidationDefinition{
						Required: []string{"uuid", "file", "name"},
					},
				},
				TypeName: "TestType",
			}
			design.Design = &design.APIDefinition{
				Types: map[string]*design.UserTypeDefinition{
					"TestType": testType,
				},
				Name:        "testapi",
				Title:       "dummy API with no resource",
				Description: "I told you it's dummy",
				Consumes:    design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
								Payload:          testType,
								PayloadMultipart: true,
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("treat non-required param as pointer type", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			content, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(content)).Should(ContainSubstring("if payload.Param != nil"))
			Ω(string(content)).Should(ContainSubstring("if payload.Image != nil"))
			Ω(string(content)).Should(ContainSubstring("if payload.Addr != nil"))
			Ω(string(content)).ShouldNot(ContainSubstring("if payload.UUID != nil"))
			Ω(string(content)).ShouldNot(ContainSubstring("if payload.File != nil"))
			Ω(string(content)).ShouldNot(ContainSubstring("if payload.Name != nil"))
			Ω(string(content)).Should(ContainSubstring("tmp_Param := *payload.Param"))
			Ω(string(content)).Should(ContainSubstring("tmp_UUID := payload.UUID"))
			Ω(string(content)).Should(ContainSubstring("fh, err := os.Open(payload.File)"))
			Ω(string(content)).Should(ContainSubstring("fh, err := os.Open(*payload.Image)"))
		})
	})

	Context("querystring params with struct:field:name metadata", func() {
		BeforeEach(func() {
			codegen.TempCount = 0
			o := design.Object{
				"foo": &design.AttributeDefinition{
					Type:     design.String,
					Metadata: dslengine.MetadataDefinition{"struct:field:name": []string{"MetaFoo"}},
				},
			}
			design.Design = &design.APIDefinition{
				Name:     "testapi",
				Consumes: design.DefaultEncoders,
				Resources: map[string]*design.ResourceDefinition{
					"foo": {
						Name: "foo",
						Actions: map[string]*design.ActionDefinition{
							"show": {
								Name: "show",
								Routes: []*design.RouteDefinition{
									{
										Verb: "GET",
										Path: "",
									},
								},
								QueryParams: &design.AttributeDefinition{Type: o},
							},
						},
					},
				},
			}
			fooRes := design.Design.Resources["foo"]
			showAct := fooRes.Actions["show"]
			showAct.Parent = fooRes
			showAct.Routes[0].Parent = showAct
		})

		It("generates param initialization code that uses the param name given in the design", func() {
			Ω(genErr).Should(BeNil())
			Ω(files).Should(HaveLen(9))
			c, err := os.ReadFile(filepath.Join(outDir, "client", "foo.go"))
			Ω(err).ShouldNot(HaveOccurred())
			content := string(c)
			Ω(string(content)).Should(ContainSubstring("ShowFoo(ctx context.Context, path string, metaFoo *string)"))
			Ω(string(content)).Should(ContainSubstring("NewShowFooRequest(ctx context.Context, path string, metaFoo *string)"))
			Ω(string(content)).Should(ContainSubstring("if metaFoo != nil {"))
			Ω(string(content)).Should(ContainSubstring("values.Set(\"foo\", *metaFoo)"))
		})
	})
})

var _ = Describe("NewGenerator", func() {
	var generator *genclient.Generator

	var args = struct {
		api         *design.APIDefinition
		outDir      string
		target      string
		toolDirName string
		tool        string
		noTool      bool
	}{
		api: &design.APIDefinition{
			Name: "test api",
		},
		target:      "app",
		toolDirName: "test_dir",
		tool:        "mycli",
		noTool:      true,
	}

	Context("with options all options set", func() {
		BeforeEach(func() {

			generator = genclient.NewGenerator(
				genclient.API(args.api),
				genclient.OutDir(args.outDir),
				genclient.Target(args.target),
				genclient.ToolDirName(args.toolDirName),
				genclient.Tool(args.tool),
				genclient.NoTool(args.noTool),
			)
		})

		It("has all public properties set with expected value", func() {
			Ω(generator).ShouldNot(BeNil())
			Ω(generator.API.Name).Should(Equal(args.api.Name))
			Ω(generator.OutDir).Should(Equal(args.outDir))
			Ω(generator.Target).Should(Equal(args.target))
			Ω(generator.ToolDirName).Should(Equal(args.toolDirName))
			Ω(generator.Tool).Should(Equal(args.tool))
			Ω(generator.NoTool).Should(Equal(args.noTool))
		})

	})
})

const clientHeaderTmpl = `// Code generated by shogoagen {{ .version }}, DO NOT EDIT.
//
// API "testapi": {{.title}}
//
// Command:
// $ shogoagen
// --out=$(GOPATH){{sep}}src{{sep}}{{.tmpDir}}
// --design={{.design}}
// --version={{.version}}
`
