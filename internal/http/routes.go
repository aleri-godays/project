package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//AddRoutes creates all routes
func (s *Server) AddRoutes(handler projectHandler) {
	s.e.GET("/health", s.Health)

	api := s.e.Group("/api/v1")

	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(s.config.JWTSecret),
	}))

	api.GET("/project/:id", handler.Project)
	api.GET("/project", handler.Projects)

	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
