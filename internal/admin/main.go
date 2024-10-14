package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/data"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func Init(g *echo.Group) {
	g.GET("", admin)

	g.GET("/db", testDb2)

	g.GET("/mkpage", mkpage)

	g.GET("/test", test)

	g.GET("/edit", edit)
}

func test(c echo.Context) error {
	content := new(data.Block)
	content.BlockName = "blocks.blockPadd"
	content.BlockOps = make(map[string]interface{})
	content.BlockOps["color"] = "red"
	hello := new(data.Block)
	hello.BlockName = "blocks.blockTest"
	hello.BlockOps = make(map[string]interface{})
	hello.BlockOps["name"] = "WORLD"
	content.Children = append(content.Children, *hello)
	chContent := new(data.Block)
	chContent.BlockName = "blocks.blockPadd"
	chContent.BlockOps = make(map[string]interface{})
	chContent.BlockOps["color"] = "green"
	chContent.Children = append(chContent.Children, *hello)
	chContent.Children = append(chContent.Children, *hello)
	content.Children = append(content.Children, *chContent)

	component, err := content.ToComponent(0)
	if err != nil {
		Logger.Error().Err(err).Msg("failed to render component")
		return err
	}

	return Render(c, http.StatusOK, component)
}

func mkpage(c echo.Context) (err error) {
	surrealdb.Delete(Database, models.Table("Page"))
	surrealdb.Delete(Database, models.Table("Group"))
	outs := make([]interface{}, 3)

	page := new(data.Page)
	id := page.NewRecordID("rootpage")
	page.ID = &id
	page.Method = "GET"
	page.Path = ""
	page.Name = "Root"

	content := new(data.Block)
	content.BlockName = "blocks.blockPadd"
	content.BlockOps = make(map[string]interface{})
	content.BlockOps["color"] = "red"
	hello := new(data.Block)
	hello.BlockName = "blocks.blockTest"
	hello.BlockOps = make(map[string]interface{})
	hello.BlockOps["name"] = "WORLD"
	content.Children = append(content.Children, *hello)
	chContent := new(data.Block)
	chContent.BlockName = "blocks.blockPadd"
	chContent.BlockOps = make(map[string]interface{})
	chContent.BlockOps["color"] = "green"
	chContent.Children = append(chContent.Children, *hello)
	chContent.Children = append(chContent.Children, *hello)
	content.Children = append(content.Children, *chContent)

	page.Block = *content

	page, err = page.Insert(page)
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[0] = err.Error()
	} else {
		outs[0] = page
	}

	group := new(data.Group)
	group.ID = nil
	group.Prefix = "/test"
	group, err = group.Insert(group)
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[1] = err.Error()
	} else {
		outs[1] = group
	}

	page = new(data.Page)
	page.ID = nil
	page.Parent = group.ID
	page.Method = "GET"
	page.Path = "/test"
	page.Name = "Root"
	page.Block.BlockName = ":)))))"
	page.Block.BlockOps = nil
	page, err = page.Insert(page)
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[2] = err.Error()
	} else {
		outs[2] = page
	}

	pretty, err := json.Marshal(outs)
	c.JSONBlob(http.StatusOK, pretty)
	return nil
}

func admin(c echo.Context) error {
	rt := new(RouteTree)
	BuildRouteTree(rt, 0)

	return Render(c, http.StatusOK, ShowRoutes(rt))
}

func edit(c echo.Context) error {
	idStr := c.QueryParams().Get("id")
	id := models.ParseRecordID(idStr)
	switch id.Table {
	case "Page":
		return editPage(c, id)
	default:
		return fmt.Errorf("Could not edit Object of type %s", id.Table)
	}
}

func editPage(c echo.Context, id *models.RecordID) error {
	var page data.Page
	page, err := (&page).FromID(*id)
	if err != nil {
		pretty, err := json.Marshal(err.Error())
		if err != nil {
			return err
		}
		c.JSONBlob(http.StatusInternalServerError, pretty)
		return nil
	}

	// todo: make an editor
	return Render(c, http.StatusOK, Editor(page.Block, page.ID))
}

func testDb(c echo.Context) error {
	var page data.Page
	page, err := (&page).FromID(page.NewRecordID("rootpage"))
	if err != nil {
		pretty, err := json.Marshal(err.Error())
		if err != nil {
			return err
		}
		c.JSONBlob(http.StatusInternalServerError, pretty)
		return nil
	}

	return Render(c, http.StatusOK, ShowPage(&page))
}

func testDb2(c echo.Context) error {
	var page data.Page
	pages, err := (&page).FromParentID(nil)

	pretty, err := json.Marshal(pages)
	if err != nil {
		pretty, err := json.Marshal(err.Error())
		if err != nil {
			return err
		}
		c.JSONBlob(http.StatusInternalServerError, pretty)
		return nil
	}

	c.JSONBlob(http.StatusInternalServerError, pretty)
	return nil
}

func ContentComponent(c data.Block) templ.Component {
	comp, _ := c.ToComponent(0)
	return comp
}
