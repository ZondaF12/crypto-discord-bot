package command

import (
	"strconv"

	"github.com/ZondaF12/crypto-bot/cmd/handlers"
	"github.com/bwmarrin/discordgo"
)

func FollowCrypto(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	margs := make([]interface{}, 0, len(options))

	if option, ok := optionMap["coin-symbol"]; ok {

		margs = append(margs, option.StringValue())
	}

	var quantity float32
	if option, ok := optionMap["quantity"]; ok {
		margs = append(margs, option.StringValue())

		// "var float float32" up here somewhere
		value, _ := strconv.ParseFloat(margs[1].(string), 32)

		quantity = float32(value)
	}

	if opt, ok := optionMap["currency"]; ok {
		margs = append(margs, opt.StringValue())
	} else {
		margs = append(margs, "GBP")
	}

	embed := handlers.GetDiscordEmbed(margs[0].(string), margs[2].(string), quantity)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
