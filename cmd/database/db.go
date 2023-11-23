package database

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ZondaF12/crypto-bot/config"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func InitDB(env config.EnvVars) error {
	if env.MONGODB_URI == "" {
		return errors.New("you must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(env.MONGODB_URI))
	if err != nil {
		return err
	}

	db = client.Database("PriceAlerts")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}

type PriceAlert struct {
	Coin      string `json:"coin" bson:"coin"`
	GuildID   string `json:"GuildID" bson:"GuildID"`
	ChannelID string `json:"ChannelID" bson:"ChannelID"`
}

func GetPriceAlerts() []PriceAlert {
	coll := GetDBCollection("Price Alerts Col")

	// find all priceAlerts
	priceAlerts := make([]PriceAlert, 0)
	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	// iterate over the cursor
	for cursor.Next(context.Background()) {
		priceAlert := PriceAlert{}
		err := cursor.Decode(&priceAlert)
		if err != nil {
			fmt.Println(err)
		}
		priceAlerts = append(priceAlerts, priceAlert)
	}

	return priceAlerts
}

func CreatePriceAlert(options []*discordgo.ApplicationCommandInteractionDataOption, guildId string) error {
	// validate the body
	channelId := FilterChannelID(options[1].StringValue())
	newPriceAlert := PriceAlert{Coin: strings.ToUpper(options[0].StringValue()), GuildID: guildId, ChannelID: channelId}

	// create the price alert
	coll := GetDBCollection("Price Alerts Col")
	_, err := coll.InsertOne(context.Background(), newPriceAlert)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func RemovePriceAlert(options []*discordgo.ApplicationCommandInteractionDataOption, guildId string) int64 {
	coll := GetDBCollection("Price Alerts Col")

	channelId := FilterChannelID(options[1].StringValue())
	filter := bson.D{{"coin", strings.ToUpper(options[0].StringValue())}, {"ChannelID", channelId}, {"GuildID", guildId}}

	result, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
	}

	return result.DeletedCount
}

func FilterChannelID(channelId string) string {
	// Define a regular expression to match numbers
	re := regexp.MustCompile("[0-9]+")

	// Find all matches in the input string
	matches := re.FindAllString(channelId, -1)

	// Print the results
	var newChannelId string
	for _, match := range matches {
		newChannelId = match
	}

	return newChannelId
}
