package command

import (
	"fmt"
	"strings"

	"github.com/ZondaF12/crypto-bot/cmd/database"
	"github.com/bwmarrin/discordgo"
)

func FollowCrypto(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	err := database.CreatePriceAlert(options, i.GuildID)
	if err != nil {
		fmt.Println(err)
	}

	content := fmt.Sprintf("Price Alert for %s scheduled in %s", strings.ToUpper(options[0].StringValue()), options[1].StringValue())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
