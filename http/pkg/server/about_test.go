package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestAbout(t *testing.T) {
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	About(wr, req)
	assert.Equal(t, wr.Code, http.StatusOK)
}
