package route

import (
	"net/http"
	"tmail/config"
	"tmail/internal/api"

	"github.com/labstack/echo/v5"
	"github.com/sunls24/gox/server"
)

func Register(e *echo.Echo, cfg *config.Config) {
	g := e.Group("/api")
	g.POST("/report", server.WrapReplyResp(api.Report))
	g.GET("/turnstile/status", server.WrapResp(api.TurnstileStatus))

	protected := g
	if cfg.TurnstileEnabled() {
		g.POST("/turnstile/verify", server.WrapReq(api.TurnstileVerify))
		protected = g.Group("", requireTurnstile(cfg))
	}
	protected.GET("/fetch", server.Wrap(api.Fetch))
	protected.GET("/fetch/latest", server.WrapReply(api.FetchLatest))
	protected.GET("/fetch/:id", server.Wrap(api.FetchDetail))
	protected.GET("/download/:id", server.WrapReply(api.Download))
	protected.GET("/domain", server.WrapResp(api.DomainList))
}

func requireTurnstile(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if !api.TurnstileVerified(c.Request(), cfg) {
				return c.JSON(http.StatusUnauthorized, server.Envelope{
					Code:    -1,
					Message: "请先完成人机验证",
				})
			}
			return next(c)
		}
	}
}
