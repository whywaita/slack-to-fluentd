package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/nlopes/slack"
	"github.com/whywaita/slack_lib"
)

func structToJsonTagMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	return result
}

func main() {
	fluentHost := os.Getenv("FLUENT_HOST")
	if fluentHost == "" {
		fluentHost = "localhost"
	}
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK_TOKEN must be set")
	}

	flogger, err := fluent.New(fluent.Config{FluentHost: fluentHost})
	if err != nil {
		fmt.Println(err)
	}
	defer flogger.Close()

	tag := "slack.posts"
	api := slack.New(slackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		// fmt.Print("Event Received: ")
		err = nil
		switch ev := msg.Data.(type) {

		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			// fmt.Println("Infos:", ev.Info)
			// fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev.Msg)
			r, err := slack_lib.ConvertReadableName(api, ev)
			if err != nil {
				fmt.Println(err)
				continue
			}

			msg := structToJsonTagMap(r)
			err = flogger.Post(tag, msg)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case *slack.PresenceChangeEvent:
			// fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			// fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			// fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			// fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			fmt.Printf("Unexpected: %v\n", msg.Data)

		}
	}
}
