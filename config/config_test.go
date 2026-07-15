package config

import "testing"

func TestMustNewTurnstileDisabled(t *testing.T) {
	t.Setenv("TURNSTILE_SITE_KEY", "")
	t.Setenv("TURNSTILE_SECRET_KEY", "")
	t.Setenv("TURNSTILE_COOKIE_TTL", "0s")

	cfg := MustNew()
	if cfg.TurnstileEnabled() {
		t.Fatal("expected Turnstile to be disabled")
	}
}

func TestMustNewTurnstileEnabled(t *testing.T) {
	t.Setenv("TURNSTILE_SITE_KEY", "test-site-key")
	t.Setenv("TURNSTILE_SECRET_KEY", "test-secret-key")
	t.Setenv("TURNSTILE_COOKIE_TTL", "1h")

	cfg := MustNew()
	if !cfg.TurnstileEnabled() {
		t.Fatal("expected Turnstile to be enabled")
	}
}

func TestMustNewRejectsPartialTurnstileConfig(t *testing.T) {
	for _, tc := range []struct {
		name      string
		siteKey   string
		secretKey string
	}{
		{name: "missing secret key", siteKey: "test-site-key"},
		{name: "missing site key", secretKey: "test-secret-key"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TURNSTILE_SITE_KEY", tc.siteKey)
			t.Setenv("TURNSTILE_SECRET_KEY", tc.secretKey)
			assertPanics(t, MustNew)
		})
	}
}

func TestMustNewRejectsInvalidTurnstileTTL(t *testing.T) {
	t.Setenv("TURNSTILE_SITE_KEY", "test-site-key")
	t.Setenv("TURNSTILE_SECRET_KEY", "test-secret-key")
	t.Setenv("TURNSTILE_COOKIE_TTL", "0s")

	assertPanics(t, MustNew)
}

func assertPanics(t *testing.T, fn func() *Config) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected function to panic")
		}
	}()
	fn()
}
