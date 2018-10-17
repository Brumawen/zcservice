package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// GetRequest holds the search criteria to be used to search for services
type GetRequest struct {
	ServiceType string `json:"serviceType"` // The search service type
	Domain      string `json:"domain"`      // The search domain.  For local networks, default of "local" is fine.
	WaitTime    int    `json:"waitTime"`    // The maximum amount of time to wait for a response
}

// CreateResponse creates a response from this request
func (e *GetRequest) CreateResponse() GetResponse {
	return GetResponse{
		ServiceType: e.ServiceType,
		Domain:      e.Domain,
	}
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (e *GetRequest) ReadFrom(r io.ReadCloser) error {
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
func (e *GetRequest) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (e *GetRequest) Serialize() (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (e *GetRequest) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &e)
	e.SetDefaults()
	return err
}

// SetDefaults checks the values and sets the defaults
func (e *GetRequest) SetDefaults() {
	if e.Domain == "" {
		e.Domain = "local"
	}
}
