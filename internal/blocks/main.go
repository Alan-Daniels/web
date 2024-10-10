package blocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"

	"github.com/a-h/templ"
)

type Block interface {
	Render(map[string]interface{}, context.Context, io.Writer) error
}
type Blocks map[string]Block

type _Block[T any] struct {
	comp func(T) templ.Component
}

func (bl *_Block[T]) Render(args map[string]interface{}, ctx context.Context, w io.Writer) error {
	t, err := bl.unmarshal(args)
	if err != nil {
		return err
	}
	return bl.comp(*t).Render(ctx, w)
}

func (bl *_Block[T]) unmarshal(args map[string]interface{}) (*T, error) {
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
	if blocks == nil {
		blocks = make(Blocks)
	}
	blocks[name] = block
}
