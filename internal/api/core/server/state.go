package server

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type state struct {
	engine   *gin.Engine
	server   *http.Server
	listener *net.Listener
}

func newState(engine *gin.Engine) *state {
	return &state{engine: engine}
}
