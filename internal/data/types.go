package data

import (
	"context"
	"fmt"
	"reflect"
	"time"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
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
var rootArgs blocks.RootArgs

func Init(e *echo.Group) error {
	components = make(map[string]templ.Component)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	(&Group{}).Init(&Group{g: e}, ctx)

	initRoot()

	return ctx.Err()
}

func initRoot() {
	var head templ.Component
	var foot templ.Component

	_head, err := FromID[Template](NewRecordID[Template]("head"))
	if err != nil {
		Logger.Error().Err(err).Msg("could not retrieve head, default to empty head")
		head = empty()
	} else {
		head, err = _head.Block.ToComponent(0)
		if err != nil {
			Logger.Error().Err(err).Msg("render head err")
			head = empty()
		}
	}
	_foot, err := FromID[Template](NewRecordID[Template]("foot"))
	if err != nil {
		Logger.Error().Err(err).Msg("could not retrieve foot, default to empty foot")
		foot = empty()
	} else {
		foot, err = _foot.Block.ToComponent(0)
		if err != nil {
			Logger.Error().Err(err).Msg("render foot err")
			foot = empty()
		}
	}

	rootArgs = blocks.RootArgs{
		Head: head,
		Foot: foot,
	}
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

type RootBlock struct {
	Children []Block
	Args     blocks.RootArgs
}

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

type HasRecordID interface {
	GetID() *models.RecordID
}

type RecordID[T any] struct {
	ID *models.RecordID `json:"id,omitempty"`
}

func (r RecordID[T]) GetID() *models.RecordID {
	return r.ID
}

func NewRecordID[T HasRecordID, K RecordID[T]](id string) models.RecordID {
	return models.NewRecordID(typeName[T](), id)
}

func (r *RecordID[T]) GetIDString() string {
	if r.ID == nil {
		return "<nil>"
	}
	return r.ID.String()
}

func Insert[T HasRecordID](t *T) (*T, error) {
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

func Update[T HasRecordID](t *T, id models.RecordID) (*T, error) {
	item, err := surrealdb.Update[T](Database, id, t)
	return item, err
}

func FromID[T HasRecordID](id models.RecordID) (item *T, err error) {
	item, err = surrealdb.Select[T](Database, id)
	if (*item).GetID() == nil {
		return item, fmt.Errorf("no record for id")
	}
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
