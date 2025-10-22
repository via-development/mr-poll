package suggestionModule

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
)

func (m *SuggestionModule) CreateSugggestion(suggestion *schema.Suggestion, suggestionChannel *schema.SuggestionChannel) (*discord.Message, util.UtilError) {
	// Channel permission check
	channel, uErr := util.PrepareChannel(m.client.Client, suggestion.GuildIdSnowflake(), suggestion.ChannelIdSnowflake(), discord.PermissionSendMessages)
	if uErr != nil {
		return nil, uErr
	}

	if !suggestion.AnonymousAuthor {
		err := m.FetchSuggestionUser(suggestion)
		if err != nil {
			return nil, util.NewFaultError(err)
		}
	}

	message, err := m.client.Rest.CreateMessage(channel.ID(), discord.MessageCreate{
		Embeds: []discord.Embed{
			m.MakeSuggestionEmbed(suggestion, suggestionChannel),
		},
		Components: MakeSuggestionComponents(suggestion),
	})
	if err != nil {
		return nil, util.NewNaturalError(err)
	}

	suggestion.MessageId = message.ID.String()
	if err = m.db.Create(suggestion).Error; err != nil {
		return nil, util.NewFaultError(err)
	}

	return message, nil
}

func (m *SuggestionModule) VoteSuggestion(suggestion *schema.Suggestion, userId string, upvote bool) (string, error) {
	action := ""
	for i, upvoter := range suggestion.Upvotes {
		if upvoter == userId {
			suggestion.Upvotes = append(suggestion.Upvotes[:i], suggestion.Upvotes[i+1:]...)
			if upvote {
				action = "removed"
			} else {
				action = "changed"
			}
		}
	}

	if action == "" {
		for i, downvoter := range suggestion.Downvotes {
			if downvoter == userId {
				suggestion.Downvotes = append(suggestion.Downvotes[:i], suggestion.Downvotes[i+1:]...)
				if !upvote {
					action = "removed"
				} else {
					action = "changed"
				}
			}
		}
	}

	if action != "removed" {
		action = "added"
		if upvote {
			suggestion.Upvotes = append(suggestion.Upvotes, userId)
		} else {
			suggestion.Downvotes = append(suggestion.Downvotes, userId)
		}
	}

	return action, m.db.UpdateColumns(&schema.Suggestion{
		Upvotes:   suggestion.Upvotes,
		Downvotes: suggestion.Downvotes,
	}).Where(&schema.Suggestion{
		MessageId: suggestion.MessageId,
	}).Error
}

func (m *SuggestionModule) FetchSuggestionUser(suggestion *schema.Suggestion) error {
	if suggestion.User() != nil {
		return nil
	}

	u, err := m.db.FetchUser(m.client.Client, suggestion.UserId)
	if err != nil {
		return err
	}
	suggestion.SetUser(*u)

	return nil
}

func (m *SuggestionModule) DeployGuildSuggestCommand(guildId snowflake.ID, suggestionChannels []schema.SuggestionChannel) error {
	gcommands, err := m.client.Rest.GetGuildCommands(m.client.ID(), guildId, true)
	if err != nil {
		return err
	}

	if len(suggestionChannels) == 0 {
		if len(gcommands) == 0 {
			return nil
		}

		var commandId *snowflake.ID
		for _, gcommand := range gcommands {
			if gcommand.Name() == "suggest" {
				ci := gcommand.ID()
				commandId = &ci
			}
		}

		if commandId == nil {
			return nil
		}

		return m.client.Rest.DeleteGuildCommand(m.client.ID(), guildId, *commandId)
	}

	var subcommands []discord.ApplicationCommandOption
	for _, sc := range suggestionChannels {
		subcommands = append(subcommands, discord.ApplicationCommandOptionSubCommand{
			Name:        sc.Name,
			Description: fmt.Sprintf("The %s suggestion channel", sc.Name),
		})
	}

	var commandId *snowflake.ID
	for _, gcommand := range gcommands {
		if gcommand.Name() == "suggest" {
			ci := gcommand.ID()
			commandId = &ci
		}
	}

	if commandId != nil {
		_, err = m.client.Rest.UpdateGuildCommand(m.client.ID(), guildId, *commandId, discord.SlashCommandUpdate{
			Options: &subcommands,
		})
	} else {
		_, err = m.client.Rest.CreateGuildCommand(m.client.ID(), guildId, discord.SlashCommandCreate{
			Name:        "suggest",
			Description: "The suggest command",
			Options:     subcommands,
		})
	}

	return err
}

func (m *SuggestionModule) SendSuggestionChannelPanel(suggestionChannel *schema.SuggestionChannel) (*snowflake.ID, error) {
	message, err := m.client.Rest.CreateMessage(snowflake.MustParse(suggestionChannel.ChannelId), discord.MessageCreate{
		Content: "Welcome to the suggestion channel!",
		Embeds: []discord.Embed{
			MakeSuggestionChannelPanelEmbed(suggestionChannel),
		},
	})

	if err != nil {
		return nil, err
	}
	return &message.ID, nil
}

// EMBEDS

func (m *SuggestionModule) MakeSuggestionEmbed(suggestion *schema.Suggestion, suggestionChannel *schema.SuggestionChannel) discord.Embed {
	suggestionEmbed := discord.Embed{
		Description: suggestion.Description,
	}

	if suggestionChannel.EmbedColor != nil {
		suggestionEmbed.Color = *suggestionChannel.EmbedColor
	} else {
		suggestionEmbed.Color = util.EmbedColor
	}

	if suggestion.Title != nil {
		suggestionEmbed.Title = *suggestion.Title
		suggestionEmbed.URL = util.BotVotePage(m.config.WebsiteURL)
	}

	user := suggestion.User()
	if suggestion.AnonymousAuthor {
		suggestionEmbed.Author = &discord.EmbedAuthor{
			Name: "Anonymous suggested",
		}
		suggestionEmbed.Thumbnail = &discord.EmbedResource{
			URL: "https://cdn.discordapp.com/embed/avatars/0.png",
		}
	} else {
		suggestionEmbed.Author = &discord.EmbedAuthor{
			Name: user.SafeName() + " suggested",
		}
		suggestionEmbed.Thumbnail = &discord.EmbedResource{
			URL: "https://ava.viadev.xyz/" + user.UserId,
		}
		suggestionEmbed.Footer = &discord.EmbedFooter{
			Text: "@" + user.Username,
		}
	}

	if suggestion.ShowCounts {
		upvoteEmoji := util.EmojiString(util.UpvoteEmoji)
		downvoteEmoji := util.EmojiString(util.DownvoteEmoji)
		if suggestionChannel.UpvoteEmoji != nil {
			upvoteEmoji = *suggestionChannel.UpvoteEmoji
		}
		if suggestionChannel.DownvoteEmoji != nil {
			downvoteEmoji = *suggestionChannel.DownvoteEmoji
		}

		suggestionEmbed.Fields = []discord.EmbedField{
			{
				Name: "Votes",
				Value: fmt.Sprintf(
					"%s ` %d ` %s ` %d `",
					upvoteEmoji, len(suggestion.Upvotes),
					downvoteEmoji, len(suggestion.Downvotes),
				),
			},
		}
	}

	return suggestionEmbed
}

func (m *SuggestionModule) MakeProcessedSuggestionEmbed(suggestion *schema.Suggestion, suggestionChannel *schema.SuggestionChannel, approved bool) discord.Embed {
	suggestionEmbed := m.MakeSuggestionEmbed(suggestion, suggestionChannel)

	state := "[Approved]"
	if !approved {
		state = "[Denied]"
	}
	suggestionEmbed.Author.Name = fmt.Sprintf("%s\n%s", state, suggestionEmbed.Author.Name)

	if approved {
		suggestionEmbed.Color = util.GreenColor
	} else {
		suggestionEmbed.Color = util.RedColor
	}

	return suggestionEmbed
}

func MakeSuggestionChannelPanelEmbed(suggestionChannel *schema.SuggestionChannel) discord.Embed {
	panelEmbed := discord.Embed{}

	if suggestionChannel.Panel.Description != nil {
		panelEmbed.Description = *suggestionChannel.Panel.Description
	}

	if suggestionChannel.Panel.AuthorName != nil {
		panelEmbed.Author = &discord.EmbedAuthor{
			Name: *suggestionChannel.Panel.AuthorName,
		}
	}

	return panelEmbed
}

// COMPONENTS

func MakeSuggestionComponents(suggestion *schema.Suggestion) []discord.LayoutComponent {
	return []discord.LayoutComponent{
		discord.ActionRowComponent{
			Components: []discord.InteractiveComponent{
				discord.ButtonComponent{
					Label:    "0",
					CustomID: "dummy",
					Style:    discord.ButtonStyleSecondary,
					Disabled: true,
				},
				discord.ButtonComponent{
					Emoji:    &discord.ComponentEmoji{ID: util.UpvoteEmoji},
					CustomID: "suggestion:upvote",
					Style:    discord.ButtonStyleSuccess,
				},
				discord.ButtonComponent{
					Emoji:    &discord.ComponentEmoji{ID: util.DownvoteEmoji},
					CustomID: "suggestion:downvote",
					Style:    discord.ButtonStyleDanger,
				},
			},
		},
	}
}

// MISC

func SuggestionSubmitModal(channelId string) discord.ModalCreate {
	return discord.ModalCreate{
		Title:    "New Suggestion!",
		CustomID: "suggest:submit:" + channelId,
		Components: []discord.LayoutComponent{
			discord.LabelComponent{
				Label: "Title",
			},
			discord.ActionRowComponent{
				Components: []discord.InteractiveComponent{
					discord.TextInputComponent{
						CustomID: "title",
						Style:    discord.TextInputStyleShort,
					},
				},
			},
			discord.LabelComponent{
				Label: "Description",
			},
			discord.ActionRowComponent{
				Components: []discord.InteractiveComponent{
					discord.TextInputComponent{
						CustomID: "description",
						Style:    discord.TextInputStyleParagraph,
					},
				},
			},
		},
	}
}
