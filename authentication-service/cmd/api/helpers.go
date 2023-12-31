package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only single JSON value")
	}

	return nil
}


func (app *Config) witeJson(w http.ResponseWriter,data any,status int,headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key,val := range headers[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	_,err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error,status ...int) error {
	statusCode := http.StatusBadRequest
	
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.witeJson(w,payload,statusCode)
}