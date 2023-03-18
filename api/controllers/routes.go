package controllers

import (
	"github.com/mertture/audit-log/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route	
	s.Router.POST("/api/users/register", middlewares.SetMiddlewareJSON(s.Register))
	s.Router.POST("/api/users/login", middlewares.SetMiddlewareJSON(s.Login))

	
	s.Router.POST("/api/logs", middlewares.SetMiddlewareAuthentication(s.CreateEvent))
	s.Router.GET("/api/logs", middlewares.SetMiddlewareAuthentication(s.GetAllEvents))
	s.Router.GET("/api/logs/:id", middlewares.SetMiddlewareAuthentication(s.GetEventByID))
	s.Router.DELETE("/api/logs/:id", middlewares.SetMiddlewareAuthentication(s.DeleteEvent))
	






}
