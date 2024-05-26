package commonwriter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eifzed/makaji/lib/common/commonerr"
	"github.com/pkg/errors"
)

type key struct{}

var errCtxKey key

type CustomError struct {
	ErrorMessage []*commonerr.ErrorFormat `json:"error_messages"`
}
type RequestInfo struct {
	StartRequest time.Time
	Host         string
	SourceIP     string
	RequestURL   string
	Method       string
	Error        error
	HTTPStatus   int
	UserAgent    string
}

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
	Links  *Links      `json:"links,omitempty"`
}

type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type Message struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	ErrorMessage []Message `json:"error_messages"`
}

func (resp *Response) Write(w http.ResponseWriter, r *http.Request, status int) (int, error) {
	if resp.Errors != nil {
		// if errors then data should be empty and vise versa
		resp.Data = nil
		// setError(r)
	}
	w.Header().Set("Content-Type", "application/json")
	respByte, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeLen, writeErr := w.Write([]byte(`{"error":["Internal Server Error"]}`))
		if writeErr != nil {
			return writeLen, writeErr
		}
		return writeLen, err
	}
	w.WriteHeader(status)
	return w.Write(respByte)
}

// TODO: enhance for loggin error
// func setError(r *http.Request) {
// 	ctx := r.Context()
// 	ctx = context.WithValue(ctx, errCtxKey, true)
// 	(*r) = *r.WithContext(ctx)
// }

func RespondOKWithData(ctx context.Context, w http.ResponseWriter, data interface{}) {
	sendResponseJSONData(w, nil, http.StatusOK, data)
}

func RespondOKWithMessage(ctx context.Context, w http.ResponseWriter, message string) {
	sendResponseJSONData(w, nil, http.StatusOK, &Message{Message: message})
}

func sendResponseJSONData(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, error) {
	resp := Response{Data: data}
	return resp.Write(w, r, status)
}

func sendResponseJSONError(w http.ResponseWriter, r *http.Request, status int, err *commonerr.ErrorMessage) (int, error) {
	resp := Response{Errors: err}
	return resp.Write(w, r, status)
}

func RespondOKWithByte(ctx context.Context, w http.ResponseWriter, byteData []byte) {
	w.Write(byteData)
}

func RespondOK(ctx context.Context, w http.ResponseWriter) {
	sendResponseJSONData(w, nil, http.StatusOK, &Message{Message: "OK"})
}

func RespondError(ctx context.Context, w http.ResponseWriter, errValue error) error {
	//TODO: handle uding custom error
	switch errCause := errors.Cause(errValue).(type) {
	case *commonerr.ErrorMessage:
		SetErrorFormat(ctx, w, errCause.Code, errCause)
	default:
		RespondDefaultError(ctx, w, errValue)
	}
	return nil
}

func RespondDefaultError(ctx context.Context, w http.ResponseWriter, errValue error) {
	SetErrorFormat(ctx, w, http.StatusInternalServerError, &commonerr.ErrorMessage{
		ErrorList: []*commonerr.ErrorFormat{{
			ErrorName:        "Error",
			ErrorDescription: errValue.Error(),
		}},
	})
}

func SetErrorFormat(ctx context.Context, w http.ResponseWriter, errorCode int, errMessage *commonerr.ErrorMessage) error {
	_, err := sendResponseJSONError(w, nil, errorCode, errMessage)
	return err
}
