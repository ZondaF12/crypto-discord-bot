package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ZondaF12/crypto-bot/cmd/handlers"
	"github.com/bwmarrin/discordgo"
)

func ConvertCrypto(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	value, err := strconv.ParseFloat(options[1].StringValue(), 32)
	if err != nil {
		fmt.Println(err)
	}
	quantity := float32(value)

	var currency string
	if len(options) == 2 {
		currency = "GBP"
	} else {
		currency = strings.ToUpper(options[2].StringValue())
	}

	embed := handlers.GetDiscordEmbed(strings.ToUpper(options[0].StringValue()), currency, quantity, false)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
