package response

import (
	"encoding/json"
	"net/http"
)

type Response map[string]interface{}

func New(status bool, message string) Response {
	m := make(Response)
	m["status"] = status
	m["message"] = message
	return m
}

func (m Response) SetResponse(status bool, message string) {
	m["status"] = status
	m["message"] = message
}

func (m Response) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}

func (m Response) SetString(key string, value string) {
	m[key] = value
}

func (m Response) Respond(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

