package core

import (
	"context"
	"strconv"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/sharding"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/database"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const Intents = gateway.IntentGuilds | gateway.IntentMessageContent | gateway.IntentGuildMessages

type Client struct {
	config  *config.Config
	log     *zap.Logger
	db      *database.Database
	modules map[string]Module

	*bot.Client
}

type clientParams struct {
	fx.In

	Config *config.Config
	Log    *zap.Logger
	Db     *database.Database
}

func New(lc fx.Lifecycle, p clientParams) (*Client, error) {
	b := &Client{
		config:  p.Config,
		log:     p.Log,
		db:      p.Db,
		modules: map[string]Module{},
	}

	var err error
	b.Client, err = disgo.New(p.Config.BotToken,
		bot.WithHTTPServerConfigOpts(p.Config.BotPublicKey,
			httpserver.WithURL("/bot"),
			httpserver.WithAddress(":"+strconv.Itoa(b.config.BotPort)),
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

func (b *Client) Register(m Module) {
	b.modules[m.Name()] = m
}

func (b *Client) Start(ctx context.Context) error {
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

func (b *Client) Stop(ctx context.Context) error {
	b.log.Info("bot is stopping")
	b.Close(ctx)
	return nil
}
