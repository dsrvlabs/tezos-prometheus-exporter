package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	c := e.NewContext(req, rr)

	err := health(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)

	rawBody, err := ioutil.ReadAll(rr.Body)

	assert.Nil(t, err)
	assert.Equal(t, "OK", string(rawBody))
}
