package helper

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Error struct {
	Message string `json:"message"`
}

func ReadJSONRequest(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func JSONResponse(w http.ResponseWriter, code int, msg string, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := &Response{
		Message: msg,
		Data:    data,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err = w.Write(b); err != nil {
		return err
	}

	return nil
}

func JSONError(w http.ResponseWriter, code int, msg string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := &Error{
		Message: msg,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err = w.Write(b); err != nil {
		return err
	}

	return nil
}
