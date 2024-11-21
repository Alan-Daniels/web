package data

import (
	"fmt"
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/blocks"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func (r *Page) Init(p Branch) error {
	var rblock RootBlock
	if r.Block.BlockName == "blocks.root" {
		rblock = RootBlock{
			Children: r.Block.Children,
		}
	} else {
		rblock = RootBlock{
			Children: []Block{r.Block},
		}
	}

	comp, err := rblock.ToComponent()
	if err != nil {
		Logger.Error().Err(err).Any("page", *r).Msg("Can't register endpoint!")
		return err
	}

	rt := p.Add(r.Method, r.Path, r.Handler(comp))
	rt.Name = r.GetIDString()
	return nil
}

func (r *Page) RInit(ctx echo.Context, c templ.Component) error {
	_, ok := components[r.GetIDString()]
	newHandler := r.Handler(c)
	if !ok {
		var gr *Group
		var err error
		path := r.Path
		parent := r.Parent
		for parent != nil {
			gr, err = FromID[Group](*parent)
			if err != nil {
				break
			}
			path = fmt.Sprintf("%s%s", gr.Prefix, path)
			parent = gr.Parent
		}
		Logger.Debug().Str("path", path).Bool("new", !ok).Send()
		ctx.Echo().Add(r.Method, path, newHandler)
	}
	return nil
}

func (r *Page) Handler(comp templ.Component) func(c echo.Context) error {
	name := r.GetIDString()
	components[name] = comp
	return func(c echo.Context) error {
		cmp := components[name]
		return Render(c, http.StatusOK, blocks.RootPage(rootArgs, cmp))
	}
}

func (r *WildcardRoute) Init(p Branch) error {
	return nil
}
