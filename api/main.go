package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.RouteNotFound("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
	e.GET("/callback", func(c echo.Context) error {
		fmt.Println(c.QueryParam("code"))
		return c.String(http.StatusOK, c.QueryParam("code"))
	})
	actions := e.Group("/action") // Group for api actions
	actions.Use(AuthenticateRequest)
	actions.GET("/poll/create", func(c echo.Context) error {
		fmt.Println("Yes create poll")
		return c.String(http.StatusOK, "Yes")
	})

	e.Logger.Fatal(e.Start(":4004"))
}
