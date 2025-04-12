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

	//var botId string
	//if botId, err = env.MustGet("BOT_ID"); err != nil || len(token) == 0 {
	//	panic("BOT_ID environment variable not set")
	//}

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
				//discord.ApplicationCommandOptionSubCommand{
				//	Name:        "online",
				//	Description: "View and create polls on the website.",
				//},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "yes-or-no",
					Description: "Create a yes or no poll.",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "question",
							Description: "The question for the poll.",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "timer",
							Description: "How long the poll should last",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "single-choice",
					Description: "Create a single-choice poll.",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "question",
							Description: "The question for the poll.",
							Required:    true,
						},
						pOption("1", "first", true),
						pOption("2", "second", true),
						pOption("3", "third", false),
						pOption("4", "fourth", false),
						pOption("5", "fifth", false),
						discord.ApplicationCommandOptionString{
							Name:        "timer",
							Description: "How long the poll should last",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "multiple-choice",
					Description: "Create a multiple-choice poll.",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "question",
							Description: "The question for the poll.",
							Required:    true,
						},
						pOption("1", "first", true),
						pOption("2", "second", true),
						pOption("3", "third", false),
						pOption("4", "fourth", false),
						pOption("5", "fifth", false),
						discord.ApplicationCommandOptionString{
							Name:        "timer",
							Description: "How long the poll should last",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "submit-choice",
					Description: "Create a submit-choice poll.",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "question",
							Description: "The question for the poll.",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "timer",
							Description: "How long the poll should last",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "preference",
					Description: "Configure your poll preferences",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "emojis",
							Description: "Configure your poll emojis",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "yes-or-no",
									Description: "The 2 emojis for your yes-or-no polls.",
								},
								discord.ApplicationCommandOptionString{
									Name:        "single-choice",
									Description: "The upto 10 emojis for your single-choice polls.",
								},
								discord.ApplicationCommandOptionString{
									Name:        "multiple-choice",
									Description: "The upto 10 emojis for your multiple-choice polls.",
								},
								discord.ApplicationCommandOptionString{
									Name:        "submit-choice",
									Description: "The upto 10 emojis for your submit-choice polls.",
								},
							},
						},
					},
				},
			},
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
	}

	//_, err = rests.SetGuildCommands(1199127749923709089, 976147096757497937, []discord.ApplicationCommandCreate{})
	_, err = rests.SetGlobalCommands(781304567052238858, commands)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success!")
}

func pOption(i, a string, r bool) discord.ApplicationCommandOptionString {
	return discord.ApplicationCommandOptionString{
		Name:        "option-" + i,
		Description: "The " + a + " option for the poll.",
		Required:    r,
	}
}
