package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aikchun/gotelegrambot"
	"github.com/aikchun/totobro-bot/internal/services/db"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func handleLambdaEvent(u gotelegrambot.Update) {
	d := db.NewDB()
	b := newBot(os.Getenv("BOT_USERNAME"), os.Getenv("BOT_TOKEN"), d)
	b.handle(&u)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var u gotelegrambot.Update

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &u); err != nil {
		log.Fatal(err)
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
