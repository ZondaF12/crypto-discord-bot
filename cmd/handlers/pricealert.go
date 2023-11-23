package handlers

import (
	"strings"
	"time"

	"github.com/ZondaF12/crypto-bot/cmd/database"
	"github.com/bwmarrin/discordgo"
)

func RunPriceAlert(s *discordgo.Session) {
	HandlePriceAlerts(s)
	
	priceAlerts := time.NewTicker(4 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
		select {
			case <- priceAlerts.C:
				HandlePriceAlerts(s)
			case <- quit:
				priceAlerts.Stop()
				return
			}
		}
	}()
}

func HandlePriceAlerts(s *discordgo.Session)  {
	result := database.GetPriceAlerts()

	for _, alert := range result {
		embed := GetDiscordEmbed(strings.ToUpper(alert.Coin), "GBP", 0, true)
		s.ChannelMessageSendEmbed(alert.ChannelID, embed)
	}
}