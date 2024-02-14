package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/msarifin29/simple_bank/db/sqlc"
	"github.com/msarifin29/simple_bank/token"
	"github.com/msarifin29/simple_bank/util"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	config     util.Config
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	server.setUpRoute()
	return server, nil
}

func (server *Server) setUpRoute() {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST(`/api/users`, server.createUser)
	router.POST(`/api/users/login`, server.LoginUser)

	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker))

	authRoutes.POST(`/api/accounts`, server.createAccount)
	authRoutes.GET(`/api/accounts/:id`, server.getAccount)
	authRoutes.GET(`/api/accounts/`, server.listAccount)

	authRoutes.POST(`/api/transfer`, server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
