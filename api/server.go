package api

import (
	"fmt"

	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/arkarsg/splitapp/token"
	"github.com/arkarsg/splitapp/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for debt service
type Server struct {
	config     utils.ServerConfig
	store      db.Store
	router     *gin.Engine
	tokenMaker token.TokenMaker
}

var r *gin.Engine

func (s *Server) initRoutes() {
	s.registerPing(r)
	s.registerAccount(r)
	s.registerUser(r)
	s.registerTransaction(r)
	s.registerDebt(r)
	s.registerPayment(r)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func NewServer(config utils.ServerConfig, store db.Store) (*Server, error) {
	r = gin.Default()
	pasetoTokenMaker, err := token.NewPasetoMaker(config.Token.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}
	server := &Server{
		config:     utils.GetConfig(),
		store:      store,
		router:     r,
		tokenMaker: pasetoTokenMaker,
	}
	server.initRoutes()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("moneyamount", moneyAmount)
		v.RegisterValidation("supportedcurrency", validCurrency)
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
