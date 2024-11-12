package data

import (
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (r *Page) Init(p Branch) error {
	comp, err := r.Block.ToComponent(0)
	if err != nil {
		Logger.Error().Err(err).Any("page", *r).Msg("Can't register endpoint!")
		return err
	}
	rt := p.Add(r.Method, r.Path, r.Handler(comp))
	rt.Name = r.GetIDString()
	return nil
}

func (r *Page) Handler(comp templ.Component) func(c echo.Context) error {
	name := r.GetIDString()
	components[name] = comp
	return func(c echo.Context) error {
		return Render(c, http.StatusOK, components[name])
	}
}

func (r *WildcardRoute) Init(p Branch) error {
	return nil
}
