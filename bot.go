package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aikchun/gotelegrambot"
	"github.com/aikchun/gototo"
	"github.com/aikchun/totobro-bot/internal/services/nextdrawservice"
	"github.com/aikchun/totobro-bot/internal/services/subscriptionservice"
	"github.com/uptrace/bun"
)

type Bot struct {
	token               string
	db                  *bun.DB
	handlers            map[string]func(*gotelegrambot.Update, []string)
	subscriptionService *subscriptionservice.SubscriptionService
	nextDrawService     *nextdrawservice.NextDrawService
}

func newBot(t string, db *bun.DB) Bot {

	ss := subscriptionservice.NewSubscriptionService(db)
	nds := nextdrawservice.NewNextDrawService(db)
	b := Bot{
		token:               t,
		db:                  db,
		subscriptionService: &ss,
		nextDrawService:     &nds,
		handlers:            make(map[string]func(*gotelegrambot.Update, []string)),
	}
	b.handlers["/nextdraw"] = b.nextDraw
	b.handlers["/results"] = b.results
	b.handlers["/subscribe"] = b.subscribe
	b.handlers["/start"] = b.start
	b.handlers["/help"] = b.help
	b.handlers["/fetchNextDraw"] = b.fetchNextDraw
	b.handlers["/unsubscribe"] = b.unsubscribe

	return b
}

func (b Bot) handle(u *gotelegrambot.Update) {
	s := u.Message.Text

	trimmed := strings.Trim(s, " ")
	tokens := strings.Split(trimmed, " ")
	funcName := tokens[0]
	args := tokens[1:]

	for key, handler := range b.handlers {
		if key == funcName || funcName == fmt.Sprintf("%s@TotoBroBot", key) {
			handler(u, args)
		}
	}
}

func (b Bot) nextDraw(u *gotelegrambot.Update, args []string) {
	n := gototo.GetNextDraw()

	p := gotelegrambot.SendMessagePayload{
		ChatID: u.Message.Chat.ID,
		Text:   nextDrawMessage(n),
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}

}

func nextDrawMessage(n gototo.NextDraw) string {
	isToday := isNextDrawToday(n)

	if isToday {
		return fmt.Sprintf("The draw is today!\n\nDate: %s\nPrize: %s", n.GetDate(), n.GetPrize())
	}

	return fmt.Sprintf("Next Toto Draw\n\nDate: %s\nPrize: %s", n.GetDate(), n.GetPrize())
}

func (b Bot) results(u *gotelegrambot.Update, args []string) {
	d := gototo.GetLatestDraw()

	var ws []string

	for _, s := range d.GetWinningNumbers() {
		ws = append(ws, strconv.Itoa(s))
	}

	text := `The latest Toto results

	Date: %s
	Winning Numbers: %s
	Additional Number: %d`

	t := fmt.Sprintf(text, d.GetDate(), strings.Join(ws, " "), d.GetAdditionalNumber())

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
	return `You can start with:
	/nextdraw to view the upcoming draw.
	/results for the latest draw results
	/subscribe to get toto draw alerts.
	/unsubscribe to stop getting toto draw alerts.`
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
	p := gotelegrambot.SendMessagePayload{
		ChatID: chatID,
		Text:   "You are subscribed to next draw alerts for Toto!",
	}

	if _, err := b.subscriptionService.Subscribe(chatID, 0); err != nil {
		log.Println(err)
		p.Text = "subscribe unsuccessful"
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}

}

func (b Bot) unsubscribe(u *gotelegrambot.Update, args []string) {
	chatID := u.Message.Chat.ID
	p := gotelegrambot.SendMessagePayload{
		ChatID: chatID,
		Text:   "You have unsubscribed from new Toto draw alerts!",
	}

	if err := b.subscriptionService.SoftDelete(chatID); err != nil {
		log.Println(err)
		p.Text = "unsubscribe unsuccessful"
	}

	if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
		log.Fatal(err)
	}

}

func (b Bot) fetchNextDraw(u *gotelegrambot.Update, args []string) {
	for _, p := range b.createNextDrawBroadcastPayloads(u, gototo.GetNextDraw()) {
		if _, err := gotelegrambot.SendMessage(b.token, p); err != nil {
			log.Fatalln(err)
		}
	}

}

func (b Bot) createNextDrawBroadcastPayloads(u *gotelegrambot.Update, n gototo.NextDraw) []gotelegrambot.SendMessagePayload {

	if !isValidFetchNextDrawUpdate(u) {
		log.Fatalln("Not valid payload")
	}

	subs, err := b.subscriptionService.GetAll()

	if err != nil {
		log.Fatalln(err)
	}

	_, err = b.nextDrawService.Save(n)

	isToday := isNextDrawToday(n)

	if err != nil && !isToday {
		log.Fatalln("failed to save next draw and next draw is not today")
	}

	var ps []gotelegrambot.SendMessagePayload

	for _, s := range subs {
		p := gotelegrambot.SendMessagePayload{
			ChatID: s.ChatID,
			Text:   nextDrawMessage(n),
		}
		ps = append(ps, p)
	}

	return ps

}

func isValidFetchNextDrawUpdate(u *gotelegrambot.Update) bool {
	id, err := strconv.ParseInt(os.Getenv("FETCH_NEXT_DRAW_TASK_MOCK_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// we use one super large value for id for all the objects
	// to make it incredibly unlikely for anyone but us to make this call

	return u.UpdateID == id &&
		u.Message.MessageID == id &&
		u.Message.From.ID == id &&
		u.Message.Chat.ID == id &&
		u.Message.From.FirstName == os.Getenv("FETCH_NEXT_DRAW_TASK_USER_FIRST_NAME")
}

func isNextDrawToday(n gototo.NextDraw) bool {
	now := time.Now().Add(time.Hour * 8)
	nowString := now.Format("Mon, 02 Jan 2006")

	return nowString == parseNextDrawDateString(n.GetDate())
}

func parseNextDrawDateString(dateString string) string {
	s := strings.Split(dateString, " ")

	var ds []string

	for i := 0; i < 4; i++ {
		ds = append(ds, s[i])
	}

	return strings.Join(ds, " ")
}
