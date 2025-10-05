package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/golittie/timeless/pkg/dateformat"
	"github.com/labstack/echo/v4"
	"github.com/via-development/mr-poll/bot/internal"
	"github.com/via-development/mr-poll/bot/internal/database"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type Api struct {
	echo *echo.Echo

	dummyTimezoneCache map[string]string // uuid -> user id

	client bot.Client
	log    *zap.Logger
	db     *database.GormDB
}

func (a *Api) Start(ctx context.Context) error {
	go func() {
		a.log.Error("api started")
		err := a.echo.Start(":3002")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("failed to start api", zap.Error(err))
			return
		}
	}()
	return nil
}

func (a *Api) Stop(ctx context.Context) error {
	a.log.Info("api is stopping")
	return a.echo.Close()
}

func (a *Api) PostTimezone(c echo.Context) error {
	id := c.Param("id")
	userId, ok := a.dummyTimezoneCache[id]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var body struct {
		TimezoneOffset int                   `json:"timezoneOffset"` // In minutes
		DateFormat     dateformat.DateFormat `json:"dateFormat"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	fmt.Println(userId, body)

	return nil
}

func New(lc fx.Lifecycle, mpb *internal.MPBot, log *zap.Logger, db *database.GormDB) *Api {
	e := echo.New()
	a := &Api{
		client: mpb.Client,
		log:    log,
		echo:   e,
	}
	e.HideBanner = true
	e.GET("/polls", a.GetPolls)
	e.POST("/tz/:id", a.PostTimezone)

	lc.Append(fx.Hook{
		OnStart: a.Start,
		OnStop:  a.Stop,
	})

	return a
}
