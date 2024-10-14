package data

import (
	"context"
	"fmt"
	"reflect"
	"time"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/labstack/echo/v4"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type Branch interface {
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
	Add(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route
	Init(parent Branch, ctx context.Context) error
}

type Leaf interface {
	Init(parent Branch) error
}

// --

type Group struct {
	RecordID[Group]
	HasParent[Group]
	Parent *models.RecordID `json:"parent"`
	Prefix string           `json:"prefix"`
	g      *echo.Group      `json:"-"`
}

func Init(e *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	(&Group{}).Init(&Group{g: e.Group("")}, ctx)

	return ctx.Err()
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
	RecordID[Page]
	HasParent[Page]
	Parent *models.RecordID `json:"parent"`
}

type Post struct {
	Template
	Tag string
}

// --

func typeName[T any]() string {
	return reflect.TypeOf(*(new(T))).Name()
}

type RecordID[T any] struct {
	ID *models.RecordID `json:"id,omitempty"`
}

func (r *RecordID[T]) GetIDString() string {
	return r.ID.String()
}

// static method
func (*RecordID[T]) Insert(t *T) (*T, error) {
	items, err := surrealdb.Insert[T](Database, models.Table(typeName[T]()), t)
	if err != nil {
		return nil, err
	}
	if len(*items) != 1 {
		return nil, fmt.Errorf("Expected 1 result but got %d", len(*items))
	}
	item := (*items)[0]
	return &item, err
}

// static method
func (*RecordID[T]) FromID(id models.RecordID) (item T, err error) {
	items, err := surrealdb.Select[T](Database, id)
	item = *items
	return item, err
}

type HasParent[T any] struct{}

// static method
func (*HasParent[T]) FromParentID(parentId *models.RecordID) ([]T, error) {
	results, err := surrealdb.Query[[]T](
		Database,
		fmt.Sprintf("SELECT * FROM %s WHERE parent = $parent", typeName[T]()),
		map[string]interface{}{"parent": parentId},
	)
	if err != nil {
		return nil, err
	}
	if len(*results) != 1 {
		return nil, fmt.Errorf("Expected 1 result but got %d", len(*results))
	}

	items := (*results)[0].Result

	return items, nil
}
