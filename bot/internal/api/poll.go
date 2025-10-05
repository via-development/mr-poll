package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func (a *Api) GetPolls(c echo.Context) error {
	fmt.Println("Get polls!")
	return nil
}
