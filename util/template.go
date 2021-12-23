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

func ExecuteTemplate(id, format, message string, extras map[string]string) (string, error) {
	t, err := template.New(id).Parse(format)
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
