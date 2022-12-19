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

	for {
		mailClient := mail.NewClient(conf.ImapHost, conf.ImapUser, conf.ImapPassword)

		forwardMsg := func(msg *mail.Message) {
			text := fmt.Sprintf("*%s*:\n%s", msg.Subject, msg.Body)
			telegramAPI.SendMessage(text)
			mailClient.MarkMsgSeen(msg)
		}

		for _, msg := range mailClient.FetchUnseenMsgs() {
			forwardMsg(msg)
		}

		mailClient.WaitNewMsgs(conf.ImapPollInterval)
		mailClient.Logout()
	}
}
