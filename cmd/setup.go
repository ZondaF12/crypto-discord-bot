package cmd

import (
	"fmt"
	"log"

	"github.com/ZondaF12/crypto-bot/config"
)

func Setup() error {
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	err = SetupDiscord(env)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return err
	}
	
	return nil
}