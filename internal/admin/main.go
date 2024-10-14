package admin

import (
	"encoding/json"
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/data"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/connection"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func Init(g *echo.Group) {
	g.GET("", admin)

	g.GET("/db", testDb)

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
	page.ID = &models.RecordID{Table: "Page", ID: "rootpage"}
	page.Method = "GET"
	page.Path = ""
	page.Name = "Root"

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

	page.Content = *content

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
	page.Content.BlockName = ":)))))"
	page.Content.BlockOps = nil
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

func playgroundPost(c echo.Context) error {
	query := c.Request().PostFormValue("query")
	//res, err := surrealdb.Query[interface{}](Database, query, map[string]interface{}{})
	var res connection.RPCResponse[[]surrealdb.QueryResult[interface{}]]

	err := Database.Send(&res, "query", query, map[string]interface{}{})

	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		return err
	}

	pretty, err := json.Marshal(res)

	c.JSONBlob(http.StatusOK, pretty)
	return nil
}

func testDb(c echo.Context) error {
	page, err := (&data.Page{}).FromID(models.NewRecordID("Page", "rootpage"))
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

func playground(c echo.Context) error {
	return Render(c, http.StatusOK, Playground("", ""))
}

func ContentComponent(c data.Content) templ.Component {
	comp, _ := c.ToComponent(0)
	return comp
}
