package blocks

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
)

type Block interface {
	Component(map[string]interface{}) (templ.Component, error)
}
type Blocks map[string]Block

type _Block[T any] struct {
	comp func(T) templ.Component
}

func (bl *_Block[T]) Component(args map[string]interface{}) (templ.Component, error) {
	t, err := bl.unmarshal(args)
	if err != nil {
		return nil, err
	}
	return bl.comp(*t), nil
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

func Merge(parent templ.Component, children []templ.Component) templ.Component {
	child := templruntime.GeneratedTemplate(func(input templruntime.GeneratedComponentInput) (retErr error) {
		writer, ctx := input.Writer, input.Context
		if ctxerr := ctx.Err(); ctxerr != nil {
			return ctxerr
		}
		buffer, existing := templruntime.GetBuffer(writer)
		if !existing {
			defer func() {
				err := templruntime.ReleaseBuffer(buffer)
				if retErr == nil {
					retErr = err
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		for i := range children {
			retErr = (children[i]).Render(ctx, buffer)
			if retErr != nil {
				return retErr
			}
		}
		return retErr
	})
	return templruntime.GeneratedTemplate(func(input templruntime.GeneratedComponentInput) (retErr error) {
		writer, ctx := input.Writer, input.Context
		if ctxerr := ctx.Err(); ctxerr != nil {
			return ctxerr
		}
		buffer, existing := templruntime.GetBuffer(writer)
		if !existing {
			defer func() {
				err := templruntime.ReleaseBuffer(buffer)
				if retErr == nil {
					retErr = err
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		ctx = templ.ClearChildren(ctx)
		retErr = parent.Render(templ.WithChildren(ctx, child), buffer)
		if retErr != nil {
			return retErr
		}
		return retErr
	})
}
