package data

import (
	"context"

	"github.com/labstack/echo/v4"
)

func (b *Group) initPages(e Branch) {
	pages, err := (&Page{}).FromParentID(b.GetIDString())
	if err == nil {
		for _, page := range pages {
			page.Init(e)
		}
	}
}

func (b *Group) initGroups(e Branch, ctx context.Context) {
	groups, err := (&Group{}).FromParentID(b.GetIDString())
	if err == nil {
		for _, group := range groups {
			select {
			case <-ctx.Done():
				return
			default:
				group.Init(e, ctx)
			}
		}
	}
}

func (e *Group) Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group) {
	return e.g.Group(prefix, m...)
}

func (e *Group) Add(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	return e.g.Add(method, path, handler, middleware...)
}

func (e *Group) Init(p Branch, ctx context.Context) error {
	e.g = p.Group(e.Prefix)
	e.initPages(e)
	e.initGroups(e, ctx)
	return nil
}
