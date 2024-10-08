package admin

import (
	"fmt"
	"net/http"

	. "github.com/Alan-Daniels/web/internal"
	"github.com/Alan-Daniels/web/internal/database"
	"github.com/labstack/echo/v4"
)

func Init(g *echo.Group) {
	g.GET("", admin)

	g.GET("/pl", playground)
	g.POST("/pl", playgroundPost)
}

func admin(c echo.Context) error {
	allnodes, err := Database.RootBranches()
	if err != nil {
		Logger.Error().Err(err).Msg("tried getting Root branches (classic mistake)")
		return err
	}
	c.HTML(http.StatusOK, fmt.Sprintf("%v", allnodes))
	return nil
}

func playgroundPost(c echo.Context) error {
	query := c.Request().PostFormValue("query")
	res, err := Database.Query(query, database.Map{})
	if err != nil {
		Logger.Error().Err(err).Msg("Error in the playground")
		return err
	}

	Render(c, http.StatusOK, Playground(query, res))
	return nil
}
func playground(c echo.Context) error {
	Render(c, http.StatusOK, Playground("", ""))
	return nil
}
