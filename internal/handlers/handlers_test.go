package handlers

import (
	"net/http"
	"testing"

	"github.com/ironedr/gochat/internal/assert"
)

func TestHome(t *testing.T) {

	app := createTestApp(t)
	handler := New(app)
	ts := createTestServer(t, handler.SetupRoutes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/")
	assert.Equal(t, code, http.StatusOK)

}
