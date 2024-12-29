package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthenticateRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth != "abc" {
			return c.String(http.StatusUnauthorized, "Authorization key does not correspond to a user.")
		}

		return next(c)
	}
}
