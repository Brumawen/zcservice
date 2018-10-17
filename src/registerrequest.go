package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// RegisterRequest is the registration request data sent from a microservice
type RegisterRequest struct {
	ID          string   `json:"id"`          // ID of the service
	Name        string   `json:"name"`        // Name of the service
	PortNo      int      `json:"portNo"`      // Port number of the service
	ServiceType string   `json:"serviceType"` // Type of the server
	Domain      string   `json:"domain"`      // Service domain
	Text        []string `json:"text"`        // Additional service Text
}

// CreateResponse creates a response to the current request
func (e *RegisterRequest) CreateResponse() RegisterResponse {
	return RegisterResponse{
		ID: e.ID,
	}
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (e *RegisterRequest) ReadFrom(r io.ReadCloser) error {
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
func (e *RegisterRequest) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (e *RegisterRequest) Serialize() (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (e *RegisterRequest) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &e)
	e.SetDefaults()
	return err
}

// SetDefaults checks the values and sets the defaults
func (e *RegisterRequest) SetDefaults() {
}
