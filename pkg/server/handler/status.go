package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StatusH manages status endpoints
type StatusH struct{}

// NewStatusH creates a new status handler
func NewStatusH() *StatusH { return &StatusH{} }

// CheckStatus godoc
// @Summary Checks if server is running
// @Description Checks status
// @Success 200 {string} string	"Status is OK, everything works!"
// @Router /status [get]
func (sh *StatusH) CheckStatus(c echo.Context) error {
	return c.String(http.StatusOK, "Status is OK, everything works!")
}
