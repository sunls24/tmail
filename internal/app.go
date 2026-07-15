package internal

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"tmail/config"
	"tmail/ent"
	"tmail/internal/api"
	"tmail/internal/constant"
	"tmail/internal/route"
	"tmail/internal/schedule"
	"tmail/web"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sunls24/gox"
	"github.com/sunls24/gox/server"
)

type App struct {
}

func NewApp() App {
	return App{}
}

func (App) init() {
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "06-01-02 15:04:05"})
}

func (app App) Run() error {
	app.init()
	cfg := config.MustNew()
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	if !cfg.TurnstileEnabled() {
		log.Warn().Msg("Turnstile verification is disabled")
	}

	client, err := ent.New(cfg.DB)
	if err != nil {
		return err
	}
	defer client.Close()

	return server.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), func(srv *server.Server) {
		srv.Echo.Pre(i18n)
		ctx := api.ServerContext(srv, cfg, client)
		schedule.New(ctx).Run()

		route.Register(srv.Echo, cfg)
		srv.Echo.StaticFS("/", echo.MustSubFS(web.FS, "dist"))
	})
}

func i18n(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if c.Request().URL.Path != "/" {
			return next(c)
		}

		c.Request().URL.Path += getLang(c.Request()) + "/"
		return next(c)
	}
}

func getLang(req *http.Request) string {
	al := req.Header.Get("Accept-Language")
	return gox.If(isSearchEngineBot(req) || strings.HasPrefix(al, constant.LangZh), constant.LangZh, constant.LangEn)
}

func isSearchEngineBot(req *http.Request) bool {
	ua := strings.ToLower(req.Header.Get("User-Agent"))
	for _, bot := range []string{"googlebot", "bingbot", "duckduckbot", "baiduspider", "yandexbot"} {
		if strings.Contains(ua, bot) {
			return true
		}
	}
	return false
}
