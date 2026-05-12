package route

import (
	"tmail/internal/api"

	"github.com/labstack/echo/v5"
	"github.com/sunls24/gox/server"
)

func Register(e *echo.Echo) {
	g := e.Group("/api")
	g.POST("/report", server.WrapReplyResp(api.Report))
	g.GET("/fetch", server.Wrap(api.Fetch))
	g.GET("/fetch/latest", server.WrapReply(api.FetchLatest))
	g.GET("/fetch/:id", server.Wrap(api.FetchDetail))
	g.GET("/download/:id", server.WrapReply(api.Download))
	g.GET("/domain", server.WrapResp(api.DomainList))
}
