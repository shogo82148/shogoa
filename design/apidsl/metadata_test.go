package apidsl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/shogoa/dslengine"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
	}{
		{
			name: "blank metadata string",
			key:  "",
			val:  "",
		},
		{
			name: "valid metadata string",
			key:  "struct:tag:json",
			val:  "myName,omitempty",
		},
		{
			name: "unicode metadata string",
			key:  "abc123一二三",
			val:  "˜µ≤≈ç√",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// declaration design
			dslengine.Reset()
			api := API("Example API", func() {
				Metadata(tt.key, tt.val)
				BasicAuthSecurity("password")
			})

			rd := Resource("Example Resource", func() {
				Metadata(tt.key, tt.val)
				Action("Example Action", func() {
					Metadata(tt.key, tt.val)
					Routing(
						GET("/", func() {
							Metadata(tt.key, tt.val)
						}),
					)
					Security("password", func() {
						Metadata(tt.key, tt.val)
					})
				})
				Response("Example Response", func() {
					Metadata(tt.key, tt.val)
				})
			})

			mtd := MediaType("Example MediaType", func() {
				Metadata(tt.key, tt.val)
				Attribute("Example Attribute", func() {
					Metadata(tt.key, tt.val)
				})
			})
			_ = dslengine.Run()

			// assertion
			expected := dslengine.MetadataDefinition{tt.key: {tt.val}}
			if diff := cmp.Diff(expected, api.Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
			if diff := cmp.Diff(expected, rd.Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
			if diff := cmp.Diff(expected, rd.Actions["Example Action"].Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
			if diff := cmp.Diff(expected, rd.Actions["Example Action"].Routes[0].Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
			if diff := cmp.Diff(expected, rd.Actions["Example Action"].Security.Scheme.Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
			if diff := cmp.Diff(expected, rd.Responses["Example Response"].Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
			if diff := cmp.Diff(expected, mtd.Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}

			mtdAttribute := mtd.Type.ToObject()["Example Attribute"]
			if diff := cmp.Diff(expected, mtdAttribute.Metadata); diff != "" {
				t.Errorf("unexpected metadata(-want/+got):\n%s", diff)
			}
		})
	}

	t.Run("no Metadata declaration", func(t *testing.T) {
		// declaration design
		dslengine.Reset()
		api := API("Example API", func() {})
		rd := Resource("Example Resource", func() {
			Action("Example Action", func() {})
			Response("Example Response", func() {})
		})
		mtd := MediaType("Example MediaType", func() {
			Attribute("Example Attribute", func() {})
		})

		_ = dslengine.Run()

		// assertion
		if api.Metadata != nil {
			t.Errorf("unexpected metadata: %+v", api.Metadata)
		}
		if rd.Metadata != nil {
			t.Errorf("unexpected metadata: %+v", rd.Metadata)
		}
		if mtd.Metadata != nil {
			t.Errorf("unexpected metadata: %+v", mtd.Metadata)
		}

		mtdAttribute := mtd.Type.ToObject()["Example Attribute"]
		if mtdAttribute.Metadata != nil {
			t.Errorf("unexpected metadata: %+v", mtdAttribute.Metadata)
		}
	})
}
