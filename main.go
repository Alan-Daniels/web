package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/time/rate"

	"github.com/Alan-Daniels/web/internal"
)

func main() {
	flSSL := flag.Bool("ssl", false, "whether to start with ssl")
	webroot := flag.String("root", ".", "where the files be ;)")
	flag.Parse()

	app := echo.New()

	app.Static("/assets", (*webroot)+"/assets")
	app.File("/favicon.ico", (*webroot)+"/assets/favicon.ico")

	app.GET("/", HomeHandler)
	app.GET("/now", NowHandler)
	app.GET("/about", AboutHandler)

	//projects := app.Group("/projects")
	//projects.GET("", ProjectsHandler)

	//notes := app.Group("/notes")
	//notes.GET("", NotesHandler)

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())

	// TODO: CORS, CSRF
	// site doesn't have interactability yet so these aren't critical

	// TODO: add logger for better error insight
	// TODO: figure out some analytics
	app.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Subsystem: "mysite",
		DoNotUseRequestPathFor404: true,
	}))
	go func() {
		metrics := echo.New()                                // this Echo will run on separate port 8081
		metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
		if err := metrics.Start(":8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal(err)
		}
	}()

	// rate limit to 20 requests per second
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

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
			//Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
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

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, internal.Home())
}

func AboutHandler(c echo.Context) error {
	return Render(c, http.StatusOK, internal.About())
}

func NowHandler(c echo.Context) error {
	return Render(c, http.StatusOK, internal.Now())
}
