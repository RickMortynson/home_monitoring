package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home_services_analyst/internal/config"
	"io"
	"net/http"
)

type Telegram struct {
	chatId   int
	botToken string
}

func NewTelegramService(cfg *config.TelegramConfig) Telegram {
	return Telegram{
		chatId:   cfg.TgChatId,
		botToken: cfg.TgBotToken,
	}
}

type MessagePayload struct {
	ChatId              int    `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
	MessageId           int    `json:"message_id"`
}

type TgResponse struct {
	Result struct {
		MessageId int `json:"message_id"`
	} `json:"result"`
}

func (t *Telegram) SendMessage(message string) (TgResponse, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	payload := MessagePayload{
		ChatId: t.chatId,
		Text:   message,
	}

	marshalled, err := json.Marshal(payload)
	if err != nil {
		return TgResponse{}, err
	}

	resp, err := http.Post(url, CONTENT_TYPE_JSON, bytes.NewReader(marshalled))
	if err != nil {
		return TgResponse{}, err
	}

	respAsBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return TgResponse{}, err
	}

	var tgResponse TgResponse
	if err = json.Unmarshal(respAsBytes, &tgResponse); err != nil {
		return TgResponse{}, err
	}

	return tgResponse, nil
}

func (t *Telegram) UpdateTextMessage(messageId int, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", t.botToken)

	payload := MessagePayload{
		MessageId: messageId,
		ChatId:    t.chatId,
		Text:      message,
	}

	marshalled, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = http.Post(url, CONTENT_TYPE_JSON, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}

	return nil
}
