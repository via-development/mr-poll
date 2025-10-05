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
			Name:             "mytime",
			Description:      "Set your timezone and dateformat.",
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionFloat{
					Name:        "utc-offset",
					Description: "-",
					Required:    true,
				},
				//discord.ApplicationCommandOptionFloat{
				//	Name:        "your-time",
				//	Description: "What time is it for you right now?",
				//	Required:    true,
				//},
				discord.ApplicationCommandOptionString{
					Name:        "date-format",
					Description: "---",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "DDMMYY",
							Value: "ddmmyy",
						},
						{
							Name:  "MMDDYY",
							Value: "mmddyy",
						},
						{
							Name:  "YYMMDD",
							Value: "yymmdd",
						},
					},
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "poll",
			Description: "---",
			Options: []discord.ApplicationCommandOption{
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
						pOption("1", "first", false),
						pOption("2", "second", false),
						pOption("3", "third", false),
						pOption("4", "fourth", false),
						pOption("5", "fifth", false),
						pOption("6", "sixth", false),
						pOption("7", "seventh", false),
						pOption("8", "eighth", false),
						pOption("9", "ninth", false),
						pOption("10", "tenth", false),
						discord.ApplicationCommandOptionString{
							Name:        "timer",
							Description: "How long the poll should last",
							Required:    false,
						},
						discord.ApplicationCommandOptionBool{
							Name:        "can-submit",
							Description: "Allow people to submit options for the poll.",
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
						pOption("6", "sixth", false),
						pOption("7", "seventh", false),
						pOption("8", "eighth", false),
						pOption("9", "ninth", false),
						pOption("10", "tenth", false),
						discord.ApplicationCommandOptionString{
							Name:        "timer",
							Description: "How long the poll should last",
							Required:    false,
						},
						discord.ApplicationCommandOptionBool{
							Name:        "can-submit",
							Description: "Allow people to submit options for the poll.",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "end",
					Description: "End a poll.",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "message",
							Description: "The message id for the poll.",
							Required:    true,
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
							},
						},
					},
				},
			},
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
		discord.SlashCommandCreate{
			Name:        "suggestion",
			Description: "---",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "channel",
					Description: "---",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "info",
							Description: "---",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "name",
									Description: "---",
									Required:    true,
								},
							},
						},
						{
							Name:        "add",
							Description: "---",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "name",
									Description: "---",
									Required:    true,
								},
								discord.ApplicationCommandOptionChannel{
									Name:        "channel",
									Description: "---",
									Required:    true,
								},
								discord.ApplicationCommandOptionString{
									Name:        "cooldown",
									Description: "---",
									Required:    false,
								},
							},
						},
						{
							Name:        "remove",
							Description: "---",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "name",
									Description: "---",
									Required:    true,
								},
							},
						},
						{
							Name:        "edit",
							Description: "---",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "name",
									Description: "---",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		discord.MessageCommandCreate{
			Name:             "End poll",
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
		discord.MessageCommandCreate{
			Name:             "Refresh poll",
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
		discord.MessageCommandCreate{
			Name:             "Approve suggestion",
			Contexts:         contexts,
			IntegrationTypes: intregationTypes,
		},
		discord.MessageCommandCreate{
			Name:             "Deny suggestion",
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
