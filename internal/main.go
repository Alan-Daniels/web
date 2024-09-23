package internal

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

//go:embed commit.txt
var Commit string

var TimeZoneSydney, _ = time.LoadLocation("Australia/Sydney")

var Logger zerolog.Logger

func ComponentHandler(comp func() templ.Component) func(e echo.Context) error {
	return func(e echo.Context) error {
		_, err := InitSession(e)
		if err != nil {
			return err
		}
		return Render(e, http.StatusOK, comp())
	}
}

// This custom Render replaces Echo's echo.Context.Render()
// with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func InitSession(c echo.Context) (*sessions.Session, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return sess, err
	}
	if sess.IsNew {
		Logger.Debug().Msg("New session!!!")
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return sess, err
		}
	}

	return sess, nil
}
