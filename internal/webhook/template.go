package webhook

import (
	"bytes"
	"encoding/json"
	"text/template"

	"email-forwarder/internal/email"
)

func FormatPayload(parsedEmail *email.ParsedEmail, tpl string) (*bytes.Buffer, error) {
	t, err := template.New("webhook").Funcs(template.FuncMap{
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
	}).Parse(tpl)
	if err != nil {
		return nil, err
	}

	var payload bytes.Buffer
	err = t.Execute(&payload, parsedEmail)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
