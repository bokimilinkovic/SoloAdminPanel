package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kolosek/pkg/service/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func setupEchoServer(t *testing.T, queryParams, payload string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Debug = true
	validator, err := validation.NewValidator()
	require.NoError(t, err)

	e.Validator = validator
	e.Binder = &validation.Binder{}

	// Method POST can be used for every request, since it doesn't matter much for httptest.NewRequest which one it is.
	req := httptest.NewRequest(http.MethodPost, "http://localhost"+queryParams, strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	return ctx, rec
}
