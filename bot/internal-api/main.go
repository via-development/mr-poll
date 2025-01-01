package internalApi

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/labstack/echo/v4"
)

type InternalApi struct {
	Echo      *echo.Echo
	BotClient bot.Client
}

func (api *InternalApi) Close() {
	api.Close()
}

func NewApi(client bot.Client) *InternalApi {
	api := &InternalApi{
		BotClient: client,
		Echo:      echo.New(),
	}

	//apiGroup := api.Echo.Group("/api/")
	//apiGroup.GET("", func(c echo.Context) error {
	//	return nil
	//})
	api.Echo.Logger.Fatal(api.Echo.Start(":4003"))

	return api
}
