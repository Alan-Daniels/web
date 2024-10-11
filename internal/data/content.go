package data

import (
	"fmt"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
	"github.com/a-h/templ"
)

var maxDepth = 1024

func (c *Content) ToComponent(depth int) (templ.Component, error) {
	if depth >= maxDepth {
		return nil, fmt.Errorf("depth of %d reached, this could be caused by a content dependancy loop (or a really deep dependancy tree)", maxDepth)
	}
	block, ok := (*Blocks)[c.BlockName]
	if !ok {
		return nil, fmt.Errorf("Could not find block with name '%s'", c.BlockName)
	}
	comp, err := block.Component(c.BlockOps)
	if err != nil {
		return nil, err
	}

	if len(c.Children) > 0 {
		children := make([]templ.Component, len(c.Children))
		for i := range c.Children {
			child, err := (c.Children[i]).ToComponent(depth + 1)
			if err != nil {
				return nil, err
			}
			children[i] = child
		}

		return blocks.Merge(comp, children), nil
	} else {
		return comp, nil
	}
}
