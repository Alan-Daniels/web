package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blog"
	"github.com/Alan-Daniels/web/internal/blog/data"
	"github.com/Alan-Daniels/web/internal/pages"
)

type Config struct {
	Server struct {
		Port     string `yaml:"port" envconfig:"SERVER_PORT"`
		HostName string `yaml:"hostname" envconfig:"SERVER_HOSTNAME"`
	} `yaml:"server"`
}

func NewConfig(file string) (*Config, error) {
	cfg := new(Config)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	staticDir := flag.String("static", ".", "where static files & default files are")
	//stateDir := flag.String("state", "./tmp", "where generated files & content files are")
	configFile := flag.String("config", "./default.yml", "see config.Init()")
	flag.Parse()

	config, err := NewConfig(*configFile)
	if err != nil {
		panic(err)
	}

	//webroot := flag.String("root", ".", "where the files be ;)")
	//logdir := flag.String("logdir", ".", "where to put logs")
	//port := flag.String("port", "8080", "port to listen on")
	//metricsPort := flag.String("metrics", "", "metrics port, default to no metrics")
	//testtag := flag.String("tag", "notest", "testing tag to use")
	//flag.Parse()

	//TestTag = *testtag

	app := echo.New()
	app.IPExtractor = echo.ExtractIPFromXFFHeader()

	app.Static("/assets", (*staticDir)+"/assets")
	app.File("/favicon.ico", (*staticDir)+"/assets/favicon.ico")

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

	//logfile, err := os.OpenFile((*logdir)+"/"+*testtag+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	//if err != nil {
	//	app.Logger.Fatal(err)
	//}
	//zerolog.TimestampFieldName = "t"
	//zerolog.LevelFieldName = "l"
	//zerolog.MessageFieldName = "m"
	//Logger = zerolog.New(logfile).With().Timestamp().Logger()

	app.Use(middleware.Gzip())
	app.Use(middleware.Secure())

	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// TODO: CORS, CSRF
	// site doesn't have interactability yet so these aren't critical

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

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", (config.Server.Port))))
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
