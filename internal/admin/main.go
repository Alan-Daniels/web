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
}

func mkpage(c echo.Context) error {
	Database.Query("Delete Page;", database.Map{})

	page := new(data.Page)
	page.ID = "Page:root"
	page.Method = "GET"
	page.Path = ""
	page.Name = "Root"
	page.Content.BlockName = "some block name :)"
	page.Content.BlockOps = nil

	res, err := Database.Insert("Page", page)
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
	}

	pretty, err := json.Marshal(res)
	c.JSONBlob(http.StatusOK, pretty)
	return nil
}

func admin(c echo.Context) error {
	allnodes, err := database.Unmarshal[[]data.Page](Database.Pages(""))
	for k, v := range *allnodes {
		// this isnt required, just nice for displaying
		(*allnodes)[k].ID = v.GetID()
	}
	//allnodes, err := Database.Pages("")
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
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		Render(c, http.StatusOK, Playground(query, res))
		return nil
	}

	Render(c, http.StatusOK, Playground(query, string(pretty)))
	return nil
}

func playground(c echo.Context) error {
	Render(c, http.StatusOK, Playground("", ""))
	return nil
}
