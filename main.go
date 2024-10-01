package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/rs/zerolog"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blog"
	"github.com/Alan-Daniels/web/internal/blog/data"
	"github.com/Alan-Daniels/web/internal/pages"
)

func main() {
	webroot := flag.String("root", ".", "where the files be ;)")
	metricsPort := flag.String("metrics", "", "metrics port, default to no metrics")
	flag.Parse()

	app := echo.New()
	app.IPExtractor = echo.ExtractIPFromXFFHeader()

	app.Static("/assets", (*webroot)+"/assets")
	app.File("/favicon.ico", (*webroot)+"/assets/favicon.ico")

	app.GET("/", ComponentHandler(pages.Home))
	app.GET("/now", ComponentHandler(pages.Now))
	app.GET("/about", ComponentHandler(pages.About))

	//projects := app.Group("/projects")
	//projects.GET("", ProjectsHandler)

	appblog := app.Group("/blog")
	appblog.GET("", data.IndexHandler(blog.BlogPosts))
	for _, post := range blog.BlogPosts {
		appblog.GET(fmt.Sprintf("/%s", post.SafeName), post.Handler())
	}

	logfile, err := os.OpenFile((*webroot)+"/log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		app.Logger.Fatal(err)
	}
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	Logger = zerolog.New(logfile).With().Timestamp().Logger()
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRemoteIP:  true,
		LogError:     true,
		LogHeaders:   []string{"Cookie"},
		LogMethod:    true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var msg *zerolog.Event
			if v.Error != nil {
				msg = Logger.Error().Err(v.Error)
			} else if v.Status == 429 {
				msg = Logger.Info()
			} else if v.Status >= 300 {
				msg = Logger.Warn()

			} else if v.Status >= 200 {
				msg = Logger.Info().
					Dict("Cookies", parseCookies(v.Headers["Cookie"]))
			} else {
				// catch-all
				msg = Logger.Info()
			}
			msg.
				Str("RemoteIP", v.RemoteIP).
				Str("Agent", v.UserAgent).
				Str("Method", v.Method).
				Int("status", v.Status).
				Str("URI", v.URI).
				Msg("request")

			return nil
		},
	}))

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())

	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// TODO: CORS, CSRF
	// site doesn't have interactability yet so these aren't critical

	// TODO: figure out some analytics
	if (*metricsPort) != "" {
		app.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Subsystem:                 "mysite",
			DoNotUseRequestPathFor404: true,
		}))
		go func() {
			metrics := echo.New()                                // this Echo will run on separate port 8081
			metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
			if err := metrics.Start(fmt.Sprintf(":%s", (*metricsPort))); err != nil && !errors.Is(err, http.ErrServerClosed) {
				app.Logger.Fatal(err)
			}
		}()
	}

	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(2),
		Burst:     5,
		ExpiresIn: 3 * time.Minute,
	})))

	// recover from panics
	app.Use(middleware.Recover())

	app.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	app.Logger.Fatal(app.Start(":8080"))
}

func parseCookies(headers []string) *zerolog.Event {
	cookies := zerolog.Dict()

	for _, header := range headers {
		pairs := strings.Split(header, ";")
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				continue
			}
			cookies.Str(strings.Trim(parts[0], " "), strings.Trim(parts[1], " "))
		}
	}

	return cookies
}
