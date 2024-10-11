package data

import (
	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/database"
	"github.com/labstack/echo/v4"
)

func (e *_branch) Groups() (groups *[]Group, err error) {
	groups, err = database.UnmarshalResponse[[]Group](Database.Branches(e.GetID()))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to get group branches")
	}

	return groups, err
}

func (b *_branch) initPages(e Branch) {
	pages, err := b.Pages()
	if err == nil {
		for _, page := range *pages {
			page.Init(e)
		}
	}
}

func (b *_branch) initGroups(e Branch) {
	groups, err := b.Groups()
	if err == nil {
		for _, group := range *groups {
			group.Init(e)
		}
	}
}

func (e *_branch) Pages() (pages *[]Page, err error) {
	pages, err = database.UnmarshalResponse[[]Page](Database.Pages(e.GetID()))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to get group pages")
	}
	return pages, err
}

func (e *Root) Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group) {
	return e.e.Group(prefix, m...)
}

func (e *Root) Add(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	return e.e.Add(method, path, handler, middleware...)
}

func (e *Root) Init(_ Branch) error {
	e.initPages(e)
	e.initGroups(e)
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
	e.initPages(e)
	e.initGroups(e)
	return nil
}
