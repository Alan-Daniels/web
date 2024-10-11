package admin

import (
	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/data"
)

type RouteTree struct {
	data.Group
	GroupChildren []RouteTree
	PageChildren  []data.Page
}

func NewRouteTree(g data.Group) (rt *RouteTree) {
	rt = new(RouteTree)
	rt.ID = g.ID
	rt.Parent = g.Parent
	rt.Prefix = g.Prefix
	return rt
}

func BuildRouteTree(rt *RouteTree) {
	pages, err := rt.Pages()
	if err != nil {
		Logger.Error().Err(err).Msg("Trouble getting Pages")
		rt.PageChildren = make([]data.Page, 0)
	}
	rt.PageChildren = *pages
	
	groups, err := rt.Groups()
	if err != nil {
		Logger.Error().Err(err).Msg("Trouble getting Pages")
		rt.GroupChildren = make([]RouteTree, 0)
	}
	rt.GroupChildren = make([]RouteTree, len(*groups))
	for i := range *groups {
		nrt := NewRouteTree((*groups)[i])
		BuildRouteTree(nrt)
		rt.GroupChildren[i] = *nrt
	}
}
