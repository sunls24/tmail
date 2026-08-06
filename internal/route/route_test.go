package route

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	"tmail/config"
	"tmail/internal/api"

	"github.com/labstack/echo/v5"
	"github.com/sunls24/gox/server"
)

func TestRegisterTurnstileRoutes(t *testing.T) {
	t.Run("disabled", func(t *testing.T) {
		e := echo.New()
		Register(e, &config.Config{})

		assertRoute(t, e, http.MethodPost, "/api/turnstile/verify", false)
	})

	t.Run("enabled", func(t *testing.T) {
		e := echo.New()
		Register(e, &config.Config{
			TurnstileSiteKey:   "test-site-key",
			TurnstileSecretKey: "test-secret-key",
		})

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

	t.Run("api key", func(t *testing.T) {
		cfg := &config.Config{
			APIKey:             "test-api-key",
			TurnstileSiteKey:   "test-site-key",
			TurnstileSecretKey: "test-secret-key",
		}
		srv := newTestServer(cfg)

		req := httptest.NewRequest(http.MethodGet, "/api/domain", nil)
		req.Header.Set("X-API-Key", "test-api-key")
		recorder := httptest.NewRecorder()
		srv.Echo.ServeHTTP(recorder, req)
		assertStatus(t, recorder, http.StatusOK)
	})
}

func TestReportHMAC(t *testing.T) {
	cfg := &config.Config{ReportHMACSecret: "test-secret", ReportMaxBodySize: 1 << 20}
	srv := newTestServer(cfg)

	t.Run("rejects unsigned request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/report?to=external@example.com", nil)
		srv.Echo.ServeHTTP(recorder, req)
		assertStatus(t, recorder, http.StatusUnauthorized)
	})

	t.Run("allows valid request", func(t *testing.T) {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		path := "/api/report"
		to := "external@example.com"
		message := timestamp + "\nPOST\n" + path + "\n" + to
		mac := hmac.New(sha256.New, []byte(cfg.ReportHMACSecret))
		_, _ = mac.Write([]byte(message))

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			path+"?to="+to,
			strings.NewReader("Subject: test\r\n\r\n"),
		)
		req.Header.Set("X-Tmail-Timestamp", timestamp)
		req.Header.Set("X-Tmail-Signature", hex.EncodeToString(mac.Sum(nil)))
		srv.Echo.ServeHTTP(recorder, req)
		assertStatus(t, recorder, http.StatusOK)
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
