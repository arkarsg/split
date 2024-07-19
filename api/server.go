package api

import (
	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for debt service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

var r *gin.Engine

func (s *Server) initRoutes() {
	s.registerPing(r)
	s.registerUsers(r)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func NewServer(store *db.Store) *Server {
	r = gin.Default()
	server := &Server{
		store:  store,
		router: r,
	}
	server.initRoutes()
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
