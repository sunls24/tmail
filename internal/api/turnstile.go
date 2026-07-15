package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"tmail/config"

	"github.com/sunls24/gox/server"
)

const (
	turnstileCookieName = "tmail_turnstile"
	turnstileVerifyURL  = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

type TurnstileStatusResp struct {
	Enabled  bool   `json:"enabled"`
	Verified bool   `json:"verified"`
	SiteKey  string `json:"site_key"`
}

func TurnstileStatus(ctx context.Context) (*TurnstileStatusResp, error) {
	cfg := Config(ctx)
	enabled := cfg.TurnstileEnabled()
	return &TurnstileStatusResp{
		Enabled:  enabled,
		Verified: enabled && TurnstileVerified(server.EchoContext(ctx).Request(), cfg),
		SiteKey:  cfg.TurnstileSiteKey,
	}, nil
}

type ReqTurnstileVerify struct {
	Token string `json:"token"`
}

func TurnstileVerify(ctx context.Context, req ReqTurnstileVerify) error {
	if req.Token == "" {
		return server.BadParam()
	}

	cfg := Config(ctx)
	form := url.Values{
		"secret":   {cfg.TurnstileSecretKey},
		"response": {req.Token},
	}
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		turnstileVerifyURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("turnstile siteverify returned %s", resp.Status)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if !result.Success {
		return server.ErrMsg("人机验证失败").WithStatusCode(http.StatusUnauthorized)
	}

	ec := server.EchoContext(ctx)
	ec.SetCookie(newTurnstileCookie(cfg, time.Now()))
	return nil
}

func TurnstileVerified(req *http.Request, cfg *config.Config) bool {
	cookie, err := req.Cookie(turnstileCookieName)
	if err != nil {
		return false
	}

	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || time.Now().Unix() >= expires {
		return false
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	return hmac.Equal(signature, signTurnstileCookie(parts[0], cfg.TurnstileSecretKey))
}

func newTurnstileCookie(cfg *config.Config, now time.Time) *http.Cookie {
	expires := now.Add(cfg.TurnstileCookieTTL)
	payload := strconv.FormatInt(expires.Unix(), 10)
	signature := base64.RawURLEncoding.EncodeToString(signTurnstileCookie(payload, cfg.TurnstileSecretKey))
	return &http.Cookie{
		Name:     turnstileCookieName,
		Value:    payload + "." + signature,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(cfg.TurnstileCookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   !cfg.Debug,
		SameSite: http.SameSiteLaxMode,
	}
}

func signTurnstileCookie(payload, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}
