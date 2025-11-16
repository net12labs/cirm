package apiwebserver

import (
	"encoding/json"
	"io"

	webserver "github.com/net12labs/cirm/dali/web-server"
)

type ApiRequest struct {
	*webserver.Request
	Response *ApiResponse
}
type ApiResponse struct {
	*webserver.Response
}

func (s *ApiRequest) Validate_HasBody() bool {
	return s.Request.Req.Body != nil
}

func (s *ApiRequest) ReadRequestBodyAsMap() (map[string]any, error) {

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

func (s *ApiRequest) WriteResponse(v any) error {
	return s.Request.WriteResponse(v)
}
