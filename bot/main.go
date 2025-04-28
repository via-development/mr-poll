package main

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
	"log"
	"mrpoll_bot/database"
	internalApi "mrpoll_bot/internal-api"
	"mrpoll_bot/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := env.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

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
		bot.WithEventListenerFunc(CommandHandler),
		bot.WithEventListenerFunc(ComponentHandler),
		bot.WithEventListenerFunc(ModalHandler),
		bot.WithEventListenerFunc(MessageHandler),
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
			OnGuildsReady: func(e *events.GuildsReady) {
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

	database.InitDB()

	fmt.Println("[Disgo]: Connecting...")
	if err = client.OpenShardManager(context.TODO()); err != nil {
		panic(err)
	}
	if err = client.OpenHTTPServer(); err != nil {
		panic(err)
	}
	defer client.Close(context.TODO())
	fmt.Println("[Disgo]: Operational!")

	api := internalApi.NewApi(client)
	defer api.Close()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

}
