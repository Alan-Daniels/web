package data

import (
	"context"
	"fmt"
	"reflect"
	"time"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/database"
	"github.com/labstack/echo/v4"
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
	Parent string      `json:"parent"`
	Prefix string      `json:"prefix"`
	g      *echo.Group `json:"-"`
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
	Parent string `json:"parent"`
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
	//ID *models.RecordID `json:"id,omitempty"`
	ID interface{} `json:"id,omitempty"`
}

func (r *RecordID[T]) GetIDString() string {
	//if r.ID == nil || r.ID.ID == nil {
	//	return ""
	//}
	//return r.ID.ID.(string)
	switch r.ID.(type) {
	case string:
		return r.ID.(string)
	case map[string]interface{}:
		return fmt.Sprintf("%s:%s", typeName[T](), r.ID.(map[string]interface{})["Content"].([]interface{})[1].(string))
	default:
		return ""
	}
}

// static method
func (*RecordID[T]) Insert(t *T) (*T, error) {
	item, err := database.Unmarshal[T](Database.Insert(typeName[T](), t))
	return item, err
}

// static method
func (*RecordID[T]) FromID(id string) (item T, err error) {
	items, err := database.UnmarshalResponse[[]T](Database.QueryFirst(
		fmt.Sprintf("SELECT * FROM %s WHERE id = $id", typeName[T]()),
		database.Map{"id": []string{"Page", "rootpage"}},
	))
	if len(*items) != 1 {
		return item, fmt.Errorf("Expected 1 result but got %d results", len(*items))
	}
	item = (*items)[0]
	return item, err
}

type HasParent[T any] struct{}

// static method
func (*HasParent[T]) FromParentID(partneId string) ([]T, error) {
	items, err := database.UnmarshalResponse[[]T](Database.QueryFirst(
		fmt.Sprintf("SELECT * FROM %s WHERE parent = $parent", typeName[T]()),
		database.Map{"parent": partneId},
	))
	return *items, err
}
