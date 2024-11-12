package blocks

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/a-h/templ"
)

type Block interface {
	Component(map[string]interface{}) (templ.Component, error)
	Editor(map[string]interface{}, int, int) (templ.Component, error)
	DefArgs() (*map[string]interface{}, error)
}
type Blocks map[string]Block

type _Block[T any] struct {
	comp func(T) templ.Component
	name string
}

func (bl *_Block[T]) Component(args map[string]interface{}) (templ.Component, error) {
	t, err := bl.marshal(args)
	if err != nil {
		return nil, err
	}
	return bl.comp(*t), nil
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
	return bl.unmarshal(a)
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

func registerBlock[T any](comp func(T) templ.Component) {
	block := new(_Block[T])
	block.comp = comp

	name := runtime.FuncForPC(reflect.ValueOf(comp).Pointer()).Name()
	namep := strings.Split(name, "/")
	name = namep[len(namep)-1]
	block.name = name
	if blocks == nil {
		blocks = make(Blocks)
	}
	blocks[name] = block
}
