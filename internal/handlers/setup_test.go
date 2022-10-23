package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/ironedr/gochat/internal/application"
)

type testServer struct {
	*httptest.Server
}

func createTestApp(t *testing.T) *application.Application {
	templateCache, err := NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	return &application.Application{
		TemplateCache:  templateCache,
		Logger:         log.New(io.Discard, "", 0),
		SessionManager: scs.New(),
		Env:            "test",
		DSN:            "postgres://myuser:test123@localhost:5432/gochat-test",
		BaseURL:        "localhost/",
		CSRFKey:        []byte("mySuperSecret"),
		Port:           9000,
	}
}

func createTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = cookieJar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
