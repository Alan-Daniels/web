package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/time/rate"

	"github.com/rs/zerolog"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blog"
	"github.com/Alan-Daniels/web/internal/blog/data"
	"github.com/Alan-Daniels/web/internal/pages"
)

func main() {
	flSSL := flag.Bool("ssl", false, "whether to start with ssl")
	webroot := flag.String("root", ".", "where the files be ;)")
	metricsPort := flag.String("metrics", "", "metrics port, default to no metrics")
	flag.Parse()

	app := echo.New()

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
	app.IPExtractor = echo.ExtractIPDirect()
	logger := zerolog.New(logfile)
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogRemoteIP:     true,
		LogResponseSize: true,
		LogError:        true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var msg *zerolog.Event
			if v.Error != nil {
				msg = logger.Error().Err(v.Error)
			} else if v.Status == 404 || v.Status == 302 || v.Status == 301 {
				msg = logger.Warn().Int64("RespSize", v.ResponseSize)
			} else if v.Status == 429 {
				msg = logger.Log()
			} else if v.Status == 200 {
				// prometheus already gives plenty of insight for these requests
				return nil
			} else {
				// catch-all
				msg = logger.Log()
			}
			msg.Timestamp().
				Str("RemoteIP", v.RemoteIP).
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())

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

	if *flSSL {
		customHTTPServer(app, *webroot)
	} else {
		app.Logger.Fatal(app.Start(":8080"))
	}
}

func customHTTPServer(e *echo.Echo, webroot string) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	autoTLSManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
		Cache:      autocert.DirCache(webroot + "/.cache"),
		HostPolicy: autocert.HostWhitelist("alandaniels.homes"),
	}
	s := http.Server{
		Addr:    ":4343",
		Handler: e, // set Echo as handler
		TLSConfig: &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		//ReadTimeout: 30 * time.Second, // use custom timeouts
	}
	go http.ListenAndServe(":8080", autoTLSManager.HTTPHandler(nil))
	if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}
