package workers

import (
	"github.com/juju/loggo"
	"github.com/nlopes/slack"
	"reflect"
)

type Config struct {
	Type string `json:"type"`
	Params []interface{} `json:"params"`
}

type ClientRTM struct {
	Rtm *slack.RTM
}

type Channels struct {
	Inbox  chan slack.RTMEvent
	Outbox chan slack.RTMEvent
	Quit   chan struct{}
}

type Worker interface {
	Setter( p []interface{} )
	Start()
}

var Types = make(map[string]reflect.Type)
var log = loggo.GetLogger("bot.worker")

func init() {
	log.SetLogLevel(loggo.DEBUG)
}

func Constructor( name string, rtm *slack.RTM) Worker {
	worker := reflect.New(Types[name]).Interface().(Worker)
	reflect.ValueOf(worker).Elem().FieldByName("Rtm" ).Set(reflect.ValueOf(rtm))
	log.Debugf("Worker %v created", reflect.TypeOf(worker).String())
	return worker
}