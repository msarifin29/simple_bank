package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/msarifin29/simple_bank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	server.setUpRoute()
	return server
}

func (server *Server) setUpRoute() {
	router := gin.Default()
	router.POST(`/api/accounts`, server.createAccount)
	router.GET(`/api/accounts/:id`, server.getAccount)
	router.GET(`/api/accounts/`, server.listAccount)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
