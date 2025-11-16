package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Req       *http.Request
	Resp      *http.ResponseWriter
	Path      *URLPath
	Method    string
	Response  *Response
	IsInvalid bool
	Error     string
}

type Response struct {
	StatusCode int
	Headers    http.Header
	MimeType   string
	req        *Request
}

func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	req := &Request{
		Req:       r,
		Resp:      &w,
		Path:      &URLPath{Path: r.URL.Path},
		Method:    r.Method,
		IsInvalid: false,
		Error:     "",
		Response: &Response{
			Headers: make(http.Header),
		},
	}
	req.Response.req = req
	return req
}

func (rs *Request) WriteResponse404() error {
	rs.Response.Headers.Set("Content-Type", "text/plain")
	rs.Response.Headers.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rs.Response.StatusCode = http.StatusNotFound
	(*rs.Resp).WriteHeader(rs.Response.StatusCode)
	(*rs.Resp).Write([]byte("404 Not Found"))
	return nil
}

func (req *Request) WriteResponseHTML(data any) error {
	if req.Response == nil {
		req.Response = &Response{}
	}
	req.Response.MimeType = "text/html"
	return req.WriteResponse(data)
}

func (r *Request) WriteResponse(data any) error {
	if r.Response == nil {
		r.Response = &Response{}
	}

	// Set Content-Type from MimeType if specified
	if r.Response.MimeType != "" {
		(*r.Resp).Header().Set("Content-Type", r.Response.MimeType)
	}

	// Set headers
	for key, values := range r.Response.Headers {
		for _, value := range values {
			(*r.Resp).Header().Add(key, value)
		}
	}

	// Write data
	if data != nil {
		// Set status code before writing data
		statusCode := r.Response.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		switch v := data.(type) {
		case string:
			(*r.Resp).WriteHeader(statusCode)
			_, err := (*r.Resp).Write([]byte(v))
			return err
		case []byte:
			(*r.Resp).WriteHeader(statusCode)
			_, err := (*r.Resp).Write(v)
			return err
		default:
			// Try to marshal as JSON
			jsonData, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to marshal data: %w", err)
			}
			(*r.Resp).WriteHeader(statusCode)
			// Set Content-Type to JSON only if MimeType not already set
			if r.Response.MimeType == "" {
				(*r.Resp).Header().Set("Content-Type", "application/json")
			}
			_, err = (*r.Resp).Write(jsonData)
			return err
		}
	} else {
		// No data, just write status code
		statusCode := r.Response.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		(*r.Resp).WriteHeader(statusCode)
	}

	return nil
}
