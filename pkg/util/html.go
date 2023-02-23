package util

import (
	"bytes"
	"fmt"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var inputTypes = map[string]string{
	"string": "text",
	"int":    "number",
	"bool":   "checkbox",
}

func StructToHTML(model interface{}) string {
	t := reflect.ValueOf(model).Type()
	var buf bytes.Buffer
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Fprintf(
			&buf,
			`<div class="mb-3">
				<label for="%s" class="form-label">%s</label>
				<input type="%s" class="form-control bg-white text-dark" id="%s" name="%s" required>
			</div>`,
			field.Name,
			cases.Title(language.English).String(field.Name),
			inputTypes[field.Type.Name()],
			field.Name,
			field.Name,
		)
	}
	return buf.String()
}
