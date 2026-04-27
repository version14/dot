package render

import (
	"bytes"
	"fmt"
	"text/template"
)

// Render parses tmpl as a Go text/template and executes it against data.
// Strict mode: missing keys raise an error rather than silently producing "<no value>".
//
// Generators use this to materialize .tmpl files with answers from the spec:
//
//	out, err := render.Render(raw, map[string]any{"ProjectName": "myapp"})
func Render(tmpl string, data interface{}) ([]byte, error) {
	t, err := template.New("dot").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("render: parse: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render: execute: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderString is the string-typed convenience wrapper around Render.
func RenderString(tmpl string, data interface{}) (string, error) {
	out, err := Render(tmpl, data)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// RenderInto renders tmpl and writes the result to dst. Useful when the caller
// wants to stream rather than allocate.
func RenderInto(dst *bytes.Buffer, tmpl string, data interface{}) error {
	t, err := template.New("dot").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("render: parse: %w", err)
	}
	if err := t.Execute(dst, data); err != nil {
		return fmt.Errorf("render: execute: %w", err)
	}
	return nil
}
