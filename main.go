package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aikchun/gotelegrambot"
	"github.com/aikchun/gototo"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func nextDraw(bot *gotelegrambot.Bot, u *gotelegrambot.Update, args []string) {
	n := gototo.GetNextDraw()

	t := fmt.Sprintf("Date: %s\nPrize: %s", n.GetDate(), n.GetPrize())

	p := gotelegrambot.SendMessagePayload{
		ChatID: u.Message.Chat.ID,
		Text:   t,
	}

	if _, err := bot.SendMessage(p); err != nil {
		log.Fatal(err)
	}

}

func results(bot *gotelegrambot.Bot, u *gotelegrambot.Update, args []string) {
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

	if _, err := bot.SendMessage(p); err != nil {
		log.Fatal(err)
	}

}

func handleLambdaEvent(u gotelegrambot.Update) {
	bot, err := gotelegrambot.NewBot(os.Getenv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	bot.SetUpdateHandler("/nextdraw", nextDraw)
	bot.SetUpdateHandler("/results", results)

	bot.HandleUpdate(&u)

}

func handler(w http.ResponseWriter, r *http.Request) {
	var u gotelegrambot.Update

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &u); err != nil {
		panic(err)
	}

	handleLambdaEvent(u)
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Couldn't find .env")
	}

	e := os.Getenv("ENVIRONMENT")

	if e != "dev" {
		lambda.Start(handleLambdaEvent)
	}

	p := ":8080"

	if e == "dev" {
		http.HandleFunc("/", handler)
		fmt.Printf("starting http server\n")
		fmt.Printf("listening: %s\n", p)

		log.Fatal(http.ListenAndServe(p, nil))
	}
}
