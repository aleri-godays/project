package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//Health just returns status: ok as json
func (s *Server) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.config.Version,
	})
}
