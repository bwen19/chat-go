package api

import "net/http"

type WebsocketHandler struct {
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}

func (h *WebsocketHandler) ServeHTTP(http.ResponseWriter, *http.Request) {

}
