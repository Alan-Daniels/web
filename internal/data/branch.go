package data

import (
	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/database"
	"github.com/labstack/echo/v4"
)

func (e *Root) Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group) {
	return e.e.Group(prefix, m...)
}

func (e *Root) Add(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	return e.e.Add(method, path, handler, middleware...)
}

func (e *Root) Init(_ Branch) error {
	pages, err := database.UnmarshalResponse[[]Page](Database.Pages(""))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to get root pages")
	}
	for _, page := range *pages {
		page.Init(e)
	}

	branches, err := database.UnmarshalResponse[[]Group](Database.Branches(""))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to get root branches")
	}
	for _, branch := range *branches {
		branch.Init(e)
	}
	return nil
}

func (e *Group) Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group) {
	return e.g.Group(prefix, m...)
}

func (e *Group) Add(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	return e.g.Add(method, path, handler, middleware...)
}

func (e *Group) Init(p Branch) error {
	e.g = p.Group(e.Prefix)

	pages, err := database.UnmarshalResponse[[]Page](Database.Pages(e.GetID()))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to get group pages")
	}
	for _, page := range *pages {
		page.Init(e)
	}

	branches, err := database.UnmarshalResponse[[]Group](Database.Branches(e.GetID()))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to get group branches")
	}
	for _, branch := range *branches {
		branch.Init(e)
	}
	return nil
}
