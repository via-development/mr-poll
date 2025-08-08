package main

import (
	"context"
	"github.com/via-development/mr-poll/bot/internal"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/database"
	pollModule "github.com/via-development/mr-poll/bot/internal/poll-module"
	pollUtil "github.com/via-development/mr-poll/bot/internal/poll-module/util"
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
			fx.Annotate(pollModule.New, fx.As(new(moduleUtil.Module)), fx.ResultTags(`group:"botModules"`)),
		),
		fx.Invoke(func(lc fx.Lifecycle, config *config.Config, client *internal.MPBot, db *database.GormDB, log *zap.Logger) error {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go pollUtil.EndTimedPollsLoop(client, db, log)
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
