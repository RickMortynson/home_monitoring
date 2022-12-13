package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home_services_analyst/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

const (
	CONTENT_TYPE_JSON = "application/json"
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
}

type UpdateMessagePayload struct {
	MessagePayload

	MessageId int `json:"message_id"`
}

type ImageMessagePayload struct {
	MessagePayload

	Photo   string `json:"photo"`
	Caption string `json:"caption"`
}

type TgResponse struct {
	Result struct {
		MessageId int `json:"message_id"`
	} `json:"result"`
}

func (t *Telegram) SendTextMessage(message string) (TgResponse, error) {
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

func (t *Telegram) SendImageMessage(image []byte, caption string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", t.botToken)

	var payloadBody bytes.Buffer
	mp := multipart.NewWriter(&payloadBody)

	values := map[string]io.Reader{
		"photo":   bytes.NewReader(image),
		"caption": strings.NewReader(caption),
		"chat_id": strings.NewReader(strconv.Itoa(t.chatId)),
	}

	for key, r := range values {
		if _, ok := r.(*bytes.Reader); ok {
			part, err := mp.CreateFormFile(key, "image.png")
			if err != nil {
				return err
			}
			io.Copy(part, r)
		} else {
			ff, err := mp.CreateFormField(key)
			if err != nil {
				return err
			}

			io.Copy(ff, r)
		}
	}

	if err := mp.Close(); err != nil {
		return err
	}

	if _, err := http.Post(url, mp.FormDataContentType(), &payloadBody); err != nil {
		return err
	}

	return nil
}

func (t *Telegram) UpdateTextMessage(messageId int, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", t.botToken)

	payload := UpdateMessagePayload{
		MessageId: messageId,
		MessagePayload: MessagePayload{
			ChatId: t.chatId,
			Text:   message,
		},
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
