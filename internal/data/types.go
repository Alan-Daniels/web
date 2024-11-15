package data

import (
	"context"
	"fmt"
	"reflect"
	"time"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/a-h/templ"
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

var components map[string]templ.Component

func Init(e *echo.Group) error {
	components = make(map[string]templ.Component)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	(&Group{}).Init(&Group{g: e}, ctx)

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

type Block struct {
	Children  []Block                `json:"children,omitempty"`
	BlockName string                 `json:"name"`
	BlockOps  map[string]interface{} `json:"opts"`
}

type Content[T any] struct {
	Block Block  `json:"content,omitempty"`
	Name  string `json:"name"`
	RecordID[T]
	HasParent[T]
	Parent *models.RecordID `json:"parent"`
}

type Template struct {
	Content[Template]
}

type Page struct {
	Route
	Content[Page]
}

type Post struct {
	Content[Post]
	Tag string
}

// --

// static method
func typeName[T any]() string {
	return reflect.TypeOf(*(new(T))).Name()
}

type RecordID[T any] struct {
	ID *models.RecordID `json:"id,omitempty"`
}

func NewRecordID[T any, K RecordID[T]](id string) models.RecordID {
	return models.NewRecordID(typeName[T](), id)
}

func (r *RecordID[T]) GetIDString() string {
	if r.ID == nil {
		return "<nil>"
	}
	return r.ID.String()
}

func Insert[T any, K RecordID[T]](t *T) (*T, error) {
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

func Update[T any, K RecordID[T]](t *T, id models.RecordID) (*T, error) {
	item, err := surrealdb.Update[T](Database, id, t)
	return item, err
}

func FromID[T any, K RecordID[T]](id models.RecordID) (item *T, err error) {
	item, err = surrealdb.Select[T](Database, id)
	return item, err
}

type HasParent[T any] struct{}

// static method
func (*HasParent[T]) FromParentID(parentId *models.RecordID) (items []T, err error) {
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

	items = (*results)[0].Result

	return items, nil
}
