package common

import (
	"encoding/base64"
	"fmt"
	"os"
)

type Context interface {
	Get(key string, dflt interface{}) interface{}
	GetString(key, dflt string) (string, error)
	GetBytes(key string, dflt []byte) ([]byte, error)
	GetFile(key string, dflt *os.File) (*os.File, error)
	GetRawMap() map[string]interface{}
	Set(string, interface{})
}

type context struct {
	Data map[string]interface{} `json:"data,omitempty"`
}

func NewContext() Context {
	return &context{
		Data: map[string]interface{}{},
	}
}

func (c *context) Get(key string, dflt interface{}) interface{} {
	val, ok := c.Data[key]
	if !ok {
		return dflt
	}
	return val
}

func (c *context) GetString(key, dflt string) (string, error) {
	iVal, ok := c.Data[key]
	if !ok {
		return dflt, nil
	}
	val, ok := iVal.(string)
	if !ok {
		return "",
			fmt.Errorf(
				`context key "%s" did not reference a string as expected`,
				key,
			)
	}
	return val, nil
}

func (c *context) GetBytes(key string, dflt []byte) ([]byte, error) {
	iVal, ok := c.Data[key]
	if !ok {
		return dflt, nil
	}
	valStr, ok := iVal.(string)
	if !ok {
		return nil,
			fmt.Errorf(
				`context key "%s" did not reference base64 encoded bytes as expected`,
				key,
			)
	}
	val, err := base64.StdEncoding.DecodeString(valStr)
	if err != nil {
		return nil,
			fmt.Errorf(
				`error base64 decoding string referenced by context key "%s": %s`,
				key,
				err,
			)
	}
	return val, nil
}

func (c *context) GetFile(key string, dflt *os.File) (*os.File, error) {
	filename, err := c.GetString(key, "")
	if err != nil {
		return nil, err
	}
	if filename == "" {
		return dflt, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf(`error opening file "%s": %s`, filename, err)
	}
	return file, nil
}

func (c *context) GetRawMap() map[string]interface{} {
	return c.Data
}

func (c *context) Set(key string, val interface{}) {
	c.Data[key] = val
}
