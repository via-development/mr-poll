package main

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/gofor-little/env"
	"log"
)

var contexts = []discord.InteractionContextType{
	discord.InteractionContextTypeGuild,
	discord.InteractionContextTypeBotDM,
	discord.InteractionContextTypePrivateChannel,
}
var intregationTypes = []discord.ApplicationIntegrationType{
	discord.ApplicationIntegrationTypeGuildInstall,
	discord.ApplicationIntegrationTypeUserInstall,
}

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
			Name:             "mr-poll",
			Description:      "Hi, I'm Mr Poll!",
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
		discord.SlashCommandCreate{
			Name:             "help",
			Description:      "Hi, I'm Mr Poll!",
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
		discord.SlashCommandCreate{
			Name:        "poll",
			Description: "---",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "online",
					Description: "View and create polls on the website.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "yes-or-no",
					Description: "Create a yes or no poll.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "single-choice",
					Description: "Create a single-choice poll.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "multiple-choice",
					Description: "Create a multiple-choice poll.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "submit-choice",
					Description: "Create a submit-choice poll.",
				},
			},
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
	}

	//_, err = rests.SetGuildCommands(1199127749923709089, 976147096757497937, []discord.ApplicationCommandCreate{})
	_, err = rests.SetGlobalCommands(1199127749923709089, commands)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success!")
}
