package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kurniawanew/pfm-kambing/api/middlewares"
)

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.GET("/", s.Home)

	// Login Route
	s.Router.POST("/login", s.Login)

	authorized := s.Router.Group("/")
	authorized.Use(gin.Logger())
	authorized.Use(gin.Recovery())
	authorized.Use(middlewares.SetMiddlewareAuthentication())
	{
		authorized.GET("/users", s.GetUsers)
		authorized.GET("/users/:id", s.GetUser)
		authorized.PUT("/users/:id", s.UpdateUser)
		authorized.DELETE("/users/:id", s.DeleteUser)
		authorized.POST("/users/create", s.CreateUser)

		authorized.GET("/transactions", s.GetTransactions)
		authorized.GET("/transactions/:id", s.GetTransaction)
		authorized.PUT("/transactions/:id", s.UpdateTransaction)
		authorized.DELETE("/transactions/:id", s.DeleteTransaction)
		authorized.POST("/transactions/create", s.CreateTransaction)
	}
}
