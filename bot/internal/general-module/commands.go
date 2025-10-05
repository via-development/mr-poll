package generalModule

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (m *GeneralModule) MrPollCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return interaction.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			IntroductoryEmbed(),
		},
		Components: []discord.ContainerComponent{
			HelpComponents(true),
		},
	})
}

func (m *GeneralModule) MyTimeCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
	//userData, err := db.FetchUser(interaction.Client(), interaction.User().ID.String())
	//if err != nil {
	//	return err
	//}
	//
	//cmdData := interaction.SlashCommandInteractionData()
	//offset := float32(cmdData.Float("utc-offset"))
	//
	//userData.UTCOffset = &offset
	//format := cmdData.String("date-format")
	//switch format {
	//case "ddmmyy":
	//	f := dateformat.DDMMYY
	//	userData.DateFormat = &f
	//case "mmddyy":
	//	f := dateformat.MMDDYY
	//	userData.DateFormat = &f
	//case "yymmdd":
	//	f := dateformat.YYMMDD
	//	userData.DateFormat = &f
	//default:
	//	return interaction.CreateMessage(discord.MessageCreate{
	//		Embeds: []discord.Embed{util.MakeErrorEmbed("You didn't input a valid dateformat.")},
	//	})
	//}
	//
	//err = db.Save(&userData).Error
	//if err != nil {
	//	return err
	//}
	//
	//return interaction.CreateMessage(discord.MessageCreate{
	//	Embeds: []discord.Embed{util.MakeSuccessEmbed("Your timezone and preferred date format was set!\n-# Please make sure to update your timezone if it changes")},
	//})
}
