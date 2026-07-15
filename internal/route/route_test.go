package route

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tmail/config"
	"tmail/internal/api"

	"github.com/labstack/echo/v5"
	"github.com/sunls24/gox/server"
)

func TestRegisterTurnstileRoutes(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		e := echo.New()
		Register(e, &config.Config{})

		assertRoute(t, e, http.MethodGet, "/api/turnstile/status", true)
		assertRoute(t, e, http.MethodPost, "/api/turnstile/verify", false)
	})

	t.Run("enabled", func(t *testing.T) {
		e := echo.New()
		Register(e, &config.Config{
			TurnstileSiteKey:   "test-site-key",
			TurnstileSecretKey: "test-secret-key",
		})

		assertRoute(t, e, http.MethodGet, "/api/turnstile/status", true)
		assertRoute(t, e, http.MethodPost, "/api/turnstile/verify", true)
	})
}

func TestTurnstileProtection(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		cfg := &config.Config{DomainList: []string{"example.com"}}
		srv := newTestServer(cfg)

		status := request(t, srv, http.MethodGet, "/api/turnstile/status")
		assertStatus(t, status, http.StatusOK)
		var body struct {
			Data struct {
				Enabled bool `json:"enabled"`
			} `json:"data"`
		}
		if err := json.Unmarshal(status.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Data.Enabled {
			t.Fatal("expected Turnstile status to be disabled")
		}

		domain := request(t, srv, http.MethodGet, "/api/domain")
		assertStatus(t, domain, http.StatusOK)
	})

	t.Run("enabled", func(t *testing.T) {
		cfg := &config.Config{
			TurnstileSiteKey:   "test-site-key",
			TurnstileSecretKey: "test-secret-key",
		}
		srv := newTestServer(cfg)

		status := request(t, srv, http.MethodGet, "/api/turnstile/status")
		assertStatus(t, status, http.StatusOK)
		var body struct {
			Data struct {
				Enabled  bool `json:"enabled"`
				Verified bool `json:"verified"`
			} `json:"data"`
		}
		if err := json.Unmarshal(status.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if !body.Data.Enabled || body.Data.Verified {
			t.Fatalf("unexpected Turnstile status: %+v", body.Data)
		}

		domain := request(t, srv, http.MethodGet, "/api/domain")
		assertStatus(t, domain, http.StatusUnauthorized)
	})
}

func assertRoute(t *testing.T, e *echo.Echo, method, path string, exists bool) {
	t.Helper()
	_, err := e.Router().Routes().FindByMethodPath(method, path)
	if (err == nil) != exists {
		t.Fatalf("route %s %s existence = %t, want %t", method, path, err == nil, exists)
	}
}

func newTestServer(cfg *config.Config) *server.Server {
	return server.New(func(srv *server.Server) {
		api.ServerContext(srv, cfg, nil)
		Register(srv.Echo, cfg)
	})
}

func request(t *testing.T, srv *server.Server, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	srv.Echo.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder
}

func assertStatus(t *testing.T, recorder *httptest.ResponseRecorder, want int) {
	t.Helper()
	if recorder.Code != want {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, want, recorder.Body.String())
	}
}
