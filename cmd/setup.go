package cmd

import (
	"fmt"
	"log"

	"github.com/ZondaF12/crypto-bot/cmd/database"
	"github.com/ZondaF12/crypto-bot/config"
)

func Setup() error {
	// Load ENV Variables from app.env
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	// Init DB
	err = database.InitDB(env)
	if err != nil {
		return err
	}

	// Defer closing the database
	defer database.CloseDB()

	// Init Discord
	err = SetupDiscord(env)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return err
	}

	return nil
}
