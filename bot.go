package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aikchun/gotelegrambot"
	"github.com/aikchun/gototo"
	"github.com/aikchun/totobro-bot/internal/services/subscriptionservice"
	"github.com/uptrace/bun"
)

type Bot struct {
	token               string
	db                  *bun.DB
	handlers            map[string]func(*gotelegrambot.Update, []string)
	subscriptionService *subscriptionservice.SubscriptionService
}

func newBot(t string, db *bun.DB) Bot {

	ss := subscriptionservice.NewSubscriptionService(db)
	b := Bot{
		token:               t,
		db:                  db,
		subscriptionService: &ss,
		handlers:            make(map[string]func(*gotelegrambot.Update, []string)),
	}
	b.handlers["/nextdraw"] = b.nextDraw
	b.handlers["/results"] = b.results
	b.handlers["/subscribe"] = b.subscribe
	b.handlers["/start"] = b.start
	b.handlers["/help"] = b.help

	return b
}

func (b Bot) handle(u *gotelegrambot.Update) {
	s := u.Message.Text

	trimmed := strings.Trim(s, " ")
	tokens := strings.Split(trimmed, " ")
	funcName := tokens[0]
	args := tokens[1:]

	for key, handler := range b.handlers {
		if key == funcName || key == fmt.Sprintf("%s@TotoBroBot", funcName) {
			handler(u, args)
		}
	}
}

func (b Bot) nextDraw(u *gotelegrambot.Update, args []string) {
	n := gototo.GetNextDraw()

	t := fmt.Sprintf("Date: %s\nPrize: %s", n.GetDate(), n.GetPrize())

	p := gotelegrambot.SendMessagePayload{
		ChatID: u.Message.Chat.ID,
		Text:   t,
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}

}

func (b Bot) results(u *gotelegrambot.Update, args []string) {
	d := gototo.GetLatestDraw()

	var ws []string

	for _, s := range d.GetWinningNumbers() {
		ws = append(ws, strconv.Itoa(s))
	}

	t := fmt.Sprintf("The latest Toto results:\nDate: %s\nWinning Numbers: %s\nAdditional Number: %d", d.GetDate(), strings.Join(ws, " "), d.GetAdditionalNumber())

	p := gotelegrambot.SendMessagePayload{
		ChatID: u.Message.Chat.ID,
		Text:   t,
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}

}

func (b Bot) start(u *gotelegrambot.Update, args []string) {
	p := gotelegrambot.SendMessagePayload{
		ChatID: u.Message.Chat.ID,
		Text:   fmt.Sprintf("%s\n%s", b.startMessage(), b.helpMessage()),
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}
}

func (b Bot) startMessage() string {
	return "Hello I am your bro for Toto."
}

func (b Bot) helpMessage() string {
	return "You can start with:\n/nextdraw to view the upcoming draw.\n/results for the latest draw results\n/subscribe to get toto draw alerts."
}

func (b Bot) help(u *gotelegrambot.Update, args []string) {
	p := gotelegrambot.SendMessagePayload{
		ChatID: u.Message.Chat.ID,
		Text:   b.helpMessage(),
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}
}

func (b Bot) subscribe(u *gotelegrambot.Update, args []string) {
	chatID := u.Message.Chat.ID

	if _, err := b.subscriptionService.Subscribe(chatID, 0); err != nil {
		log.Fatal(err)
	}

	p := gotelegrambot.SendMessagePayload{
		ChatID: chatID,
		Text:   "You are subscribed to next draw alerts for Toto!",
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}

}
