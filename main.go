package main

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	botToken    = "YOUR_BOT_TOKEN_HERE"
	banDuration = 5 * time.Minute
)

var newUsers = make(map[int64]time.Time)

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleNewMembers(update.Message)

		if update.Message.Text != "" {
			handleMessage(bot, update.Message)
		}
	}
}

func handleNewMembers(message *tgbotapi.Message) {
	if message.NewChatMembers != nil {
		for _, newUser := range message.NewChatMembers {
			newUsers[newUser.ID] = time.Now()
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID

	joinTime, isNewUser := newUsers[userID]
	if !isNewUser {
		return
	}

	if time.Since(joinTime) > banDuration {
		delete(newUsers, userID)
		return
	}

	if containsDrugRelatedContent(message.Text) {
		if err := banUser(bot, chatID, userID); err != nil {
			log.Printf("Error banning user: %v", err)
			return
		}
		delete(newUsers, userID)
	}
}

func containsDrugRelatedContent(text string) bool {
	drugKeywords := []string{"drugs", "cocaine", "heroin", "weed", "meth"}
	lowercaseText := strings.ToLower(text)

	for _, keyword := range drugKeywords {
		if strings.Contains(lowercaseText, keyword) {
			return true
		}
	}
	return false
}

func banUser(bot *tgbotapi.BotAPI, chatID, userID int64) error {
	banConfig := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: time.Now().Add(24 * time.Hour).Unix(),
	}

	if _, err := bot.Request(banConfig); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "User banned for attempting to sell drugs.")
	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
