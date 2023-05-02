package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/internal/app/config"
	"testing"
)

func TestApi_HandleTest(t *testing.T) {
	s := New(config.NewConfig())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	s.handleTest().ServeHTTP(rec, req)
	assert.Equal(t, rec.Body.String(), "Just test")
}
