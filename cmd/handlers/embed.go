package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bojanz/currency"
	"github.com/bwmarrin/discordgo"
	"github.com/leekchan/accounting"
)

func GetDiscordEmbed(coin string, currencyCode string, quantity float32) *discordgo.MessageEmbed {
	locale := currency.NewLocale("en")
	symbol, ok := currency.GetSymbol(strings.ToUpper(currencyCode), locale)
	if !ok {
		symbol = "£"
	}

	res := GetCoinPrice(coin, currencyCode)
	coinQuote := res.Quote[strings.ToUpper(currencyCode)]

	var totalPrice float32
	if quantity != 0 {
		totalPrice = coinQuote.Price * quantity
	} else{
		totalPrice = coinQuote.Price
	}

	ac := accounting.Accounting{Symbol: symbol, Precision: PriceRounding(totalPrice)}

	if quantity != 0 {
		return &discordgo.MessageEmbed{
			Type: discordgo.EmbedTypeRich,
			URL: "https://coinmarketcap.com/currencies/" + res.Slug,
			Title: "Cryptocurrency Price Converter",
			Description: res.Name,
			Color: 16591219,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/128x128/%d.png", res.Id),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "**Price**",
					Value: ac.FormatMoney(totalPrice),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Made by Roo#7777",
				IconURL: "https://i.ibb.co/VDMp2Bx/0e58a19b5a24f0542691313ff5106e40-1.png",
			},
			Timestamp: fmt.Sprintf("%v", time.Now().Format(time.RFC3339)),
		}
	} else{
		return &discordgo.MessageEmbed{
			Type: discordgo.EmbedTypeRich,
			URL: "https://coinmarketcap.com/currencies/" + res.Slug,
			Title: "Cryptocurrency Price Tracker",
			Description: res.Name,
			Color: 16591219,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/128x128/%d.png", res.Id),
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
	}
}

func GetCoinPrice(coin string, currency string) CoinData {
	res := FetchPrice(coin, strings.ToUpper(currency))

	return res.Data[strings.ToUpper(coin)][0]
}

func PriceRounding(price float32) int {
	if price < 0.0001 {
		return 8
	} else if price < 1 {
		return 5
	}

	return 3
}