package internal

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	client, err := ent.New(cfg.DB)
	if err != nil {
		return err
	}
	defer client.Close()

	schedule.New(client).Run()

	e := echo.New()
	e.Pre(i18n)
	e.Use(api.Middleware(client, cfg))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisablePrintStack: true,
	}))
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.DefaultHTTPErrorHandler(err, c)
		//goland:noinspection GoTypeAssertionOnErrors
		if _, ok := err.(*echo.HTTPError); !ok {
			log.Err(err).Send()
		}
	}

	route.Register(e)
	e.StaticFS("/", echo.MustSubFS(web.FS, "dist"))
	return e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}

func i18n(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path != "/" {
			return next(c)
		}

		c.Request().URL.Path += getLang(c.Request()) + "/"
		return next(c)
	}
}

func getLang(req *http.Request) string {
	if isGoogleBot(req) {
		return constant.LangZh
	}

	if al := req.Header.Get("Accept-Language"); strings.HasPrefix(al, constant.LangZh) {
		return constant.LangZh
	}

	return constant.LangEn
}

func isGoogleBot(req *http.Request) bool {
	ua := req.Header.Get("User-Agent")
	return strings.Contains(strings.ToLower(ua), "googlebot")
}
