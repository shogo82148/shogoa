package codegen

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/template"

	"github.com/shogo82148/shogoa/design"
)

var (
	validationFuncs = template.FuncMap{
		"tabs":     Tabs,
		"slice":    toSlice,
		"oneof":    oneof,
		"constant": constant,
		"goifyAtt": GoifyAtt,
		"add":      Add,
	}
	enumValT     = template.Must(template.New("enum").Funcs(validationFuncs).Parse(enumValTmpl))
	formatValT   = template.Must(template.New("format").Funcs(validationFuncs).Parse(formatValTmpl))
	patternValT  = template.Must(template.New("pattern").Funcs(validationFuncs).Parse(patternValTmpl))
	minMaxValT   = template.Must(template.New("minMax").Funcs(validationFuncs).Parse(minMaxValTmpl))
	lengthValT   = template.Must(template.New("length").Funcs(validationFuncs).Parse(lengthValTmpl))
	requiredValT = template.Must(template.New("required").Funcs(validationFuncs).Parse(requiredValTmpl))
)

// Validator is the code generator for the 'Validate' type methods.
type Validator struct {
	arrayValT *template.Template
	hashValT  *template.Template
	userValT  *template.Template
	seen      map[string]*bytes.Buffer
}

// NewValidator instantiates a validate code generator.
func NewValidator() *Validator {
	v := &Validator{seen: map[string]*bytes.Buffer{}}
	fm := template.FuncMap{
		"tabs":             Tabs,
		"slice":            toSlice,
		"oneof":            oneof,
		"constant":         constant,
		"goifyAtt":         GoifyAtt,
		"add":              Add,
		"recurseAttribute": v.recurseAttribute,
	}
	v.arrayValT = template.Must(template.New("array").Funcs(fm).Parse(arrayValTmpl))
	v.hashValT = template.Must(template.New("hash").Funcs(fm).Parse(hashValTmpl))
	v.userValT = template.Must(template.New("user").Funcs(fm).Parse(userValTmpl))
	return v
}

// Code produces Go code that runs the validation checks recursively over the given attribute.
func (v *Validator) Code(att *design.AttributeDefinition, nonzero, required, hasDefault bool, target, context string, depth int, private bool) string {
	buf := v.recurse(att, nonzero, required, hasDefault, target, context, depth, private)
	return buf.String()
}

func (v *Validator) arrayValCode(att *design.AttributeDefinition, nonzero, required, hasDefault bool, target, context string, depth int, private bool) []byte {
	a := att.Type.ToArray()
	if a == nil {
		return nil
	}

	var buf bytes.Buffer

	// Perform any validation on the array type such as MinLength, MaxLength, etc.
	validation := ValidationChecker(att, nonzero, required, hasDefault, target, context, depth, private)
	first := true
	if validation != "" {
		buf.WriteString(validation)
		first = false
	}
	val := v.Code(a.ElemType, true, false, false, "e", context+"[*]", depth+1, false)
	if val != "" {
		switch a.ElemType.Type.(type) {
		case *design.UserTypeDefinition, *design.MediaTypeDefinition:
			// For user and media types, call the Validate method
			val = RunTemplate(v.userValT, map[string]interface{}{
				"depth":  depth + 2,
				"target": "e",
			})
			val = fmt.Sprintf("%sif e != nil {\n%s\n%s}", Tabs(depth+1), val, Tabs(depth+1))
		}
		data := map[string]interface{}{
			"elemType":   a.ElemType,
			"context":    context,
			"target":     target,
			"depth":      1,
			"private":    private,
			"validation": val,
		}
		validation = RunTemplate(v.arrayValT, data)
		if !first {
			buf.WriteByte('\n')
		} else {
			first = false
		}
		buf.WriteString(validation)
	}
	_ = first // suppress: ineffectual assignment to first (ineffassign)
	return buf.Bytes()
}

func (v *Validator) hashValCode(att *design.AttributeDefinition, nonzero, required, hasDefault bool, target, context string, depth int, private bool) []byte {
	h := att.Type.ToHash()
	if h == nil {
		return nil
	}

	var buf bytes.Buffer

	// Perform any validation on the hash type such as MinLength, MaxLength, etc.
	validation := ValidationChecker(att, nonzero, required, hasDefault, target, context, depth, private)
	first := true
	if validation != "" {
		buf.WriteString(validation)
		first = false
	}
	keyVal := v.Code(h.KeyType, true, false, false, "k", context+"[*]", depth+1, false)
	if keyVal != "" {
		switch h.KeyType.Type.(type) {
		case *design.UserTypeDefinition, *design.MediaTypeDefinition:
			// For user and media types, call the Validate method
			keyVal = RunTemplate(v.userValT, map[string]interface{}{
				"depth":  depth + 2,
				"target": "k",
			})
			keyVal = fmt.Sprintf("%sif e != nil {\n%s\n%s}", Tabs(depth+1), keyVal, Tabs(depth+1))
		}
	}
	elemVal := v.Code(h.ElemType, true, false, false, "e", context+"[*]", depth+1, false)
	if elemVal != "" {
		switch h.ElemType.Type.(type) {
		case *design.UserTypeDefinition, *design.MediaTypeDefinition:
			// For user and media types, call the Validate method
			elemVal = RunTemplate(v.userValT, map[string]interface{}{
				"depth":  depth + 2,
				"target": "e",
			})
			elemVal = fmt.Sprintf("%sif e != nil {\n%s\n%s}", Tabs(depth+1), elemVal, Tabs(depth+1))
		}
	}
	if keyVal != "" || elemVal != "" {
		data := map[string]interface{}{
			"depth":          1,
			"target":         target,
			"keyValidation":  keyVal,
			"elemValidation": elemVal,
		}
		validation = RunTemplate(v.hashValT, data)
		if !first {
			buf.WriteByte('\n')
		} else {
			first = false
		}
		buf.WriteString(validation)
	}
	_ = first // suppress: ineffectual assignment to first (ineffassign)
	return buf.Bytes()
}

func (v *Validator) recurse(att *design.AttributeDefinition, nonzero, required, hasDefault bool, target, context string, depth int, private bool) *bytes.Buffer {
	var (
		buf   = new(bytes.Buffer)
		first = true
	)

	// Break infinite recursions
	switch dt := att.Type.(type) {
	case *design.MediaTypeDefinition:
		if buf, ok := v.seen[dt.TypeName]; ok {
			return buf
		}
		v.seen[dt.TypeName] = buf
	case *design.UserTypeDefinition:
		if buf, ok := v.seen[dt.TypeName]; ok {
			return buf
		}
		v.seen[dt.TypeName] = buf
	}

	if o := att.Type.ToObject(); o != nil {
		if ds, ok := att.Type.(design.DataStructure); ok {
			att = ds.Definition()
		}
		validation := ValidationChecker(att, nonzero, required, hasDefault, target, context, depth, private)
		if validation != "" {
			buf.WriteString(validation)
			first = false
		}
		o.IterateAttributes(func(n string, catt *design.AttributeDefinition) error {
			validation := v.recurseAttribute(att, catt, n, target, context, depth, private)
			if validation != "" {
				if !first {
					buf.WriteByte('\n')
				} else {
					first = false
				}
				buf.WriteString(validation)
			}
			return nil
		})
	} else if a := att.Type.ToArray(); a != nil {
		buf.Write(v.arrayValCode(att, nonzero, required, hasDefault, target, context, depth, private))
	} else if h := att.Type.ToHash(); h != nil {
		buf.Write(v.hashValCode(att, nonzero, required, hasDefault, target, context, depth, private))
	} else {
		validation := ValidationChecker(att, nonzero, required, hasDefault, target, context, depth, private)
		if validation != "" {
			buf.WriteString(validation)
		}
	}
	return buf
}

func (v *Validator) recurseAttribute(att, catt *design.AttributeDefinition, n, target, context string, depth int, private bool) string {
	var validation string
	if _, ok := catt.Type.(design.DataStructure); ok {
		validation = RunTemplate(v.userValT, map[string]interface{}{
			"depth":  depth,
			"target": fmt.Sprintf("%s.%s", target, GoifyAtt(catt, n, true)),
		})
	} else {
		dp := depth
		if catt.Type.IsObject() {
			dp++
		}
		validation = v.recurse(
			catt,
			att.IsNonZero(n),
			att.IsRequired(n),
			att.HasDefaultValue(n),
			fmt.Sprintf("%s.%s", target, GoifyAtt(catt, n, true)),
			fmt.Sprintf("%s.%s", context, n),
			dp,
			private,
		).String()
	}
	if validation != "" {
		if catt.Type.IsObject() {
			validation = fmt.Sprintf("%sif %s.%s != nil {\n%s\n%s}",
				Tabs(depth), target, GoifyAtt(catt, n, true), validation, Tabs(depth))
		}
	}
	return validation
}

// ValidationChecker produces Go code that runs the validation defined in the given attribute
// definition against the content of the variable named target recursively.
// context is used to keep track of recursion to produce helpful error messages in case of type
// validation error.
// The generated code assumes that there is a pre-existing "err" variable of type
// error. It initializes that variable in case a validation fails.
// Note: we do not want to recurse here, recursion is done by the marshaler/unmarshaler code.
func ValidationChecker(att *design.AttributeDefinition, nonzero, required, hasDefault bool, target, context string, depth int, private bool) string {
	if att.Validation == nil {
		return ""
	}
	t := target
	isPointer := private || (!required && !hasDefault && !nonzero)
	if isPointer && att.Type.IsPrimitive() {
		t = "*" + t
	}
	data := map[string]interface{}{
		"attribute": att,
		"isPointer": private || isPointer,
		"nonzero":   nonzero,
		"context":   context,
		"target":    target,
		"targetVal": t,
		"string":    att.Type.Kind() == design.StringKind,
		"array":     att.Type.IsArray(),
		"hash":      att.Type.IsHash(),
		"depth":     depth,
		"private":   private,
	}
	res := validationsCode(att, data)
	return strings.Join(res, "\n")
}

func validationsCode(att *design.AttributeDefinition, data map[string]interface{}) (res []string) {
	validation := att.Validation
	if values := validation.Values; values != nil {
		data["values"] = values
		if val := RunTemplate(enumValT, data); val != "" {
			res = append(res, val)
		}
	}
	if format := validation.Format; format != "" {
		data["format"] = format
		if val := RunTemplate(formatValT, data); val != "" {
			res = append(res, val)
		}
	}
	if pattern := validation.Pattern; pattern != "" {
		data["pattern"] = pattern
		if val := RunTemplate(patternValT, data); val != "" {
			res = append(res, val)
		}
	}
	if min := validation.Minimum; min != nil {
		if att.Type == design.Integer {
			data["min"] = renderInteger(*min)
		} else {
			data["min"] = fmt.Sprintf("%f", *min)
		}
		data["isMin"] = true
		delete(data, "max")
		if val := RunTemplate(minMaxValT, data); val != "" {
			res = append(res, val)
		}
	}
	if max := validation.Maximum; max != nil {
		if att.Type == design.Integer {
			data["max"] = renderInteger(*max)
		} else {
			data["max"] = fmt.Sprintf("%f", *max)
		}
		data["isMin"] = false
		delete(data, "min")
		if val := RunTemplate(minMaxValT, data); val != "" {
			res = append(res, val)
		}
	}
	if minLength := validation.MinLength; minLength != nil {
		data["minLength"] = minLength
		data["isMinLength"] = true
		delete(data, "maxLength")
		if val := RunTemplate(lengthValT, data); val != "" {
			res = append(res, val)
		}
	}
	if maxLength := validation.MaxLength; maxLength != nil {
		data["maxLength"] = maxLength
		data["isMinLength"] = false
		delete(data, "minLength")
		if val := RunTemplate(lengthValT, data); val != "" {
			res = append(res, val)
		}
	}
	if required := validation.Required; len(required) > 0 {
		var val string
		for i, r := range required {
			if i > 0 {
				val += "\n"
			}
			data["required"] = r
			val += RunTemplate(requiredValT, data)
		}
		res = append(res, val)
	}
	return
}

// renderInteger renders a max or min value properly, taking into account
// overflows due to casting from a float value.
func renderInteger(f float64) string {
	if f > math.Nextafter(float64(math.MaxInt64), 0) {
		return fmt.Sprintf("%d", int64(math.MaxInt64))
	}
	if f < math.Nextafter(float64(math.MinInt64), 0) {
		return fmt.Sprintf("%d", int64(math.MinInt64))
	}
	return fmt.Sprintf("%d", int64(f))
}

// oneof produces code that compares target with each element of vals and ORs
// the result, e.g. "target == 1 || target == 2".
func oneof(target string, vals []interface{}) string {
	elems := make([]string, len(vals))
	for i, v := range vals {
		elems[i] = fmt.Sprintf("%s == %#v", target, v)
	}
	return strings.Join(elems, " || ")
}

// constant returns the Go constant name of the format with the given value.
func constant(formatName string) string {
	switch formatName {
	case "date":
		return "shogoa.FormatDate"
	case "date-time":
		return "shogoa.FormatDateTime"
	case "email":
		return "shogoa.FormatEmail"
	case "hostname":
		return "shogoa.FormatHostname"
	case "ipv4":
		return "shogoa.FormatIPv4"
	case "ipv6":
		return "shogoa.FormatIPv6"
	case "ip":
		return "shogoa.FormatIP"
	case "uri":
		return "shogoa.FormatURI"
	case "mac":
		return "shogoa.FormatMAC"
	case "cidr":
		return "shogoa.FormatCIDR"
	case "regexp":
		return "shogoa.FormatRegexp"
	case "rfc1123":
		return "shogoa.FormatRFC1123"
	}
	panic("unknown format") // bug
}

const (
	arrayValTmpl = `{{ tabs .depth }}for _, e := range {{ .target }} {
{{ .validation }}
{{ tabs .depth }}}`

	hashValTmpl = `{{ tabs .depth }}for {{ if .keyValidation }}k{{ else }}_{{ end }}, {{ if .elemValidation }}e{{ else }}_{{ end }} := range {{ .target }} {
{{- if .keyValidation }}
{{ .keyValidation }}{{ end }}{{ if .elemValidation }}
{{ .elemValidation }}{{ end }}
{{ tabs .depth }}}`

	userValTmpl = `{{ tabs .depth }}if err2 := {{ .target }}.Validate(); err2 != nil {
{{ tabs .depth }}	err = shogoa.MergeErrors(err, err2)
{{ tabs .depth }}}`

	enumValTmpl = `{{ $depth := or (and .isPointer (add .depth 1)) .depth }}{{/*
*/}}{{ if .isPointer }}{{ tabs .depth }}if {{ .target }} != nil {
{{ end }}{{ tabs $depth }}if !({{ oneof .targetVal .values }}) {
{{ tabs $depth }}	err = shogoa.MergeErrors(err, shogoa.InvalidEnumValueError(` + "`" + `{{ .context }}` + "`" + `, {{ .targetVal }}, {{ slice .values }}))
{{ if .isPointer }}{{ tabs $depth }}}
{{ end }}{{ tabs .depth }}}`

	patternValTmpl = `{{ $depth := or (and .isPointer (add .depth 1)) .depth }}{{/*
*/}}{{ if .isPointer }}{{ tabs .depth }}if {{ .target }} != nil {
{{ end }}{{ tabs $depth }}if ok := shogoa.ValidatePattern(` + "`{{ .pattern }}`" + `, {{ .targetVal }}); !ok {
{{ tabs $depth }}	err = shogoa.MergeErrors(err, shogoa.InvalidPatternError(` + "`" + `{{ .context }}` + "`" + `, {{ .targetVal }}, ` + "`{{ .pattern }}`" + `))
{{ tabs $depth }}}{{ if .isPointer }}
{{ tabs .depth }}}{{ end }}`

	formatValTmpl = `{{ $depth := or (and .isPointer (add .depth 1)) .depth }}{{/*
*/}}{{ if .isPointer }}{{ tabs .depth }}if {{ .target }} != nil {
{{ end }}{{ tabs $depth }}if err2 := shogoa.ValidateFormat({{ constant .format }}, {{ .targetVal }}); err2 != nil {
{{ tabs $depth }}		err = shogoa.MergeErrors(err, shogoa.InvalidFormatError(` + "`" + `{{ .context }}` + "`" + `, {{ .targetVal }}, {{ constant .format }}, err2))
{{ if .isPointer }}{{ tabs $depth }}}
{{ end }}{{ tabs .depth }}}`

	minMaxValTmpl = `{{ $depth := or (and .isPointer (add .depth 1)) .depth }}{{/*
*/}}{{ if .isPointer }}{{ tabs .depth }}if {{ .target }} != nil {
{{ end }}{{ tabs .depth }}	if {{ .targetVal }} {{ if .isMin }}<{{ else }}>{{ end }} {{ if .isMin }}{{ .min }}{{ else }}{{ .max }}{{ end }} {
{{ tabs $depth }}	err = shogoa.MergeErrors(err, shogoa.InvalidRangeError(` + "`" + `{{ .context }}` + "`" + `, {{ .targetVal }}, {{ if .isMin }}{{ .min }}, true{{ else }}{{ .max }}, false{{ end }}))
{{ if .isPointer }}{{ tabs $depth }}}
{{ end }}{{ tabs .depth }}}`

	lengthValTmpl = `{{$depth := or (and .isPointer (add .depth 1)) .depth}}{{/*
*/}}{{$target := or (and (or (or .array .hash) .nonzero) .target) .targetVal}}{{/*
*/}}{{if .isPointer}}{{tabs .depth}}if {{.target}} != nil {
{{end}}{{tabs .depth}}	if {{if .string}}utf8.RuneCountInString({{$target}}){{else}}len({{$target}}){{end}} {{if .isMinLength}}<{{else}}>{{end}} {{if .isMinLength}}{{.minLength}}{{else}}{{.maxLength}}{{end}} {
{{tabs $depth}}	err = shogoa.MergeErrors(err, shogoa.InvalidLengthError(` + "`" + `{{.context}}` + "`" + `, {{$target}}, {{if .string}}utf8.RuneCountInString({{$target}}){{else}}len({{$target}}){{end}}, {{if .isMinLength}}{{.minLength}}, true{{else}}{{.maxLength}}, false{{end}}))
{{if .isPointer}}{{tabs $depth}}}
{{end}}{{tabs .depth}}}`

	requiredValTmpl = `{{ $att := index $.attribute.Type.ToObject .required }}{{/*
*/}}{{ if or $.private (not $att.Type.IsPrimitive) }}{{ tabs $.depth }}if {{ $.target }}.{{ goifyAtt $att .required true }} == nil {
{{ tabs $.depth }}	err = shogoa.MergeErrors(err, shogoa.MissingAttributeError(` + "`" + `{{ $.context }}` + "`" + `, "{{ .required }}"))
{{ tabs $.depth }}}{{ end }}`
)
