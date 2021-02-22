package main

import (
	"fmt"

	"mail2telegram/conf"
	"mail2telegram/mail"
	"mail2telegram/telegram"
)

func main() {
	conf := conf.BuildConf()

	telegramAPI := telegram.NewAPI(conf.TelegramChatID, conf.TelegramBotToken)

	mailClient := mail.NewClient(conf.ImapHost, conf.ImapUser, conf.ImapPassword)
	defer mailClient.Logout()

	forwardMsg := func(msg *mail.Message) {
		text := fmt.Sprintf("*%s*:\n%s", msg.Subject, msg.Body)
		telegramAPI.SendMessage(text)
		mailClient.MarkMsgSeen(msg)
	}

	for _, msg := range mailClient.FetchUnseenMsgs() {
		forwardMsg(msg)
	}

	msgs := make(chan *mail.Message)

	go func() {
		for msg := range msgs {
			forwardMsg(msg)
		}
	}()

	mailClient.WaitNewMsgs(msgs, conf.ImapPollInterval)
}
