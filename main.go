package main

import (
	"log"
)

func main() {
	botToken := "7326700702:AAGYArp6trbneAXPngI4fJApx35KxwdFNPM"
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	userTracker := NewSimpleUserTracker()
	contentModerator := NewSimpleDrugModerator()

	bot, err := NewBot(botToken, userTracker, contentModerator)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Start()
}
