package main

import (
	"github.com/M0hammadUsman/snippetbox/internal/assert"
	"net/http"
	"testing"
)

// End to End
func TestPing(t *testing.T) {
	// t.Parallel() -> go test -race ./...               | -race slows down a bit
	app := application{}
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	statusCode, _, body := ts.get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}
