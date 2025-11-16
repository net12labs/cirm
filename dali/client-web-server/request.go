package clientwebserver

import (
	webserver "github.com/net12labs/cirm/dali/web-server"
)

type Request struct {
	*webserver.Request
	Response *Response
}
type Response struct {
	*webserver.Response
}
