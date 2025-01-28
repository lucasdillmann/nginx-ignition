package server

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type state struct {
	engine   *gin.Engine
	server   *http.Server
	listener *net.Listener
}

func newState(engine *gin.Engine) *state {
	return &state{engine: engine}
}
