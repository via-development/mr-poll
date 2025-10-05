package main

import (
	"context"
	"github.com/via-development/mr-poll/bot/internal"
	"github.com/via-development/mr-poll/bot/internal/api"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/database"
	generalModule "github.com/via-development/mr-poll/bot/internal/general-module"
	pollModule "github.com/via-development/mr-poll/bot/internal/poll-module"
	suggestionModule "github.com/via-development/mr-poll/bot/internal/suggestion-module"
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			zap.NewDevelopment,
			config.New,
			database.New,
			internal.NewMPBot,
			api.New,
			fx.Annotate(pollModule.New, fx.As(new(moduleUtil.Module)), fx.ResultTags(`group:"botModules"`)),
			fx.Annotate(suggestionModule.New, fx.As(new(moduleUtil.Module)), fx.ResultTags(`group:"botModules"`)),
			fx.Annotate(generalModule.New, fx.As(new(moduleUtil.Module)), fx.ResultTags(`group:"botModules"`)),
		),
		fx.Invoke(func(lc fx.Lifecycle, config *config.Config, client *internal.MPBot, db *database.GormDB, log *zap.Logger, api *api.Api, pm *pollModule.PollModule) error {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go pm.EndTimedPollsLoop()
					return nil
				},
			})

			if config.AutoMigrate {
				return db.RunMigrations()
			}
			return nil
		}),
	).Run()
}
