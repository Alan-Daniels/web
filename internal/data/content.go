package data

import (
	"fmt"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
	"github.com/a-h/templ"
)

func (c *Content) ToComponent() (templ.Component, error) {
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
			child, err := (c.Children[i]).ToComponent()
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
