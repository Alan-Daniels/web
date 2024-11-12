package data

import (
	"errors"
	"fmt"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
	"github.com/a-h/templ"
)

var maxDepth = 1024

func (c *Block) ToComponent(depth int) (comp templ.Component, err error) {
	comp, err = c.EditorComponent(depth)

	if len(c.Children) > 0 {
		errs := make([]error, 0)
		if err != nil {
			errs = append(errs, err)
		}
		children := make([]templ.Component, len(c.Children))
		for i := range c.Children {
			child, nerr := (c.Children[i]).ToComponent(depth + 1)
			errs = append(errs, nerr)
			children[i] = child
		}
		if len(errs) > 0 {
			err = errors.Join(errs...)
		}
		return blocks.Merge(comp, children), err
	} else {
		return comp, err
	}
}

func (c *Block) Form(blockId, parentBlockId int) (comp templ.Component, err error) {
	block, ok := Blocks[c.BlockName]
	if !ok {
		err = fmt.Errorf("Could not find block with name '%s'", c.BlockName)
		comp = brokenBlock(c, err)
	} else {
		comp, err = block.Editor(c.BlockOps, blockId, parentBlockId)
		if err != nil {
			comp = brokenBlock(c, err)
		}
	}

	return comp, err
}

func (c *Block) EditorComponent(depth int) (comp templ.Component, err error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("depth of %d reached, this could be caused by a content dependancy loop (or a really deep dependancy tree)", maxDepth)
	}
	block, ok := Blocks[c.BlockName]
	if !ok {
		err = fmt.Errorf("Could not find block with name '%s'", c.BlockName)
		comp = brokenBlock(c, err)
	} else {
		comp, err = block.Component(c.BlockOps)
		if err != nil {
			comp = brokenBlock(c, err)
		}
	}

	return comp, err
}
