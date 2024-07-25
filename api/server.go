package api

import (
	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for debt service
type Server struct {
	store  db.Store
	router *gin.Engine
}

var r *gin.Engine

func (s *Server) initRoutes() {
	s.registerPing(r)
	s.registerUsers(r)
	s.registerPayment(r)
	s.registerTransaction(r)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func NewServer(store db.Store) *Server {
	r = gin.Default()
	server := &Server{
		store:  store,
		router: r,
	}
	server.initRoutes()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("moneyamount", moneyAmount)
		v.RegisterValidation("supportedcurrency", validCurrency)
	}

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
