package main

import (
	"github.com/juju/loggo"
	"github.com/nlopes/slack"
	"reflect"
	"slack-bot/workers"
)

var rootLogger = loggo.GetLogger("")
var log = loggo.GetLogger("bot")

func init() {
	rootLogger.SetLogLevel(loggo.INFO)
	log.SetLogLevel(loggo.TRACE)
}

func main() {
	config := configConstructor()

	// Slack Client
	api := slack.New(config.Slack.Token)
	log.Infof("Start slack event listening...")
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Creo e inicio los Workers
	var wrkrs []workers.Worker
	for  i, w := range config.Workers {
		wrkrs = append(wrkrs, workers.Constructor(config.Workers[i].Type, rtm))
		wrkrs[i].Setter(w.Params)
		wrkrs[i].Start()
	}

	// NOTE: https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel
	for {
		select {
		case msg := <- rtm.IncomingEvents :
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:

				// Exclude himself and "Slackbot" from response
				config.Slack.ExcludedUsers = append(config.Slack.ExcludedUsers, ev.Info.User.ID)
				config.Slack.ExcludedUsers = append(config.Slack.ExcludedUsers, "USLACKBOT")

				log.Infof("Connection counter: %v", ev.ConnectionCount)
			case *slack.MessageEvent:
				if !workers.StringMatch(msg.Data.(*slack.MessageEvent).User, config.Slack.ExcludedUsers ) {
					log.Debugf("Broadcasting message from %v to workers", msg.Data.(*slack.MessageEvent).User)
					for _, w := range wrkrs {
						reflect.ValueOf(w).Elem().FieldByName("Inbox").Send(reflect.ValueOf(msg))
					}
				} else {
					log.Debugf("Message from %s has been excluded by global match", msg.Data.(*slack.MessageEvent).User)
				}
			case *slack.RTMError:
				log.Errorf("%s\n", ev.Error())
			case *slack.AckMessage:
				log.Tracef("ACKMessage: %v\n", ev)
			default:
				log.Debugf("Event (%v) received.\n", msg.Type)
			}
		}
	}

}
