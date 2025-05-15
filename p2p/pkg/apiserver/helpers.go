package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// StatusResponse is a helper struct used to wrap a string message with a status code.
type StatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteResponse helper is used to write a response to the client with the given code
// and message. If the message is a string, it will be wrapped in a StatusResponse
// struct. Otherwise, the message will be encoded as JSON.
func WriteResponse(w http.ResponseWriter, code int, message any) error {
	var b bytes.Buffer
	switch message := message.(type) {
	case string:
		err := json.NewEncoder(&b).Encode(StatusResponse{Code: code, Message: message})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to encode status response: %w", err)
		}
	default:
		err := json.NewEncoder(&b).Encode(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return fmt.Errorf("failed to encode response: %w", err)
		}
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprintln(w, b.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}

func WriteError(w http.ResponseWriter, code int, err error) {
	_ = WriteResponse(w, code, err.Error())
}

// MethodHandler helper is used to wrap a handler and ensure that the request method
// matches the given method. If the method does not match, a 405 is returned.
func MethodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			err := WriteResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		handler(w, r)
	}
}

// BindJSON helper is used to bind the request body to the given type.
func BindJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var body T

	if r.Body == nil {
		return body, errors.New("no body")
	}
	//nolint:errcheck
	defer r.Body.Close()

	return body, json.NewDecoder(r.Body).Decode(&body)
}
