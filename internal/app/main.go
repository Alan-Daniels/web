package app

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/admin"
	"github.com/Alan-Daniels/web/internal/data"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func Init() error {
	app := echo.New()
	app.IPExtractor = echo.ExtractIPFromXFFHeader()

	app.Static("/assets", (RootDir)+"/assets")
	app.File("/favicon.ico", (RootDir)+"/assets/favicon.ico")

	group := app.Group("", middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(2),
		Burst:     5,
		ExpiresIn: 3 * time.Minute,
	})))
	err := data.Init(group)
	if err != nil {
		Logger.Error().Err(err).Msg("Failed to build pages, only admin will be accessable")
	}

	adm := app.Group("/admin")
	admin.Init(adm)

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())

	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// TODO: CORS, CSRF

	// recover from panics
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			Logger.Error().Stack().Err(err).Msg("Encountered a PANIC while serving an endpoint")
			return err
		},
	}))

	app.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	return app.Start(fmt.Sprintf(":%s", (Config.Server.Port)))
}
