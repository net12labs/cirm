package web_server

import (
	"encoding/json"
	"fmt"
	"io"

	webserver "github.com/net12labs/cirm/mali/web-server"
)

type Request struct {
	*webserver.Request
	Response *Response
}
type Response struct {
	*webserver.Response
}

func (s *Request) Validate_HasBody() bool {
	return s.Request.Req.Body != nil
}

func (s *Request) ReadRequestBodyAsMap() (map[string]any, error) {

	var result map[string]any
	bodyBytes, err := io.ReadAll(s.Request.Req.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Request) WriteResponse(v any) error {
	fmt.Println("Writing response:", v)
	return s.Request.WriteResponse(v)
}
