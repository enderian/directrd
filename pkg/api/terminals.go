package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/enderian/directrd/pkg/types"
	"github.com/labstack/echo"
)

func terminalsGroup(g *echo.Group) {
	g.GET("s", indexTerminals)
	g.POST("", createTerminal)

	g.GET("/:name", showTerminal, setTerminal)
	g.DELETE("/:name", deleteTerminal, setTerminal)
	g.PUT("/:name", updateTerminal, setTerminal)
	g.POST("/:name/execute", execCommand, setTerminal)
}

func indexTerminals(c echo.Context) error {
	var terminals []*types.Terminal
	query := ctx.DB()

	if c.QueryParam("room") != "" {
		room, _ := strconv.Atoi(c.QueryParam("room"))
		query = query.Where("room_id = ?", room)
	}
	if err := query.Find(&terminals).Error; err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, terminals)
}

func createTerminal(c echo.Context) error {
	terminal := types.Terminal{}
	_ = c.Bind(&terminal)

	if err := ctx.DB().Save(&terminal).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, terminal)
}

var setTerminal echo.MiddlewareFunc = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		terminal := &types.Terminal{}
		if err := ctx.DB().Where("name = ?", c.Param("name")).Find(&terminal).Error; err != nil {
			return fmt.Errorf("terminal with name %s not found", c.Param("name"))
		}

		c.Set("terminal", terminal)
		return next(c)
	}
}

func showTerminal(c echo.Context) error {
	return c.JSON(http.StatusOK, c.Get("terminal").(*types.Terminal))
}

func deleteTerminal(c echo.Context) error {
	if err := ctx.DB().Delete(c.Get("terminal").(*types.Terminal)).Error; err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func updateTerminal(c echo.Context) error {
	terminal := &types.Terminal{}
	terminal.ID = c.Get("terminal").(*types.Terminal).ID

	if err := c.Bind(&terminal); err != nil {
		return err
	}
	if err := ctx.DB().Save(&terminal).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, terminal)
}

func execCommand(c echo.Context) error {
	cmd := types.Command{}
	cmd.Terminal = c.Get("terminal").(*types.Terminal).Name
	_ = c.Bind(&cmd)

	commandQueue <- cmd

	return c.NoContent(http.StatusNoContent)
}
