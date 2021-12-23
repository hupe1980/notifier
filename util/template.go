package util

import (
	"bytes"
	"text/template"
)

func Filter(s []string, e string) bool {
	if len(s) == 0 {
		return true
	}

	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func ExecuteTemplate(id, tmpl, message string, extras map[string]string) (string, error) {
	if tmpl == "" {
		tmpl = "{{ .Message }}"
	}

	t, err := template.New(id).Parse(tmpl)
	if err != nil {
		return "", err
	}

	extras["Message"] = message

	var buf bytes.Buffer
	if err := t.Execute(&buf, extras); err != nil {
		return "", err
	}

	return buf.String(), nil
}
