package data

import (
	"fmt"
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/labstack/echo/v4"
)

func (r *Page) Init(p Branch) error {
	p.Add(r.Method, r.Path, r.Handler())
	return nil
}

func (r *Page) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		Logger.Debug().Any("route", r).Msg("starting handler for page!")
		c.HTML(http.StatusOK, fmt.Sprintf("this is a page! %v", r))
		return nil
	}
}
func (r *WildcardRoute) Init(p Branch) error {
	return nil
}
