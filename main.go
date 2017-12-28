package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/jinzhu/copier"
	"github.com/nlopes/slack"
	"github.com/whywaita/slack_lib"
)

func StructToJsonTagMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	return result
}

func convertReadableName(api *slack.Client, ev *slack.MessageEvent) (slack.Msg, error) {
	var err error
	result := slack.Msg{}
	msg := ev.Msg

	copier.Copy(&result, &msg)

	rUser, err := api.GetUserInfo(msg.User)
	if err != nil {
		return slack.Msg{}, err
	}

	_, channelName, err := slack_lib.GetFromName(api, ev)
	if err != nil {
		return slack.Msg{}, err
	}

	rTeam, err := api.GetTeamInfo()
	if err != nil {
		return slack.Msg{}, err
	}

	result.User = rUser.Name
	result.Channel = channelName
	result.Team = rTeam.Name

	return result, nil
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
		switch ev := msg.Data.(type) {

		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			// fmt.Println("Infos:", ev.Info)
			// fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev.Msg)
			r, err := convertReadableName(api, ev)
			if err != nil {
				fmt.Println(err)
				continue
			}

			msg := StructToJsonTagMap(r)
			err = flogger.Post(tag, msg)
			if err != nil {
				fmt.Println(err)
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
