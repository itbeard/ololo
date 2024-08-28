package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	banDuration = 5 * time.Minute
	banMessage  = "User banned for attempting to sell drugs."
)

type Bot struct {
	api              *tgbotapi.BotAPI
	userTracker      UserTracker
	contentModerator ContentModerator
}

func NewBot(token string, tracker UserTracker, moderator ContentModerator) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, userTracker: tracker, contentModerator: moderator}, nil
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.handleNewMembers(update.Message)

		if update.Message.Text != "" {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleNewMembers(message *tgbotapi.Message) {
	if message.NewChatMembers != nil {
		for _, newUser := range message.NewChatMembers {
			b.userTracker.AddUser(newUser.ID)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID

	if !b.userTracker.IsNewUser(userID) {
		return
	}

	if b.contentModerator.IsForbiddenContent(message.Text) {
		if err := b.banUser(chatID, userID); err != nil {
			log.Printf("Error banning user: %v", err)
			return
		}
		b.userTracker.RemoveUser(userID)
	}
}

func (b *Bot) banUser(chatID, userID int64) error {
	banConfig := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: time.Now().Add(24 * time.Hour).Unix(),
	}

	if _, err := b.api.Request(banConfig); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, banMessage)
	if _, err := b.api.Send(msg); err != nil {
		return err
	}

	return nil
}
