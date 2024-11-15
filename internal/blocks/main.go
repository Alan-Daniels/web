package blocks

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/Alan-Daniels/web/internal/blocks/types"
	"github.com/a-h/templ"
)

type Block interface {
	Component(map[string]interface{}, []*templ.Component) (templ.Component, error)
	Editor(map[string]interface{}, int, int) (templ.Component, error)
	DefArgs() (*map[string]interface{}, error)
}
type Blocks map[string]Block

type _Block[T any] struct {
	Comp func(T, []*templ.Component) templ.Component
	Name string
}

func (bl *_Block[T]) Component(args map[string]interface{}, children []*templ.Component) (templ.Component, error) {
	t, err := bl.marshal(args)
	if err != nil {
		return nil, err
	}
	return bl.Comp(*t, children), nil
}

func (bl *_Block[T]) Editor(args map[string]interface{}, blockId, parentBlockId int) (templ.Component, error) {
	t, err := bl.marshal(args)
	if err != nil {
		return nil, err
	}
	form := bl.Form(*t, blockId, parentBlockId)
	return form, nil
}

func (bl *_Block[T]) DefArgs() (*map[string]interface{}, error) {
	a := new(T)

	ref := reflect.TypeFor[T]()
	refvalue := reflect.ValueOf(a)
	fieldLen := ref.NumField()

	for i := 0; i < fieldLen; i++ {
		fieldvalue := refvalue.Elem().Field(i)
		def, ok := fieldvalue.Type().MethodByName("Default")
		if !ok {
			continue
		}
		ret := def.Func.Call([]reflect.Value{fieldvalue})
		fieldvalue.Set(ret[0])
	}

	return bl.unmarshal(a)
}

func (bl *_Block[T]) FormField(field reflect.StructField, value reflect.Value, blockId, parentBlockId int) templ.Component {
	def, ok := value.Type().MethodByName("FormField")
	if !ok {
		return types.DefaultFormField(field, value, blockId, parentBlockId)
	}
	ret := def.Func.Call([]reflect.Value{value, reflect.ValueOf(field), reflect.ValueOf(blockId), reflect.ValueOf(parentBlockId)})[0]

	return ret.Interface().(templ.Component)
}

func (bl *_Block[T]) marshal(args map[string]interface{}) (*T, error) {
	t := new(T)

	var jsonBytes []byte
	jsonBytes, err := json.Marshal(args)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialise response '%+v' to object: %w", args, err)
	}

	err = json.Unmarshal(jsonBytes, t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (bl *_Block[T]) unmarshal(args *T) (*map[string]interface{}, error) {
	t := new(map[string]interface{})

	var jsonBytes []byte
	jsonBytes, err := json.Marshal(args)

	if err != nil {
		return nil, fmt.Errorf("failed to deserialise response '%+v' to object: %w", args, err)
	}

	err = json.Unmarshal(jsonBytes, t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

var blocks Blocks

func Init() *Blocks {
	return &blocks
}

func registerBlock[T any](comp func(T, []*templ.Component) templ.Component) {
	block := new(_Block[T])
	block.Comp = comp

	name := runtime.FuncForPC(reflect.ValueOf(comp).Pointer()).Name()
	namep := strings.Split(name, "/")
	name = namep[len(namep)-1]
	block.Name = name
	if blocks == nil {
		blocks = make(Blocks)
	}
	blocks[name] = block
}
