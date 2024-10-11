package admin

import (
	"encoding/json"
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/data"
	"github.com/Alan-Daniels/web/internal/database"
	"github.com/labstack/echo/v4"
)

func Init(g *echo.Group) {
	g.GET("", admin)

	g.GET("/pl", playground)
	g.POST("/pl", playgroundPost)

	g.GET("/mkpage", mkpage)

	g.GET("/test", test)
}

func test(c echo.Context) error {
	content := new(data.Content)
	content.BlockName = "blocks.blockPadd"
	content.BlockOps = make(map[string]interface{})
	content.BlockOps["color"] = "red"
	hello := new(data.Content)
	hello.BlockName = "blocks.blockTest"
	hello.BlockOps = make(map[string]interface{})
	hello.BlockOps["name"] = "WORLD"
	content.Children = append(content.Children, *hello)
	chContent := new(data.Content)
	chContent.BlockName = "blocks.blockPadd"
	chContent.BlockOps = make(map[string]interface{})
	chContent.BlockOps["color"] = "green"
	chContent.Children = append(chContent.Children, *hello)
	chContent.Children = append(chContent.Children, *hello)
	content.Children = append(content.Children, *chContent)

	component, err := content.ToComponent()
	if err != nil {
		Logger.Error().Err(err).Msg("failed to render component")
		return err
	}

	return Render(c, http.StatusOK, component)
}

func mkpage(c echo.Context) error {
	Database.Query("DELETE Page;DELETE Group;", database.Map{})
	outs := make([]interface{}, 3)

	page := new(data.Page)
	page.ID = "Page:root"
	// no parent
	page.Method = "GET"
	page.Path = ""
	page.Name = "Root"
	page.Content.BlockName = "some block name :)"
	page.Content.BlockOps = nil
	page, err := database.Unmarshal[data.Page](Database.Insert("Page", page))
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[0] = err.Error()
	} else {
		page.ID = page.GetID()
		outs[0] = page
	}

	group := new(data.Group)
	group.Prefix = "/test"
	group, err = database.Unmarshal[data.Group](Database.Insert("Group", group))
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[1] = err.Error()
	} else {
		group.ID = group.GetID()
		outs[1] = group
	}

	page = new(data.Page)
	page.ID = nil
	page.Parent = group.GetID()
	page.Method = "GET"
	page.Path = "/test"
	page.Name = "Root"
	page.Content.BlockName = ":)))))"
	page.Content.BlockOps = nil
	page, err = database.Unmarshal[data.Page](Database.Insert("Page", page))
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[2] = err.Error()
	} else {
		page.ID = page.GetID()
		outs[2] = page
	}

	pretty, err := json.Marshal(outs)
	c.JSONBlob(http.StatusOK, pretty)
	return nil
}

func admin(c echo.Context) error {
	allnodes, err := database.UnmarshalResponse[[]data.Page](Database.Pages(""))
	for k, v := range *allnodes {
		// this isnt required, just nice for displaying
		(*allnodes)[k].ID = v.GetID()
	}
	if err != nil {
		Logger.Error().Err(err).Msg("tried getting Root branches (classic mistake)")
		return err
	}
	pretty, err := json.Marshal(allnodes)
	c.JSONBlob(http.StatusOK, pretty)
	return nil
}

func playgroundPost(c echo.Context) error {
	query := c.Request().PostFormValue("query")
	res, err := Database.Query(query, database.Map{})
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		return err
	}

	pretty, err := json.Marshal(res)

	c.JSONBlob(http.StatusOK, pretty)
	return nil
}

func playground(c echo.Context) error {
	return Render(c, http.StatusOK, Playground("", ""))
}
