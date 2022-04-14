package df

import (
	"fmt"
	"text/template"

	"github.com/iancoleman/strcase"
)

type Object struct {
	Name      string           `json:"name"`
	Id        bool             `json:"id,omitempty"`
	Named     bool             `json:"named,omitempty"`
	Typed     bool             `json:"typed,omitempty"`
	SubTypes  *[]string        `json:"subtypes,omitempty"`
	SubTypeOf *string          `json:"subtypeof,omitempty"`
	Fields    map[string]Field `json:"fields"`
}

type Field struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Multiple    bool    `json:"multiple,omitempty"`
	ElementType *string `json:"elements,omitempty"`
	Legend      string  `json:"legend"`
}

func (f Field) TypeLine(objects map[string]Object) string {
	n := f.Name

	if n == "Id" || n == "Name" {
		n = n + "_"
	}

	m := ""
	if f.Multiple {
		m = "[]"
	}
	t := f.Type
	if f.Type == "array" {
		t = "[]*" + objects[*f.ElementType].Name
	}
	if f.Type == "map" {
		t = "map[int]*" + objects[*f.ElementType].Name
	}
	if f.Type == "object" {
		t = "*" + f.Name
	}
	j := fmt.Sprintf("`json:\"%s\" legend:\"%s\"`", strcase.ToLowerCamel(f.Name), f.Legend)
	return fmt.Sprintf("%s %s%s %s", n, m, t, j)
}

func (f Field) StartAction() string {
	n := f.Name

	if n == "Id" || n == "Name" {
		n = n + "_"
	}

	if f.Type == "object" {
		p := fmt.Sprintf("v, _ := parse%s(d, &t)", f.Name)
		if !f.Multiple {
			return fmt.Sprintf("%s\nobj.%s = v", p, n)
		} else {
			return fmt.Sprintf("%s\nobj.%s = append(obj.%s, v)", p, n, n)
		}
	}

	if f.Type == "array" || f.Type == "map" {
		el := strcase.ToCamel(*f.ElementType)
		gen := fmt.Sprintf("parse%s", el)

		if f.Type == "array" {
			return fmt.Sprintf("parseArray(d, &obj.%s, %s)", f.Name, gen)
		}

		if f.Type == "map" {
			return fmt.Sprintf("obj.%s = make(map[int]*%s)\nparseMap(d, &obj.%s, %s)", f.Name, el, f.Name, gen)
		}
	}

	if f.Type == "int" || f.Type == "string" {
		return "data = nil"
	}

	return ""
}

func (f Field) EndAction() string {
	n := f.Name

	if n == "Id" || n == "Name" {
		n = n + "_"
	}

	if !f.Multiple {
		if f.Type == "int" {
			return fmt.Sprintf("obj.%s = n(data)", n)
		} else if f.Type == "string" {
			return fmt.Sprintf("obj.%s = string(data)", n)
		}
	} else {
		if f.Type == "int" {
			return fmt.Sprintf("obj.%s = append(obj.%s, n(data))", n, n)
		} else if f.Type == "string" {
			return fmt.Sprintf("obj.%s = append(obj.%s, string(data))", n, n)
		}
	}

	return ""
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by legendsbrowser; DO NOT EDIT.
package model

import (
	"encoding/xml"
	"strconv"
	"github.com/iancoleman/strcase"
)

{{- range $name, $obj := .Objects }}
type {{ $obj.Name }} struct {
	{{- range $fname, $field := $obj.Fields }}
	{{- if not (and (eq $fname "type") (not (not $obj.SubTypes))) }}
	{{ $field.TypeLine $.Objects }}
	{{- end }}
	{{- end }}
	{{- if not (not $obj.SubTypes) }}
	Details any
	{{- end }}
}

{{- if $obj.Id }}
func (x *{{ $obj.Name }}) Id() int { return x.Id_ }
{{- end }}
{{- if $obj.Named }}
func (x *{{ $obj.Name }}) Name() string { return x.Name_ }
{{- end }}
{{- end }}

// Parser

func n(d []byte) int {
	v, _ := strconv.Atoi(string(d))
	return v
}

{{- range $name, $obj := .Objects }}
func parse{{ $obj.Name }}(d *xml.Decoder, start *xml.StartElement) (*{{ $obj.Name }}, error) {
	var (
		obj = {{ $obj.Name }}{}
		data []byte
	)
	for {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			{{- range $fname, $field := $obj.Fields }}
			case "{{ $fname }}":
				{{ $field.StartAction }}
			{{- end }}
			default:
				// fmt.Println("unknown field", t.Name.Local)
				d.Skip()
			}

		case xml.CharData:
			data = append(data, t...)

		case xml.EndElement:
			if t.Name.Local == start.Name.Local {
				return &obj, nil
			}

			switch t.Name.Local {
			{{- range $fname, $field := $obj.Fields }}
			case "{{ $fname }}":
				{{- if and (eq $fname "type") (not (not $obj.SubTypes)) }}
				var err error
				switch strcase.ToCamel(string(data)) {
				{{- range $sub := $obj.SubTypes }}
				case "{{ $sub }}":
					obj.Details, err = parse{{ $obj.Name }}{{ $sub }}(d, start)
				{{- end }}
				default:
					d.Skip()
				}
				if err != nil {
					return nil, err
				}
				return &obj, nil
				{{- else }}
				{{ $field.EndAction }}
				{{- end }}
			{{- end }}
			default:
				// fmt.Println("unknown field", t.Name.Local)
			}
		}
	}
}
{{- end }}
`))
