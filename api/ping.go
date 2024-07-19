package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerPing(r *gin.Engine) {
	r.GET("/ping", pong)
}

func pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
