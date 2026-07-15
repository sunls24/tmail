package api

import (
	"net/http/httptest"
	"testing"
	"time"
	"tmail/config"
)

func TestTurnstileVerified(t *testing.T) {
	cfg := &config.Config{
		TurnstileSecretKey: "test-secret",
		TurnstileCookieTTL: time.Hour,
		Debug:              true,
	}
	req := httptest.NewRequest("GET", "/api/domain", nil)
	req.AddCookie(newTurnstileCookie(cfg, time.Now()))

	if !TurnstileVerified(req, cfg) {
		t.Fatal("expected cookie to be valid")
	}
}

func TestTurnstileVerifiedRejectsTamperedCookie(t *testing.T) {
	cfg := &config.Config{
		TurnstileSecretKey: "test-secret",
		TurnstileCookieTTL: time.Hour,
		Debug:              true,
	}
	cookie := newTurnstileCookie(cfg, time.Now())
	cookie.Value += "tampered"
	req := httptest.NewRequest("GET", "/api/domain", nil)
	req.AddCookie(cookie)

	if TurnstileVerified(req, cfg) {
		t.Fatal("expected tampered cookie to be rejected")
	}
}

func TestTurnstileVerifiedRejectsExpiredCookie(t *testing.T) {
	cfg := &config.Config{
		TurnstileSecretKey: "test-secret",
		TurnstileCookieTTL: time.Hour,
		Debug:              true,
	}
	req := httptest.NewRequest("GET", "/api/domain", nil)
	req.AddCookie(newTurnstileCookie(cfg, time.Now().Add(-2*time.Hour)))

	if TurnstileVerified(req, cfg) {
		t.Fatal("expected expired cookie to be rejected")
	}
}

func TestNewTurnstileCookieTTL(t *testing.T) {
	cfg := &config.Config{
		TurnstileSecretKey: "test-secret",
		TurnstileCookieTTL: 6 * time.Hour,
		Debug:              true,
	}
	now := time.Unix(1_700_000_000, 0)
	cookie := newTurnstileCookie(cfg, now)

	if cookie.MaxAge != 6*60*60 {
		t.Fatalf("unexpected max age: %d", cookie.MaxAge)
	}
	if !cookie.Expires.Equal(now.Add(6 * time.Hour)) {
		t.Fatalf("unexpected expires: %s", cookie.Expires)
	}
}
