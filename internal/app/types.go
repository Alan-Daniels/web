package app

import "github.com/labstack/echo/v4"

type Branch interface {
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Init() error
}
type Leaf interface {
	echo.HandlerFunc
}

// --

type Group struct {
	echo.Group
	Parent *Branch
	Prefix string
}

func (e *Group) Init() error {
	panic("not done")
}

type Root struct {
	e *echo.Echo
}

func (e *Root) Init() error {
	panic("not done")
}

type Route struct {
	echo.Route
	Parent *Branch
}

type WildcardRoute struct {
	Group
	Route
	Parent *Branch
	Tag    string
}

// --

type Content struct {
	Children  []Content
	BlockName string
	BlockOps  map[string]interface{}
}

type Template struct {
	Content Content
	Name    string
}

type Page struct {
	Template
}

type Post struct {
	Template
	Tag string
}
