package command

import (
	"strings"

	"github.com/ZondaF12/crypto-bot/cmd/handlers"
	"github.com/bwmarrin/discordgo"
)

func CheckCryptoPrice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	var currency string
	if len(options) == 1 {
		currency = "GBP"
	} else {
		currency = strings.ToUpper(options[1].StringValue())
	}

	embed := handlers.GetDiscordEmbed(strings.ToUpper(options[0].StringValue()), currency, 0)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
