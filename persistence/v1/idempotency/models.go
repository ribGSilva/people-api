package idempotency

import (
	"bytes"
	"fmt"
	"text/template"
)

type Idempotency struct {
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
	Status   int    `json:"stats"`
	Response string `json:"resp"`
}

type search struct {
	Id       string
	Endpoint string
}

const idempotencyKey = "app.person.idempotency.{{.Id}}.endpoint.{{.Endpoint}}"

func buildKey(s search) (string, error) {
	t, err := template.New("idempotency").Parse(idempotencyKey)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, s); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return tpl.String(), nil
}
