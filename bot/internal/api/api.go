package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golittie/timeless/pkg/dateformat"
	"github.com/labstack/echo/v4"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/core"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Api struct {
	echo *echo.Echo

	DummyTimezoneCache map[string]string // uuid -> user id

	client *core.Client
	log    *zap.Logger
	db     *database.Database
	config *config.Config
	redis  *redis.Client
}

func (a *Api) Start(ctx context.Context) error {
	go func() {
		a.log.Error("api started")
		err := a.echo.Start(":" + strconv.Itoa(a.config.ApiPort))
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
	userId := a.redis.Get(context.Background(), redis.TimezoneKey(id))

	if userId.Err() != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var body struct {
		TimezoneOffset int                   `json:"offset"` // In minutes
		DateFormat     dateformat.DateFormat `json:"dateFormat"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	fmt.Println(userId.Val(), body)

	return c.NoContent(http.StatusOK)
}

func New(lc fx.Lifecycle, client *core.Client, log *zap.Logger, db *database.Database, config *config.Config, redis *redis.Client) *Api {
	e := echo.New()
	a := &Api{
		client: client,
		log:    log,
		echo:   e,
		config: config,
		redis:  redis,

		DummyTimezoneCache: map[string]string{},
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
