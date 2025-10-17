package main

import (
	"context"

	"github.com/via-development/mr-poll/bot/internal/api"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/core"
	"github.com/via-development/mr-poll/bot/internal/database"
	generalModule "github.com/via-development/mr-poll/bot/internal/general-module"
	pollModule "github.com/via-development/mr-poll/bot/internal/poll-module"
	"github.com/via-development/mr-poll/bot/internal/redis"
	suggestionModule "github.com/via-development/mr-poll/bot/internal/suggestion-module"
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
			redis.New,
			core.New,
			api.New,
			fx.Annotate(pollModule.New, fx.As(new(core.Module)), fx.ResultTags(`group:"botModules"`)),
			fx.Annotate(suggestionModule.New, fx.As(new(core.Module)), fx.ResultTags(`group:"botModules"`)),
			fx.Annotate(generalModule.New, fx.As(new(core.Module)), fx.ResultTags(`group:"botModules"`)),
		),
		fx.Invoke(func(lc fx.Lifecycle, p struct {
			fx.In

			Config  *config.Config
			Client  *core.Client
			Db      *database.Database
			Log     *zap.Logger
			Api     *api.Api
			Modules []core.Module `group:"botModules"`
		}) error {
			for _, module := range p.Modules {
				p.Client.Register(module)
			}
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					//for _, mod := range modules {
					//	if (*mod).Name() == "polls" {
					//		pm := mod.(*pollModule.PollModule)
					//	}
					//}
					//go pm.EndTimedPollsLoop()
					return nil
				},
			})

			if p.Config.AutoMigrate {
				return p.Db.RunMigrations()
			}
			return nil
		}),
	).Run()
}
