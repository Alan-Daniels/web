package data

import (
	"errors"
	"fmt"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
	"github.com/a-h/templ"
)

var maxDepth = 1024

func (c *Content) ToComponent(depth int) (comp templ.Component, err error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("depth of %d reached, this could be caused by a content dependancy loop (or a really deep dependancy tree)", maxDepth)
	}
	block, ok := Blocks[c.BlockName]
	if !ok {
		err = fmt.Errorf("Could not find block with name '%s'", c.BlockName)
		comp = brokenContent(c, err)
	} else {
		comp, err = block.Component(c.BlockOps)
		if err != nil {
			comp = brokenContent(c, err)
		}
	}

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
