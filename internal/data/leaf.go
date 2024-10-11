package data

import (
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/labstack/echo/v4"
)

func (r *Page) Init(p Branch) error {
	comp, err := r.Content.ToComponent(0)
	if err != nil {
		Logger.Error().Err(err).Any("page", *r).Msg("Can't register endpoint!")
		return err
	}
	p.Add(r.Method, r.Path, func(c echo.Context) error {
		return Render(c, http.StatusOK, comp)
	})
	return nil
}

func (r *WildcardRoute) Init(p Branch) error {
	return nil
}
