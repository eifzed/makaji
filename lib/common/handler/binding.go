package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/schema"
)

const (
	ContentURLEncoded string = "application/x-www-form-urlencoded"
	ContentJSON       string = "application/json"
	ContentType       string = "Content-Type"
)

type routeContext struct{}

func Bind(r *http.Request, param interface{}, ignoreInvalidKeys ...bool) error {
	if r.Method == http.MethodGet {
		err := decodeSchemaRequest(r, param, ignoreInvalidKeys...)
		if err != nil {
			return err
		}
		return nil
	}

	contentType := filterContentType(r.Header.Get(ContentType))
	switch contentType {
	case ContentURLEncoded:
		err := decodeSchemaRequest(r, param, ignoreInvalidKeys...)
		if err != nil {
			return err
		}
	case ContentJSON:
		decoder := json.NewDecoder(r.Body)
		if len(ignoreInvalidKeys) > 0 {
			if !ignoreInvalidKeys[0] {
				decoder.DisallowUnknownFields()
			}
		}
		err := decoder.Decode(param)
		if err != nil {
			return errors.New("error binding")
		}
		// TODO: add json validator

	}
	return nil
}

func filterContentType(contentType string) string {
	for i, char := range contentType {
		if char == ' ' || char == ';' {
			return contentType[:i]
		}
	}
	return contentType
}

func decodeSchemaRequest(r *http.Request, param interface{}, ignoreInvalidKeys ...bool) error {
	sourceDecode := r.URL.Query()
	decoder := schema.NewDecoder()
	ignoreKeys := true
	if len(ignoreInvalidKeys) > 0 {
		ignoreKeys = ignoreInvalidKeys[0]
	}
	decoder.IgnoreUnknownKeys(ignoreKeys)
	err := decoder.Decode(param, sourceDecode)
	// TODO: add url schema validator
	if err != nil {
		// TODO: return 403
		return errors.New("error binding")
	}
	return nil
}
