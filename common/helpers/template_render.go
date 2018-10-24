package helpers

import (
	"bytes"
	"fmt"
	"text/template"
)

func Render(tpl []byte, obj interface{}) (string, error) {
	goTemplate, err := template.New(
		"",
	).Parse(string(tpl))
	if err != nil {
		return "", fmt.Errorf("error creating Go template: %s", err)
	}
	var templateBuffer bytes.Buffer
	err = goTemplate.Execute(&templateBuffer, obj)
	if err != nil {
		return "", fmt.Errorf(
			"error rendering Go template into ARM template: %s",
			err,
		)
	}
	return templateBuffer.String(), nil
}
