package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Config defines the configuration for the web server
type Config struct {
	ID                 string `json:"id"`                 // ID of the ZeroConf microservice
	Name               string `json:"name"`               // Name of the service
	DefaultServiceType string `json:"defaultServiceType"` // Default Service Type to use
}

// ReadFromFile will read the configuration settings from the specified file
func (c *Config) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteToFile will write the configuration settings to the specified file
func (c *Config) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (c *Config) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteTo serializes the entity and writes it to the http response
func (c *Config) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// SetDefaults checks the values and sets the defaults
func (c *Config) SetDefaults() {
	mustSave := false
	if c.ID == "" {
		// Generate a new GUID for this ID
		if uuid, err := uuid.NewV4(); err == nil {
			c.ID = strings.Replace(uuid.String(), "-", "", -1)
			mustSave = true
		}
	}
	if c.Name == "" {
		c.Name = "ZCService"
	}
	if c.DefaultServiceType == "" {
		c.DefaultServiceType = "_zcservice._tcp"
		mustSave = true
	}
	// Todo
	if mustSave {
		c.WriteToFile("config.json")
	}
}
