package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// GetResponse holds the response data for a GetRequest call
type GetResponse struct {
	ServiceType string        `json:"serviceType"` // The service type
	Domain      string        `json:"domain"`      // The domain
	Services    []ServiceItem `json:"services"`    // The list of services
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (e *GetResponse) ReadFrom(r io.ReadCloser) error {
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
func (e *GetResponse) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (e *GetResponse) Serialize() (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (e *GetResponse) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &e)
	e.SetDefaults()
	return err
}

// SetDefaults checks the values and sets the defaults
func (e *GetResponse) SetDefaults() {
}
