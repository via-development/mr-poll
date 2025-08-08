package core

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/sharding"
	"github.com/gofor-little/env"
	"github.com/via-development/mr-poll/bot/util"
)

func CreateBot() bot.Client {
	var err error
	var token string
	if token, err = env.MustGet("BOT_TOKEN"); err != nil || len(token) == 0 {
		panic("BOT_TOKEN environment variable not set")
	}

	var publicKey string
	if publicKey, err = env.MustGet("BOT_PUBLIC_KEY"); err != nil || len(token) == 0 {
		panic("BOT_PUBLIC_KEY environment variable not set")
	}

	client, err := disgo.New(token,
		bot.WithHTTPServerConfigOpts(publicKey,
			httpserver.WithURL("/bot"),
			httpserver.WithAddress(":4002"),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds, cache.FlagChannels, cache.FlagMembers, cache.FlagRoles),
		),
		bot.WithEventListenerFunc(commandHandler),
		bot.WithEventListenerFunc(componentHandler),
		bot.WithEventListenerFunc(modalHandler),
		bot.WithEventListenerFunc(messageHandler),
		bot.WithShardManagerConfigOpts(
			sharding.WithShardIDs(util.Config.ShardIds...),
			sharding.WithShardCount(util.Config.ShardCount),
			sharding.WithAutoScaling(true),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentGuilds, gateway.IntentMessageContent, gateway.IntentGuildMessages),
				gateway.WithCompress(true),
				gateway.WithPresenceOpts(gateway.WithCustomActivity("/mr-poll | Not made with AI!")),
			),
		),
		bot.WithEventListeners(&events.ListenerAdapter{
			Ready: func(e *events.GuildsReady) {
				fmt.Printf("Shard %d online!\n", e.ShardID())
				err = e.Client().SetPresenceForShard(
					context.Background(), e.ShardID(),
					gateway.WithCustomActivity(fmt.Sprintf("/mr-poll | %s Shard (%d)", util.ShardNames[e.ShardID()], e.ShardID())),
				)
				if err != nil {
					fmt.Println(err)
				}
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	return client
}
