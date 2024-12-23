package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/data"
	"github.com/labstack/echo/v4"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

type BlockCounter struct {
	state int
}

func (bId *BlockCounter) Next() int {
	bId.state += 1
	return bId.state
}

func (bId *BlockCounter) Reset() int {
	state := bId.state
	bId.state = 0
	return state
}

func NewBlockCounter(start int) *BlockCounter {
	bId := BlockCounter{state: start}
	return &bId
}

var routes map[string]*echo.Route
var g *echo.Group

func Init(g *echo.Group) {
	routes = make(map[string]*echo.Route)

	g.GET("", admin)
	g.GET("/mkpage", mkpage)
	g.GET("/preview", preview)

	routes["editor"] = g.GET("/edit", eEditor)
	routes["editor.save"] = g.POST("/edit", eEditorSave)
	routes["editor.block.create"] = g.GET("/edit/block", eBlockCreate)
	routes["editor.block.update"] = g.POST("/edit/block", eBlockUpdate)
}

func preview(c echo.Context) error {
	blockName := c.QueryParam("block")

	block, ok := Blocks[blockName]
	if !ok {
		Logger.Warn().Msgf("cannot find block with name '%s'", blockName)
		return c.JSON(http.StatusNotFound, fmt.Sprintf("not found '%s'", blockName))
	}

	args, err := block.DefArgs()
	if err != nil {
		Logger.Err(err).Send()
		return c.JSON(http.StatusInternalServerError, err)
	}

	testChild1 := new(data.Block)
	testChild1.BlockName = "blocks.blockTest"
	testChild1.BlockOps = map[string]interface{}{
		"name": "test child 1",
	}

	testChild2 := new(data.Block)
	testChild2.BlockName = "blocks.blockTest"
	testChild2.BlockOps = map[string]interface{}{
		"name": "test child 2",
	}

	testChild3 := new(data.Block)
	testChild3.BlockName = "blocks.blockTest"
	testChild3.BlockOps = map[string]interface{}{
		"name": "test child 3",
	}

	bl := new(data.Block)
	bl.BlockName = blockName
	bl.BlockOps = *args
	bl.Children = append(bl.Children, *testChild1, *testChild2, *testChild3)

	comp, err := bl.ToComponent(0)
	if err != nil {
		Logger.Err(err).Send()
		return c.JSON(http.StatusInternalServerError, err)
	}

	return Render(c, http.StatusOK, Preview(comp))
}

func mkpage(c echo.Context) (err error) {
	surrealdb.Delete(Database, models.Table("Page"))
	surrealdb.Delete(Database, models.Table("Group"))
	outs := make([]interface{}, 3)

	page := new(data.Page)
	id := data.NewRecordID[data.Page]("rootpage")
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

	page, err = data.Insert(page)
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		outs[0] = err.Error()
	} else {
		outs[0] = page
	}

	group := new(data.Group)
	group.ID = nil
	group.Prefix = "/test"
	group, err = data.Insert(group)
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
	page, err = data.Insert(page)
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

	return Render(c, http.StatusOK, PageRoutes(rt))
}

func eBlockCreate(c echo.Context) error {
	params := c.QueryParams()
	idStr := params.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("id must be an int")
	}
	parentStr := params.Get("parentid")
	parent, err := strconv.Atoi(parentStr)
	if err != nil {
		return fmt.Errorf("parent id must be an int")
	}
	name := params.Get("name")

	b, ok := Blocks[name]
	if !ok {
		return fmt.Errorf("cant find block with name")
	}

	defops, err := b.DefArgs()
	if err != nil {
		return err
	}
	block := new(data.Block)
	block.BlockName = name
	block.BlockOps = *defops

	return Render(c, http.StatusOK, TemplEditor(*block, id, parent))
}

func eBlockUpdate(c echo.Context) error {
	var block data.Block

	err := c.Bind(&block)
	if err != nil {
		Logger.Err(err).Send()
		return err
	}

	Logger.Debug().Any("block", block).Send()
	bct := NewBlockCounter(0)

	component, err := block.EditorComponent(0, EditBlockChildren(block.Children, bct))
	if err != nil {
		Logger.Error().Err(err).Msg("failed to render component")

		c.JSON(http.StatusInternalServerError, err.Error())
		return err
	}

	return Render(c, http.StatusOK, component)
}

func eEditorSave(c echo.Context) error {
	idStr := c.QueryParams().Get("id")
	id := models.ParseRecordID(idStr)
	var block data.Block

	err := c.Bind(&block)
	if err != nil {
		Logger.Err(err).Send()
		return err
	}

	root := data.EnsureBlockRoot(block)

	switch id.Table {
	case "Page":
		return savePage(c, id, root)
	default:
		return fmt.Errorf("Could not edit Object of type %s", id.Table)
	}
}

func eEditor(c echo.Context) error {
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
	var page *data.Page
	page, err := data.FromID[data.Page](*id)
	if err != nil {
		pretty, err := json.Marshal(err.Error())
		if err != nil {
			return err
		}
		c.JSONBlob(http.StatusInternalServerError, pretty)
		return nil
	}

	block := data.EnsureBlockRoot(page.Block)

	return Render(c, http.StatusOK, PageEditor(block, page.ID))
}

func savePage(c echo.Context, id *models.RecordID, content data.RootBlock) error {
	var page *data.Page

	page, err := data.FromID[data.Page](*id)
	if err != nil {
		pretty, err := json.Marshal(err.Error())
		if err != nil {
			return err
		}
		c.JSONBlob(http.StatusInternalServerError, pretty)
		return nil
	}

	comp, err := content.ToComponent()
	if err != nil {
		Logger.Error().Err(err).Msg("to component error")
		return err
	}

	page.Block = data.Block{
		Children:  content.Children,
		BlockName: "blocks.root",
	}
	page, err = data.Update(page, *page.ID)
	if err != nil {
		Logger.Error().Err(err).Msg("DB update error")
		return c.JSON(http.StatusInternalServerError, err)
	}

	page.RInit(c, comp)

	return c.JSON(http.StatusOK, page)
}
