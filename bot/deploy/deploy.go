package main

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/gofor-little/env"
	"log"
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

	rests := rest.New(rest.NewClient(token))
	commands := []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "mr-poll",
			Description: "Hi, I'm Mr Poll!",
		},
		discord.SlashCommandCreate{
			Name:        "help",
			Description: "Hi, I'm Mr Poll!",
		},
		discord.SlashCommandCreate{
			Name:        "poll",
			Description: "lol",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "online",
					Description: "Create a poll for Discord using our website.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "yes-or-no",
					Description: "Create a yes or no poll.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "single-choice",
					Description: "Create a yes or no poll.",
				},
			},
		},
	}

	_, err = rests.SetGuildCommands(1199127749923709089, 976147096757497937, commands)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success!")
}
