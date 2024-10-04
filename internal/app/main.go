package app

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/Alan-Daniels/web/internal"
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

	//app.GET("/", ComponentHandler(pages.Home))
	//app.GET("/now", ComponentHandler(pages.Now))
	//app.GET("/about", ComponentHandler(pages.About))

	//appblog := app.Group("/blog")
	//appblog.GET("", data.IndexHandler(blog.BlogPosts))
	//for _, post := range blog.BlogPosts {
	//	appblog.GET(fmt.Sprintf("/%s", post.SafeName), post.Handler())
	//}

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())

	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// TODO: CORS, CSRF

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

	return app.Start(fmt.Sprintf(":%s", (Config.Server.Port)))
}
