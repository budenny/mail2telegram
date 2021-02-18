package conf

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Conf represents application configuration
type Conf struct {
	TelegramChatID   int64
	TelegramBotToken string

	ImapHost         string
	ImapUser         string
	ImapPassword     string
	ImapPollInterval time.Duration
}

// BuildConf builds conf from environment
func BuildConf() *Conf {
	c := &Conf{}
	c.TelegramChatID = getEnvVarInt("TELEGRAM_CHAT_ID")
	c.TelegramBotToken = getEnvVarStr("TELEGRAM_BOT_TOKEN")
	c.ImapHost = getEnvVarStr("MAIL_IMAP_HOST")
	c.ImapUser = getEnvVarStr("MAIL_USER")
	c.ImapPassword = getEnvVarStr("MAIL_PASSWORD")
	c.ImapPollInterval = time.Duration(getEnvVarInt("MAIL_POLL_INTERVAL_SEC")) * time.Second
	return c
}

func getEnvVarStr(name string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		log.Fatal("missing " + name + " var")
	}
	return value
}

func getEnvVarInt(name string) int64 {
	valueStr := getEnvVarStr(name)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
