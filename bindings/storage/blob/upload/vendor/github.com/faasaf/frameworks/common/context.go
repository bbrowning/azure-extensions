package common

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

type Context interface {
	Get(key string, dflt interface{}) interface{}
	GetString(key, dflt string) (string, error)
	GetBytes(key string, dflt []byte) ([]byte, error)
	GetFile(key string, dflt *os.File) (*os.File, error)
	GetHTTPRequest(key string) (*http.Request, error)
	GetRawMap() map[string]interface{}
	SetHTTPRequest(string, *http.Request) error
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

func (c *context) GetHTTPRequest(key string) (*http.Request, error) {
	file, err := c.GetFile(key, nil)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, nil
	}
	return http.ReadRequest(bufio.NewReader(file))
}

func (c *context) GetRawMap() map[string]interface{} {
	return c.Data
}

func (c *context) SetHTTPRequest(key string, req *http.Request) error {
	reqBites, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(reqBites); err != nil {
		return err
	}

	c.Data[key] = file.Name()

	return nil
}

func (c *context) Set(key string, val interface{}) {
	c.Data[key] = val
}
