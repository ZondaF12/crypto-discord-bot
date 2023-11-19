package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ZondaF12/crypto-bot/cmd/request"
	"github.com/ZondaF12/crypto-bot/config"
	"github.com/bojanz/currency"
	"github.com/bwmarrin/discordgo"
	"github.com/leekchan/accounting"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "c",
			Description: "Command for demonstrating options",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "coin-symbol",
					Description: "Coin Symbol",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency",
					Description: "Currency",
					Required:    false,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"c": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response
			margs := make([]interface{}, 0, len(options))

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["coin-symbol"]; ok {
				// Option values must be type asserted from interface{}.
				// Discordgo provides utility functions to make this simple.
				margs = append(margs, option.StringValue())
			}

			if opt, ok := optionMap["currency"]; ok {
				margs = append(margs, opt.StringValue())
			} else {
				margs = append(margs, "GBP")
			}

			res := request.FetchPrice(margs[0].(string), strings.ToUpper(margs[1].(string)))

			locale := currency.NewLocale("en")
			symbol, ok := currency.GetSymbol(strings.ToUpper(margs[1].(string)), locale)
			if !ok {
				symbol = "£"
			}

			coinQuote := res.Data[strings.ToUpper(margs[0].(string))][0].Quote[strings.ToUpper(margs[1].(string))]
			ac := accounting.Accounting{Symbol: symbol, Precision: PriceRounding(coinQuote.Price)}

			embed := &discordgo.MessageEmbed{
				Type: discordgo.EmbedTypeRich,
				URL: "https://coinmarketcap.com/currencies/" + res.Data[strings.ToUpper(margs[0].(string))][0].Slug,
				Title: "Cryptocurrency Price Tracker",
				Description: res.Data[strings.ToUpper(margs[0].(string))][0].Name,
				Color: 16591219,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/128x128/%d.png", res.Data[strings.ToUpper(margs[0].(string))][0].Id),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "**Price (24h)**",
						Value: fmt.Sprintf("%s (%.2f%%)", ac.FormatMoney(coinQuote.Price), coinQuote.Percent_change_24h),
					},
					{
						Name: "**7 Day Percentage Change**",
						Value: fmt.Sprintf("%.2f%%", coinQuote.Percent_change_7d),
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Made by Roo#7777",
					IconURL: "https://i.ibb.co/VDMp2Bx/0e58a19b5a24f0542691313ff5106e40-1.png",
				},
				Timestamp: fmt.Sprintf("%v", time.Now().Format(time.RFC3339)),
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
		},
	}
)

func SetupDiscord(config config.EnvVars) error {
	s, err := discordgo.New("Bot " + config.TOKEN)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "704855585434239017", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "704855585434239017", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")

	return nil
}

func PriceRounding(price float32) int {
	if price < 1 {
		return 5
	} else if price < 10 {
		return 3
	}

	return 2
}