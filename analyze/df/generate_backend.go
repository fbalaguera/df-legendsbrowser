package df

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/robertjanetzko/LegendsBrowser2/backend/util"
)

func enumValue(s string) string {
	return strcase.ToCamel(strings.ReplaceAll(strings.ReplaceAll(s, "2", " two"), ":", "_"))
}
func enumString(s string) string { return strcase.ToDelimited(s, ' ') }

var backendTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"enum":       enumValue,
	"enumString": enumString,
}).Parse(`// Code generated by legendsbrowser; DO NOT EDIT.
package model

import (
	"github.com/robertjanetzko/LegendsBrowser2/backend/util"
	"fmt"
	"encoding/json"
)

func InitSameFields() {
	sameFields = map[string]map[string]map[string]bool{
		{{- range $name, $obj := $.Objects }}
		"{{$obj.Name}}": {
			{{- range $field := ($obj.LegendFields "plus") }}
			{{- if ne 0 (len ($obj.LegendFields "base")) }}
			"{{$field.Name}}": {
				{{- range $field2 := ($obj.LegendFields "base") }}
				{{- if eq $field.Type $field2.Type }}
				"{{ $field2.Name }}": true,
				{{- end }}
				{{- end }}
			},
			{{- end }}
			{{- end }}
		},
		{{- end }}
	}
}


{{- range $name, $obj := $.Objects }}

{{- range $fname, $field := $obj.Fields }}
{{- if eq $field.Type "enum" }}
type {{ $obj.Name }}{{ $field.Name }} int
const (
	{{ $obj.Name }}{{ $field.Name }}_Unknown {{ $obj.Name }}{{ $field.Name }} = iota
	{{- range $i, $v := $field.UniqueEnumValues }}
	{{ $obj.Name }}{{ $field.Name }}_{{ enum $v }}
	{{- end }}
)

func parse{{ $obj.Name }}{{ $field.Name }}(s string) {{ $obj.Name }}{{ $field.Name }} {
	switch s {
	{{- range $i, $v := $field.EnumValues }}
	case "{{ $v }}":
		return {{ $obj.Name }}{{ $field.Name }}_{{ enum $v }}
	{{- end }}
	}
	return {{ $obj.Name }}{{ $field.Name }}_Unknown
}

func (s {{ $obj.Name }}{{ $field.Name }}) String() string {
	switch s {
	{{- range $i, $v := $field.UniqueEnumValues }}
	case {{ $obj.Name }}{{ $field.Name }}_{{ enum $v }}:
		return "{{ enumString $v }}"
	{{- end }}
	}
	return "unknown"
}

func (s {{ $obj.Name }}{{ $field.Name }}) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

{{- end }}
{{- end }}

type {{ $obj.Name }} struct {
	{{- range $fname, $field := $obj.Fields }}
	{{- if not (and (eq $fname "type") (not (not $obj.SubTypes))) }}
	{{- if not ($field.SameField $obj) }}
	{{ $field.TypeLine }} // {{ $fname }}
	{{- end }}
	{{- end }}
	{{- end }}
	{{- if not (not $obj.SubTypes) }}
	Details {{ $obj.Name }}Details
	{{- end }}
	{{- range $fname, $field := $obj.Additional }}
	{{ $field.TypeLine }} // {{ $fname }}
	{{- end }}
}

func New{{ $obj.Name }}() *{{ $obj.Name }} {
	return &{{ $obj.Name }}{
		{{- range $fname, $field := $obj.Fields }}{{- if $field.MustInit }}{{- if not ($field.SameField $obj) }}
		{{ $field.Init }}
		{{- end }}{{- end }}{{- end }}
		{{- range $fname, $field := $obj.Additional }}{{- if $field.MustInit }}
		{{ $field.Init }}
		{{- end }}{{- end }}
	}
}

{{- if $obj.Id }}
func (x *{{ $obj.Name }}) Id() int { return x.Id_ }
func (x *{{ $obj.Name }}) setId(id int) { x.Id_ = id }
{{- end }}
{{- if $obj.Named }}
func (x *{{ $obj.Name }}) Name() string { return x.Name_ }
{{- end }}
{{- if $obj.SubType }}
func (x *{{ $obj.Name }}) Type() string { return "{{ $obj.SubType }}" }
{{- end }}
{{- if $obj.SubType }}
func (x *{{ $obj.Name }}) RelatedToEntity(id int) bool { return {{ $obj.RelatedToEntity }} }
func (x *{{ $obj.Name }}) RelatedToHf(id int) bool { return {{ $obj.RelatedToHf }} }
func (x *{{ $obj.Name }}) RelatedToArtifact(id int) bool { return {{ $obj.RelatedToArtifact }} }
func (x *{{ $obj.Name }}) RelatedToSite(id int) bool { return {{ $obj.RelatedToSite }} }
func (x *{{ $obj.Name }}) RelatedToStructure(siteId, id int) bool { return {{ $obj.RelatedToStructure }} }
func (x *{{ $obj.Name }}) RelatedToRegion(id int) bool { return {{ $obj.RelatedToRegion }} }
func (x *{{ $obj.Name }}) RelatedToWorldConstruction(id int) bool { return {{ $obj.RelatedToWorldConstruction }} }
func (x *{{ $obj.Name }}) RelatedToWrittenContent(id int) bool { return {{ $obj.RelatedToWrittenContent }} }
func (x *{{ $obj.Name }}) RelatedToDanceForm(id int) bool { return {{ $obj.RelatedToDanceForm }} }
func (x *{{ $obj.Name }}) RelatedToMusicalForm(id int) bool { return {{ $obj.RelatedToMusicalForm }} }
func (x *{{ $obj.Name }}) RelatedToPoeticForm(id int) bool { return {{ $obj.RelatedToPoeticForm }} }
{{- end }}

func (x *{{ $obj.Name }}) CheckFields() {
	{{- range $field := ($obj.LegendFields "plus") }}
	{{- if not ($field.SameField $obj) }}
	{{- range $field2 := ($obj.LegendFields "base") }}
		{{- if eq $field.Type $field2.Type }}
		{{- if eq $field.Type "int" }}
			if x.{{ $field.Name}} != x.{{ $field2.Name}} {
				sameFields["{{$obj.Name}}"]["{{ $field.Name}}"]["{{ $field2.Name}}"] = false
			}
		{{- end }}
		{{- if eq $field.Type "string" }}
			if x.{{ $field.Name}} != x.{{ $field2.Name}} && x.{{ $field.Name}} != "" && x.{{ $field2.Name}} != "" {
				sameFields["{{$obj.Name}}"]["{{ $field.Name}}"]["{{ $field2.Name}}"] = false
			}
		{{- end }}		{{- end }}
		{{- end }}
		{{- end }}
	{{- end }}
}

func (x *{{ $obj.Name }}) MarshalJSON() ([]byte, error) {
	d := make(map[string]any)
	{{- range $fname, $field := $obj.Fields }}{{- if not ($field.SameField $obj) }}{{- if not (and (eq $fname "type") (not (not $obj.SubTypes))) }}
	{{ $field.JsonMarshal }}
	{{- end }}{{- end }}{{- end }}
	{{- if not (not $obj.SubTypes) }}
	d["details"] = x.Details
	{{- end }}
	{{- range $fname, $field := $obj.Additional }}
	{{ $field.JsonMarshal }}
	{{- end }}
	return json.Marshal(d)
}

{{- end }}

// Parser

{{- range $name, $obj := $.Objects }}
{{- range $plus := $.Modes }}
func parse{{ $obj.Name }}{{ if $plus }}Plus{{ end }}(p *util.XMLParser{{ if $plus }}, obj *{{ $obj.Name }}{{ end }}) (*{{ $obj.Name }}, error) {
	{{- if not $plus }}
	var obj = New{{ $obj.Name }}()
	{{- end }}
	{{- if $plus }}
	if obj == nil {
		obj = New{{ $obj.Name }}()
	}
	{{- end }}

	for {
		t, n, err := p.Token()
		if err != nil {
			return nil, err
		}
		switch t {
		case util.StartElement:
			switch n {
			{{- range $fname, $field := $obj.Fields }}
			{{- if $field.Active $plus }}
			case "{{ $fname }}":
						{{- if and (eq $fname "type") (not (not $obj.SubTypes)) }}
						data, err := p.Value()
						if err != nil {
							return nil, err
						}
						switch string(data) {
						{{- range $sub := ($obj.ActiveSubTypes $plus) }}
						case "{{ $sub.Case }}":
							{{- if eq 1 (len $sub.Options) }}
							{{- if not $plus }}
							obj.Details, err = parse{{ $sub.Name }}(p)
							{{- else }}
							obj.Details, err = parse{{ $sub.Name }}Plus(p, obj.Details.(*{{ $sub.Name }}))
							{{- end }}
							{{- else }}
							switch details := obj.Details.(type) {
								{{- range $opt := $sub.Options }}
							case *{{ $opt}}:
								obj.Details, err = parse{{ $opt }}Plus(p, details)
								{{- end }}
							default:
								fmt.Println("unknown subtype option", obj.Details)
								p.Skip()
							}
							{{- end }}
						{{- end }}
						default:
							p.Skip()
						}
						if err != nil {
							return nil, err
						}
						return obj, nil
						
						{{- else }}
								{{ $field.StartAction $obj $plus }}
								{{- end }}
								{{- end }}
								{{- end }}
			default:
				// fmt.Println("unknown field", n)
				p.Skip()
			}

		case util.EndElement:
				obj.CheckFields()
				return obj, nil
		}
	}
}
{{- end }}
{{- end }}
`))

var sameFields map[string]map[string]string

func LoadSameFields() error {
	data, err := ioutil.ReadFile("same.json")
	if err != nil {
		return err
	}

	sameFields = make(map[string]map[string]string)
	json.Unmarshal(data, &sameFields)
	return nil
}

func GenerateBackendCode(objects *Metadata) error {
	file, _ := json.MarshalIndent(objects, "", "  ")
	_ = ioutil.WriteFile("model.json", file, 0644)

	f, err := os.Create("../backend/model/model.go")
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	err = backendTemplate.Execute(&buf, struct {
		Objects *Metadata
		Modes   []bool
	}{
		Objects: objects,
		Modes:   []bool{false, true},
	})
	if err != nil {
		return err
	}
	p, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("WARN: could not format source", err)
		p = buf.Bytes()
	}
	_, err = f.Write(p)
	return err
}

func (f Field) FixedName() string {
	n := f.Name
	if n == "Id" || n == "Name" || n == "Type" {
		n = n + "_"
	}
	return n
}

func (f Field) TypeLine() string {
	n := f.FixedName()

	m := ""
	if f.Multiple {
		m = "[]"
	}
	t := f.Type
	if f.Type == "array" {
		t = "[]*" + *f.ElementType
	}
	if f.Type == "map" {
		t = "map[int]*" + *f.ElementType
	}
	if f.Type == "object" {
		t = "*" + *f.ElementType
	}
	if f.Type == "enum" {
		t = *f.ElementType
	}
	j := fmt.Sprintf("`json:\"%s\" legend:\"%s\"`", strcase.ToLowerCamel(f.Name), f.Legend)
	return fmt.Sprintf("%s %s%s %s", n, m, t, j)
}

func (f Field) MustInit() bool {
	return f.Type == "map" || (f.Type == "int" && !f.Multiple) || strings.HasPrefix(f.Type, "map[")
}

func (f Field) Init() string {
	n := f.FixedName()

	if f.Type == "map" {
		return fmt.Sprintf("%s: make(map[int]*%s),", n, *f.ElementType)
	} else if strings.HasPrefix(f.Type, "map[") {
		return fmt.Sprintf("%s: make(%s),", n, f.Type)
	}
	if f.Type == "int" && !f.Multiple {
		return fmt.Sprintf("%s: -1,", n)
	}

	return ""
}

func (f Field) StartAction(obj Object, plus bool) string {
	n := f.FixedName()

	if f.Type == "object" {
		var p string
		if !plus {
			p = fmt.Sprintf("v, _ := parse%s(p)", *f.ElementType)
		} else {
			p = fmt.Sprintf("v, _ := parse%sPlus(p, &%s{})", *f.ElementType, *f.ElementType)
		}
		if !f.Multiple {
			return fmt.Sprintf("%s\nobj.%s = v", p, n)
		} else {
			return fmt.Sprintf("%s\nobj.%s = append(obj.%s, v)", p, n, n)
		}
	}

	if f.Type == "array" || f.Type == "map" {
		gen := fmt.Sprintf("parse%s", *f.ElementType)

		if f.Type == "array" {
			if !plus {
				return fmt.Sprintf("parseArray(p, &obj.%s, %s)", f.Name, gen)
			} else {
				gen = fmt.Sprintf("parse%sPlus", *f.ElementType)
				return fmt.Sprintf("parseArrayPlus(p, &obj.%s, %s)", f.Name, gen)
			}
		}

		if f.Type == "map" {
			if !plus {
				return fmt.Sprintf("parseMap(p, &obj.%s, %s)", f.Name, gen)
			} else {
				gen = fmt.Sprintf("parse%sPlus", *f.ElementType)
				return fmt.Sprintf("parseMapPlus(p, &obj.%s, %s)", f.Name, gen)
			}
		}
	}

	if f.Type == "int" || f.Type == "string" || f.Type == "bool" || f.Type == "enum" {
		n := f.FixedName()

		if f.Name == n {
			n = f.CorrectedName(obj)
		}

		s := "data, err := p.Value()\nif err != nil { return nil, err }\n"

		if !f.Multiple {
			if f.Type == "int" {
				return fmt.Sprintf("%sobj.%s = num(data)", s, n)
			} else if f.Type == "string" {
				return fmt.Sprintf("%sobj.%s = txt(data)", s, n)
			} else if f.Type == "bool" {
				s := "_, err := p.Value()\nif err != nil { return nil, err }\n"
				return fmt.Sprintf("%sobj.%s = true", s, n)
			} else if f.Type == "enum" {
				return fmt.Sprintf("%sobj.%s = parse%s%s(txt(data))", s, n, obj.Name, f.CorrectedName(obj))
			}
		} else {
			if f.Type == "int" {
				return fmt.Sprintf("%sobj.%s = append(obj.%s, num(data))", s, n, n)
			} else if f.Type == "string" {
				return fmt.Sprintf("%sobj.%s = append(obj.%s, txt(data))", s, n, n)
			} else if f.Type == "bool" {
				s := "_, err := p.Value()\nif err != nil { return nil, err }\n"
				return fmt.Sprintf("%sobj.%s = append(obj.%s, true)", s, n, n)
			} else if f.Type == "enum" {
				return fmt.Sprintf("%sobj.%s = append(obj.%s, parse%s%s(txt(data)))", s, n, n, obj.Name, f.CorrectedName(obj))
			}
		}
	}

	return ""
}

func (f Field) EndAction(obj Object) string {
	n := f.FixedName()

	if f.Name == n {
		n = f.CorrectedName(obj)
	}

	if !f.Multiple {
		if f.Type == "int" {
			return fmt.Sprintf("obj.%s = n(data)", n)
		} else if f.Type == "string" {
			return fmt.Sprintf("obj.%s = string(data)", n)
		} else if f.Type == "bool" {
			return fmt.Sprintf("obj.%s = true", n)
		} else if f.Type == "enum" {
			return fmt.Sprintf("obj.%s = parse%s%s(string(data))", n, obj.Name, f.CorrectedName(obj))
		}
	} else {
		if f.Type == "int" {
			return fmt.Sprintf("obj.%s = append(obj.%s, n(data))", n, n)
		} else if f.Type == "string" {
			return fmt.Sprintf("obj.%s = append(obj.%s, string(data))", n, n)
		} else if f.Type == "enum" {
			return fmt.Sprintf("obj.%s = append(obj.%s, parse%s%s(string(data)))", n, n, obj.Name, f.CorrectedName(obj))
		}
	}

	return ""
}

var entityRegex, _ = regexp.Compile("(civ|civ_id|enid|[^d]*entity(_id)?|^entity(_id)?|^source|^destination|^involved)(_?[0-9])?$")
var hfRegex, _ = regexp.Compile("(hfid|hf_id|hist_figure_id|histfig_id|histfig|bodies|_hf)")
var artifactRegex, _ = regexp.Compile("(item(_id)?|artifact_id)$")
var siteRegex, _ = regexp.Compile("(site_id|site)[0-9]?$")
var structureRegex, _ = regexp.Compile("(structure(_id)?)$")
var regionRegex, _ = regexp.Compile("(region_id|srid)$")
var worldConstructionRegex, _ = regexp.Compile("(wcid)$")
var writtenContentRegex, _ = regexp.Compile("^wc_id$")

var noRegex, _ = regexp.Compile("^XXXXX$")

func (obj Object) RelatedToEntity() string {
	return obj.Related("entity", entityRegex, "")
}
func (obj Object) RelatedToHf() string {
	return obj.Related("hf", hfRegex, "")
}
func (obj Object) RelatedToArtifact() string {
	return obj.Related("artifact", artifactRegex, "")
}
func (obj Object) RelatedToSite() string {
	return obj.Related("site", siteRegex, "")
}
func (obj Object) RelatedToStructure() string {
	return obj.Related("structure", structureRegex, "x.RelatedToSite(siteId)")
}
func (obj Object) RelatedToRegion() string {
	return obj.Related("region", regionRegex, "")
}
func (obj Object) RelatedToWorldConstruction() string {
	return obj.Related("worldConstruction", worldConstructionRegex, "")
}
func (obj Object) RelatedToWrittenContent() string {
	return obj.Related("writtenContent", writtenContentRegex, "")
}
func (obj Object) RelatedToDanceForm() string {
	return obj.Related("danceForm", noRegex, "")
}
func (obj Object) RelatedToMusicalForm() string {
	return obj.Related("musicalForm", noRegex, "")
}
func (obj Object) RelatedToPoeticForm() string {
	return obj.Related("poeticForm", noRegex, "")
}

func (obj Object) Related(relation string, regex *regexp.Regexp, init string) string {
	var list []string
	for n, f := range obj.Fields {
		if f.Type == "int" && !f.SameField(obj) && (relation == f.Related || regex.MatchString(n)) {
			if !f.Multiple {
				list = append(list, fmt.Sprintf("x.%s == id", f.Name))
			} else {
				list = append(list, fmt.Sprintf("containsInt(x.%s, id)", f.Name))
			}
		}
	}
	sort.Strings(list)
	if len(list) > 0 {
		l := strings.Join(list, " || ")
		if init == "" {
			return l
		} else {
			return init + "&& (" + l + ")"
		}
	}
	return "false"
}

func (obj Object) LegendFields(t string) []Field {
	var list []Field
	for _, f := range obj.Fields {
		if f.Name != "Name" && f.Name != "Id" && f.Name != "Type" && f.Legend == t && f.Type != "object" && !f.Multiple {
			list = append(list, f)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}

func (f Field) JsonMarshal() string {
	n := f.FixedName()

	if f.Type == "int" && !f.Multiple {
		return fmt.Sprintf(`if x.%s != -1 { d["%s"] = x.%s }`, n, strcase.ToLowerCamel(f.Name), n)
	}
	if f.Type == "enum" && !f.Multiple {
		return fmt.Sprintf(`if x.%s != 0 { d["%s"] = x.%s }`, n, strcase.ToLowerCamel(f.Name), n)
	}
	return fmt.Sprintf(`d["%s"] = x.%s`, strcase.ToLowerCamel(f.Name), n)
}

func (f Field) SameField(obj Object) bool {
	if f.Legend != "plus" {
		return false
	}
	_, ok := sameFields[obj.Name][f.Name]
	// fmt.Println(obj.Name, f.Name, ok)
	return ok
}

func (f Field) CorrectedName(obj Object) string {
	if f.Legend != "plus" {
		if f.Name == "LocalId" && obj.Id {
			return "Id_"
		}
		return f.Name
	}
	n, ok := sameFields[obj.Name][f.Name]
	if ok {
		return n
	}
	return f.Name
}

func (f Field) UniqueEnumValues() []string {
	vs := make(map[string]bool)
	for _, k := range *f.EnumValues {
		vs[enumValue(k)] = true
	}
	list := util.Keys(vs)
	sort.Strings(list)
	return list
}
