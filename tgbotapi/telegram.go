package tgbotapi

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MessengerService a telegram implementation of bbdcbot.MessengerService
//Uses telegram to converse with the user
type MessengerService struct {
	Bot    *tgbotapi.BotAPI
	ChatID int64
}

//NewMessengerService creates a new service, given a telegram bot API token
//and a telegram chat ID which indicates where messages will be sent
func NewMessengerService(token string, chatID int64) (*MessengerService, error) {
	ms := &MessengerService{
		ChatID: chatID,
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &MessengerService{}, errors.New("Failed to start telegram bot: " + err.Error())
	}

	ms.Bot = bot
	return ms, nil
}

//Alert sends a telegram message to the user with the ChatID
//in the struct
func (ms *MessengerService) Alert(msg string) {
	telegramMsg := tgbotapi.NewMessage(ms.ChatID, msg)
	ms.Bot.Send(telegramMsg)
}
