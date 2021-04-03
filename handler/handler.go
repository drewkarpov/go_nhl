package handler

import "net/http"

type IHandler interface {
	LoggingRequest(r http.Request, statusCode int)
}
