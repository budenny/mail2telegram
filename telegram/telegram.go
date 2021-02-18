package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// API provides basic telegram functionality
type API struct {
	chatID  int64
	sendURI string
}

//NewAPI instantiates api
func NewAPI(chatID int64, botToken string) *API {
	sendURI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	log.Printf("Telegram uri is:" + sendURI)
	return &API{chatID, sendURI}
}

type reqBody struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// SendMessage ...
func (api *API) SendMessage(text string) {
	req := &reqBody{
		ChatID:    api.chatID,
		Text:      text,
		ParseMode: "markdown",
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post(api.sendURI, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("unexpected status" + res.Status)
	}

	log.Printf("Sent Message: '%s'\n", text)
}
