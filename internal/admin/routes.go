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

func BuildRouteTree(rt *RouteTree, depth int) {
	if depth > 1024 {
		Logger.Error().Msg("Reached route max depth")
		return
	}
	pages, err := (&data.Page{}).FromParentID(rt.ID)
	if err != nil {
		Logger.Error().Err(err).Msg("Trouble getting Pages")
		rt.PageChildren = make([]data.Page, 0)
	}
	rt.PageChildren = pages

	groups, err := (&data.Group{}).FromParentID(rt.ID)
	if err != nil {
		groups = make([]data.Group, 0)
	}
	rt.GroupChildren = make([]RouteTree, len(groups))
	for i := range groups {
		nrt := NewRouteTree((groups)[i])
		BuildRouteTree(nrt, depth+1)
		rt.GroupChildren[i] = *nrt
	}
}
