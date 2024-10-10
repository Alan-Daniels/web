package data

import (
	"github.com/labstack/echo/v4"
)

type Branch interface {
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
	Add(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	Init(parent Branch) error
}

type Leaf interface {
	Init(parent Branch) error
}

// --

type Group struct {
	RecordID
	Parent string      `json:"parent"`
	Prefix string      `json:"prefix"`
	g      *echo.Group `json:"-"`
}

func Init(e *echo.Echo) {
	r := Root{e: e}
	r.Init(&r)
}

type Root struct {
	e *echo.Echo
}

// --

type Route struct {
	echo.Route
}

type WildcardRoute struct {
	Group
	Route
	Tag string
}

// --

type Content struct {
	Children  []Content              `json:"children,omitempty"`
	BlockName string                 `json:"name"`
	BlockOps  map[string]interface{} `json:"opts"`
}

type Template struct {
	Content Content `json:"content,omitempty"`
	Name    string  `json:"name"`
}

type Page struct {
	Route
	Template
	RecordID
	Parent string `json:"parent"`
}

type Post struct {
	Template
	Tag string
}

type RecordID struct {
	ID interface{} `json:"id,omitempty"`
}

func (r *RecordID) GetID() string {
	switch r.ID.(type) {
	case string:
		return r.ID.(string)
	case map[string]interface{}:
		return r.ID.(map[string]interface{})["Content"].([]interface{})[1].(string)
	default:
		return ""
	}
}
