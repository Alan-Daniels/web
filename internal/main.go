package internal

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

//go:embed commit.txt
var Commit string

var TimeZoneSydney, _ = time.LoadLocation("Australia/Sydney")

func ComponentHandler(comp func() templ.Component) func(e echo.Context) error {
	return func(e echo.Context) error {
		return Render(e, http.StatusOK, comp())
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
