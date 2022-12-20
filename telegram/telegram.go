package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func ensureTruncated(s string, maxLen int) string {
	if len(s) >= maxLen {
		delim := "\n\n...\n\n"
		nDelim := len(delim)
		return s[:maxLen/2-nDelim] + delim + s[len(s)-maxLen/2+nDelim:]
	}
	return s
}

// SendMessage ...
func (api *API) SendMessage(text string) {

	// odd numbers of "_" break Telegram parser in "markdown" mode
	// so let's workarund it
	text = strings.Replace(text, "_", "\\_", -1)

	// Telegram API doesn't allow messages longer than 4000 chars
	text = ensureTruncated(text, 4000)

	req := &reqBody{
		ChatID:    api.chatID,
		Text:      text,
		ParseMode: "markdown",
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sending Message: '%s'\n", text)

	res, err := http.Post(api.sendURI, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		bodyStr := ""
		if err == nil {
			bodyStr = string(body)
		}
		log.Fatalf("Send Message: unexpected status: %v\nResp:%v", res.Status, bodyStr)
	}

	log.Println("Message sent")
}
