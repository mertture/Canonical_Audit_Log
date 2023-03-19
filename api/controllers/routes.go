package controllers

import (
	"github.com/mertture/audit-log/api/middlewares"
)

func (s *Server) initializeRoutes() {

	// Home Route	
	s.Router.POST("/api/users/register", middlewares.SetMiddlewareJSON(s.Register))
	s.Router.POST("/api/users/login", middlewares.SetMiddlewareJSON(s.Login))

	
	s.Router.POST("/api/events", middlewares.SetMiddlewareAuthentication(s.CreateEvent))
	s.Router.GET("/api/events", middlewares.SetMiddlewareAuthentication(s.GetAllEvents))
	s.Router.GET("/api/events/:id", middlewares.SetMiddlewareAuthentication(s.GetEventByID))
	s.Router.DELETE("/api/events/:id", middlewares.SetMiddlewareAuthentication(s.DeleteEvent))

	s.Router.GET("/api/events/types/:id", middlewares.SetMiddlewareAuthentication(s.GetEventByTypeID))
	






}
