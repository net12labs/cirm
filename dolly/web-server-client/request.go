package clientwebserver

import (
	webserver "github.com/net12labs/cirm/mali/web-server"
)

type Request struct {
	*webserver.Request
	Response *Response
}
type Response struct {
	*webserver.Response
}
