package internal

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/sharding"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/database"
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const Intents = gateway.IntentGuilds | gateway.IntentMessageContent | gateway.IntentGuildMessages

type MPBotParams struct {
	fx.In

	Config  *config.Config
	Log     *zap.Logger
	Db      *database.GormDB
	Modules []moduleUtil.Module `group:"botModules"`
}

type MPBot struct {
	config  *config.Config
	log     *zap.Logger
	db      *database.GormDB
	modules map[string]moduleUtil.Module

	bot.Client
}

func NewMPBot(lc fx.Lifecycle, p MPBotParams) (*MPBot, error) {
	mods := map[string]moduleUtil.Module{}
	for _, m := range p.Modules {
		mods[m.Name()] = m
		p.Log.Info(m.Name() + " module loaded")
	}

	b := &MPBot{
		config:  p.Config,
		log:     p.Log,
		db:      p.Db,
		modules: mods,
	}

	var err error
	b.Client, err = disgo.New(p.Config.BotToken,
		bot.WithHTTPServerConfigOpts(p.Config.BotPublicKey,
			httpserver.WithURL("/bot"),
			httpserver.WithAddress(":3001"),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds, cache.FlagChannels, cache.FlagMembers, cache.FlagRoles),
		),
		bot.WithEventListenerFunc(b.HandleReady),
		bot.WithEventListenerFunc(b.HandleCommandInteraction),
		bot.WithEventListenerFunc(b.HandleComponentInteraction),
		bot.WithEventListenerFunc(b.HandleModalSubmitInteraction),
		bot.WithEventListenerFunc(b.HandleMessage),
		bot.WithShardManagerConfigOpts(
			sharding.WithShardIDs(p.Config.ShardIds...),
			sharding.WithShardCount(p.Config.ShardCount),
			sharding.WithAutoScaling(true),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(Intents),
				gateway.WithCompress(true),
				gateway.WithAutoReconnect(true),
			),
		),
	)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: b.Start,
		OnStop:  b.Stop,
	})

	return b, nil
}

func (b *MPBot) Start(ctx context.Context) error {
	b.log.Info("starting sharding manager")
	if err := b.OpenShardManager(ctx); err != nil {
		return err
	}
	b.log.Info("sharding manager started")
	b.log.Info("starting http server")
	if err := b.OpenHTTPServer(); err != nil {
		return err
	}
	b.log.Info("http server started")
	return nil
}

func (b *MPBot) Stop(ctx context.Context) error {
	b.log.Info("bot is stopping")
	b.Close(ctx)
	return nil
}
