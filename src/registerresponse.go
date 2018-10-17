package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// RegisterResponse holds the response data for a RegisterRequest call
type RegisterResponse struct {
	ID string `json:"id"` // ID of the service registration
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (e *RegisterResponse) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			err = json.Unmarshal(b, &e)
		}
	}
	e.SetDefaults()
	return err
}

// WriteTo serializes the entity and writes it to the http response
func (e *RegisterResponse) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (e *RegisterResponse) Serialize() (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (e *RegisterResponse) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &e)
	e.SetDefaults()
	return err
}

// SetDefaults checks the values and sets the defaults
func (e *RegisterResponse) SetDefaults() {
}
