package context

import (
	"github.com/labstack/echo/v4"
)

// Context represents custom package level context which extends `echo.Context`.
type Context struct {
	echo.Context
}

// New creates a new custom context.
func New(ectx echo.Context) *Context {
	return &Context{Context: ectx}
}

// StdLogger returns standard logger instance.
func (c *Context) StdLogger() *StdLogger {
	logger, ok := c.Logger().(*StdLogger)
	if !ok {
		return nil
	}

	return logger
}
