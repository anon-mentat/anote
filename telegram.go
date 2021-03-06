package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ui18n "github.com/unknwon/i18n"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Telegram group ID consts
const (
	tAnonBalkan       = -1001161265502
	tAnon             = -1001361489843
	tAnonTaxi         = -1001422544298
	tAnonTaxiPrv      = -1001271198034
	tAnonOps          = -1001213539865
	tAnonShout        = -1001453693349
	tAnonShoutPreview = -1001484971271
)

func initBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(conf.TelegramAPIKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(tAnonOps, "Anote daemon successfully started. 🚀")
	bot.Send(msg)

	return bot
}

func logTelegram(message string) {
	msg := tgbotapi.NewMessage(tAnonOps, message)
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}

func messageTelegram(message string, groupID int64) error {
	message = strings.Replace(message, "\\n", "\n", -1)
	msg := tgbotapi.NewMessage(groupID, message)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	if err != nil &&
		!strings.Contains(err.Error(), "blocked") &&
		!strings.Contains(err.Error(), "chat not found") &&
		!strings.Contains(err.Error(), "initiate") &&
		!strings.Contains(err.Error(), "deactivated") {
		logTelegram("[telegram.go - 50]" + err.Error() + " ### user: " + strconv.Itoa(int(groupID)) + " ### message: " + message)
	}
	return err
}

func sendInvestmentMessages(investment float64, newPrice float64) {
	msg := fmt.Sprintf(ui18n.Tr(lang, "newPurchase"), investment)
	msgHr := fmt.Sprintf(ui18n.Tr(langHr, "newPurchase"), investment)

	if newPrice > float64(0) {
		msg += "\n\n"
		msg += fmt.Sprintf(ui18n.Tr(lang, "priceRise"), newPrice)

		msgHr += "\n\n"
		msgHr += fmt.Sprintf(ui18n.Tr(langHr, "priceRise"), newPrice)
	}

	msg += "\n\n"
	msg += ui18n.Tr(lang, "purchaseHowto")

	msgHr += "\n\n"
	msgHr += ui18n.Tr(langHr, "purchaseHowto")

	messageTelegram(msg, tAnonOps)
	messageTelegram(msg, tAnon)
	messageTelegram(msgHr, tAnonBalkan)
}

// TelegramUpdate struct represent webhook update data from Telegram
type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID                          int    `json:"id"`
			Title                       string `json:"title"`
			Type                        string `json:"type"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
		} `json:"chat"`
		Date           int `json:"date"`
		ReplyToMessage struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID                          int    `json:"id"`
				Title                       string `json:"title"`
				Type                        string `json:"type"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"reply_to_message"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
		NewChatParticipant struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"new_chat_participant"`
		NewChatMember struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"new_chat_member"`
		NewChatMembers []struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"new_chat_members"`
	} `json:"message"`
}

func clean() {
	var users []*User
	db.Find(&users)
	for i, u := range users {
		if u.MiningActivated != nil {
			err := messageTelegram("Miner check.", int64(u.TelegramID))
			if err != nil &&
				(strings.Contains(err.Error(), "blocked") ||
					strings.Contains(err.Error(), "chat not found") ||
					strings.Contains(err.Error(), "initiate") ||
					strings.Contains(err.Error(), "deactivated")) {
				db.Delete(u)
			} else if err != nil {
				logTelegram("[telegram.go - 160]" + err.Error())
			}
		}
		log.Println(i)
	}
	log.Println("done cleaning")
}

func clean1() {
	var users []*User
	db.Unscoped().Find(&users)
	for i, u := range users {
		if len(u.Address) == 0 {
			u.Address = u.Nickname
			db.Unscoped().Save(u)
		}
		log.Println(i)
	}
	log.Println("done cleaning")
}
