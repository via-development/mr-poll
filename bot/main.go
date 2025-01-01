package main

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/gofor-little/env"
	"log"
	internalApi "mrpoll_bot/internal-api"
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
			cache.WithCaches(cache.FlagGuilds, cache.FlagChannels),
		),
		bot.WithEventListenerFunc(CommandHandler),
		bot.WithEventListenerFunc(ComponentHandler),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds),
		),
	)
	if err != nil {
		panic(err)
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}
	if err = client.OpenHTTPServer(); err != nil {
		panic(err)
	}
	defer client.Close(context.TODO())

	api := internalApi.NewApi(client)
	defer api.Close()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

}
