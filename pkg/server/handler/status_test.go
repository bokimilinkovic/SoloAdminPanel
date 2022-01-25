package handler_test

import (
	"net/http"
	"testing"

	"github.com/kolosek/pkg/server/handler"
	"github.com/stretchr/testify/assert"
)

func TestStatusH_CheckStatus(t *testing.T) {
	ctx, rec := setupEchoServer(t, "", "")
	statusHandler := handler.NewStatusH()
	err := statusHandler.CheckStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, rec.Body.String(), "Status is OK, everything works!")
	assert.Equal(t, http.StatusOK, rec.Code)
}
