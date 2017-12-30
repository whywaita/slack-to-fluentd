package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/nlopes/slack"
	"github.com/whywaita/slack_lib"
)

func structToJSONTagMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	typ := reflect.TypeOf(data)
	// should check data is struct or not
	val := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		vi := val.FieldByName(field.Name).Interface()
		// if field is struct, convert recursively
		if field.Type.Kind() == reflect.Struct {
			vi = structToJSONTagMap(vi)
		}
		if tag, ok := field.Tag.Lookup("json"); ok {
			result[tag] = vi
			continue
		}
		result[strings.ToLower(field.Name)] = vi
	}
	return result
}

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Fatal("SLACK_TOKEN must be set")
	}
	fluentHost := os.Getenv("FLUENTD_HOST")
	if fluentHost == "" {
		log.Fatal("FLUENTD_HOST must be set")
	}

	flogger, err := fluent.New(fluent.Config{FluentHost: fluentHost})
	if err != nil {
		log.Fatal(err)
	}
	defer flogger.Close()

	tag := "slack.post"
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

			msg := structToJSONTagMap(r)
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
			// fmt.Printf("Unexpected: %v\n", msg.Data)
			fmt.Printf("Unexpected: %v\n", msg)

		}
	}
}
