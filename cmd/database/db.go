package database

import (
	"context"
	"errors"
	"fmt"

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

func CreatePriceAlert(options []*discordgo.ApplicationCommandInteractionDataOption, guildId string) error {
	// validate the body
	newPriceAlert := PriceAlert{Coin: options[0].StringValue(), GuildID: guildId, ChannelID: options[1].StringValue()}

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
	filter := bson.D{{"coin", options[0].StringValue()}, {"ChannelID", options[1].StringValue()}, {"GuildID", guildId}}

	result, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
	}

	return result.DeletedCount
}
