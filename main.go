package main

import "github.com/ZondaF12/crypto-bot/cmd"

func main() {
	// setup and run app
	err := cmd.Setup()
	if err != nil {
		panic(err)
	}
}