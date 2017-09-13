package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValidationError struct {
	error
	expected string

	received string
	raw      string
}

func TotalValidator(total int) Validator {
	return func(b []byte, resp *http.Response) error {
		t := struct{ Total int }{Total: -1}
		err := json.Unmarshal(b, &t)
		if err != nil {
			return err
		}
		if t.Total != total {
			return fmt.Errorf("Received total (%d) doesn't match expected (%d)", t.Total, total)
		}
		return nil
	}
}

func RandValidator(r *int) Validator {
	return func(b []byte, resp *http.Response) error {
		t := struct{ Total int }{Total: -1}
		err := json.Unmarshal(b, &t)
		if err != nil {
			return err
		}
		*r = t.Total
		return nil
	}
}
