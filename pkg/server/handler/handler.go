package handler

import (
	"github.com/kolosek/pkg/server/context"
	"github.com/kolosek/pkg/service/log"
	"github.com/labstack/echo/v4"
)

// Handler represents a basic request handler foundation.
// Specific handlers should extend it if they need this functionality.
// By purpose this package does not contain New() function,
// since it is planned to be used as an extension for concrete handlers.
type Handler struct {
	Logger *log.Logger
}

// Context returns a custom context.
func (h *Handler) Context(ectx echo.Context) *context.Context {
	return ectx.(*context.Context)
}
