package http

import (
	"github.com/aleri-godays/project"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type projectHandler struct {
	repo project.Repository
}

func (h *projectHandler) Project(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return jsonError(c, "id must be of type int", http.StatusBadRequest)
	}

	p, err := h.repo.Get(c.Request().Context(), id)
	if err != nil {
		return jsonError(c, "could not fetch project", http.StatusInternalServerError)
	}

	return c.JSONPretty(http.StatusOK, p, "  ")
}

func (h *projectHandler) Projects(c echo.Context) error {
	ps, err := h.repo.All(c.Request().Context())
	if err != nil {
		return jsonError(c, "could not fetch projects", http.StatusInternalServerError)
	}

	return c.JSONPretty(http.StatusOK, ps, "  ")
}
