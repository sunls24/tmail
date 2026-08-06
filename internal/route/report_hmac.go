package route

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
	"tmail/config"

	"github.com/labstack/echo/v5"
	"github.com/sunls24/gox/server"
)

const reportHMACClockSkew = 5 * time.Minute

func requireReportHMAC(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if cfg.ReportHMACSecret == "" {
				return next(c)
			}

			timestamp := c.Request().Header.Get("X-Tmail-Timestamp")
			signature := c.Request().Header.Get("X-Tmail-Signature")
			if !validReportHMAC(timestamp, signature, c, cfg.ReportHMACSecret) {
				return c.JSON(http.StatusUnauthorized, server.Envelope{
					Code:    -1,
					Message: "请求签名无效",
				})
			}
			return next(c)
		}
	}
}

func validReportHMAC(timestamp, signature string, c *echo.Context, secret string) bool {
	if timestamp == "" || signature == "" {
		return false
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	delta := time.Now().Unix() - ts
	if delta < -int64(reportHMACClockSkew/time.Second) || delta > int64(reportHMACClockSkew/time.Second) {
		return false
	}

	message := timestamp + "\n" + c.Request().Method + "\n" + c.Request().URL.Path + "\n" + c.QueryParam("to")
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))
	return signature == expected
}
