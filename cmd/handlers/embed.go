package handlers

import (
	"fmt"
	"time"

	"github.com/bojanz/currency"
	"github.com/bwmarrin/discordgo"
	"github.com/leekchan/accounting"
)

func GetDiscordEmbed(coin string, currencyCode string, quantity float32, priceAlert bool) *discordgo.MessageEmbed {
	locale := currency.NewLocale("en")
	symbol, ok := currency.GetSymbol(currencyCode, locale)
	if !ok {
		symbol = "£"
	}

	res := GetCoinPrice(coin, currencyCode)
	coinQuote := res.Quote[currencyCode]

	var totalPrice float32
	if quantity != 0 {
		totalPrice = coinQuote.Price * quantity
	} else {
		totalPrice = coinQuote.Price
	}

	ac := accounting.Accounting{Symbol: symbol, Precision: PriceRounding(totalPrice)}

	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		URL:         "https://coinmarketcap.com/currencies/" + res.Slug,
		Title:       "Cryptocurrency Price Tracker",
		Description: res.Name,
		Color:       16591219,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/128x128/%d.png", res.Id),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Price (24h)**",
				Value: fmt.Sprintf("%s (%.2f%%)", ac.FormatMoney(totalPrice), coinQuote.Percent_change_24h),
			},
			{
				Name:  "**7 Day Percentage Change**",
				Value: fmt.Sprintf("%.2f%%", coinQuote.Percent_change_7d),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Made by Roo#7777",
			IconURL: "https://i.ibb.co/VDMp2Bx/0e58a19b5a24f0542691313ff5106e40-1.png",
		},
		Timestamp: fmt.Sprintf("%v", time.Now().Format(time.RFC3339)),
	}

	if quantity != 0 && !priceAlert {
		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "**Price**",
				Value: ac.FormatMoney(totalPrice),
			},
		}

		embed.Title = "Cryptocurrency Price Converter"
	}

	if priceAlert {
		embed = &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("%s Price Update", res.Name),
				IconURL: fmt.Sprintf("https://s2.coinmarketcap.com/static/img/coins/128x128/%d.png", res.Id),
			},
			Description: fmt.Sprintf("The current price of %s is **%s**", coin, ac.FormatMoney(totalPrice)),
			Color:       16591219,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Made by Roo#7777",
				IconURL: "https://i.ibb.co/VDMp2Bx/0e58a19b5a24f0542691313ff5106e40-1.png",
			},
			Timestamp: fmt.Sprintf("%v", time.Now().Format(time.RFC3339)),
		}
	}

	return embed
}

func GetCoinPrice(coin string, currency string) CoinData {
	res := FetchPrice(coin, currency)

	return res.Data[coin][0]
}

func PriceRounding(price float32) int {
	if price < 0.0001 {
		return 8
	} else if price < 1 {
		return 5
	}

	return 2
}
