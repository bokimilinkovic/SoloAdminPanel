package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getRequestBody(ctx echo.Context) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, err
	}

	ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

func bindAndValidate(request interface{}, ctx echo.Context) error {
	if err := ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}
