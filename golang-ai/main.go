package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Krognol/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go"
)

var wolframClient *wolfram.Client

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

	wolframClient = &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}
	go printCommandEvents(bot.CommandEvents())

	bot.Command("ask-wolfram <message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Examples:    []string{"who is the prime minister of india"},
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			query := r.Param("message")

			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data[:]) // convert whatever data i am getting in string format
			fmt.Println("answer", rough)
			value := gjson.Get(rough, "entities.wolfram_search_query.0.value") // .0 to access the value of first object
			answer := value.String()

			if answer == "" {
				w.Reply("I'm sorry, I couldn't understand the question.")
				return
			}

			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("there is some error", err)
			}
			fmt.Println("ANSWRR VALUE IS", answer)
			fmt.Println(res)
			fmt.Println(value)
			fmt.Println(query)
			fmt.Println(msg)
			fmt.Println(rough)
			w.Reply(res)
		},
	})

	// stop the program
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}
