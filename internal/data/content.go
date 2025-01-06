package data

import (
	"errors"
	"fmt"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
	"github.com/a-h/templ"
)

var maxDepth = 1024

func (c *RootBlock) ToComponent() (comp templ.Component, err error) {
	children := make([]*templ.Component, len(c.Children))
	if len(c.Children) > 0 {
		errs := make([]error, 0)
		for i := range c.Children {
			child, nerr := (c.Children[i]).ToComponent(1)
			errs = append(errs, nerr)
			children[i] = &child
		}
		if len(errs) > 0 {
			err = errors.Join(errs...)
		}
	}
	comp = blocks.Merge(children)
	return comp, err
}

func (c *Block) ToComponent(depth int) (comp templ.Component, err error) {
	children := make([]*templ.Component, len(c.Children))
	if len(c.Children) > 0 {
		errs := make([]error, 0)
		for i := range c.Children {
			child, nerr := (c.Children[i]).ToComponent(depth + 1)
			errs = append(errs, nerr)
			children[i] = &child
		}
		if len(errs) > 0 {
			err = errors.Join(errs...)
		}
	}
	comp, nerr := c.component(depth, children)
	err = errors.Join(nerr, err)
	return comp, err
}

func (c *Block) Form(blockId, parentBlockId int) (comp templ.Component, err error) {
	block, ok := Blocks[c.BlockName]
	if !ok {
		err = fmt.Errorf("Could not find block with name '%s'", c.BlockName)
		comp = blocks.BrokenBlock.Form(blocks.BrokenArgs{Name: c.BlockName, Err: err.Error(), Args: c.BlockOps}, blockId, parentBlockId)
	} else {
		comp, err = block.Editor(c.BlockOps, blockId, parentBlockId)
		if err != nil {
			comp = blocks.BrokenBlock.Form(blocks.BrokenArgs{Name: c.BlockName, Err: err.Error(), Args: c.BlockOps}, blockId, parentBlockId)
		}
	}

	return comp, err
}

func (c *Block) component(depth int, children []*templ.Component) (comp templ.Component, err error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("depth of %d reached, this could be caused by a content dependancy loop (or a really deep dependancy tree)", maxDepth)
	}
	block, ok := Blocks[c.BlockName]
	if !ok {
		err = fmt.Errorf("Could not find block with name '%s'", c.BlockName)
		comp = blocks.BrokenBlock.Comp(blocks.BrokenArgs{Name: c.BlockName, Err: err.Error(), Args: c.BlockOps}, children)
	} else {
		comp, err = block.Component(c.BlockOps, children)
		if err != nil {
			comp = blocks.BrokenBlock.Comp(blocks.BrokenArgs{Name: c.BlockName, Err: err.Error(), Args: c.BlockOps}, children)
		}
	}

	return comp, err
}

func (c *Block) EditorComponent(depth int, child templ.Component) (comp templ.Component, err error) {
	return c.component(depth, []*templ.Component{&child})
}

func EnsureBlockRoot(b Block) RootBlock {
	rootName := "blocks.root"
	if b.BlockName == rootName || b.BlockName == "" {
		return RootBlock{
			Children: b.Children,
			Args: rootArgs,
		}
	} else {
		Logger.Warn().Any("block", b).Msg("gotten a non-root block as a root block")
		return RootBlock{
			Children: []Block{b},
			Args: rootArgs,
		}
	}
}
