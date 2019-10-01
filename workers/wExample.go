package workers

import (
	"fmt"
	"github.com/juju/loggo"
	"github.com/nlopes/slack"
	"reflect"
)

type example struct {
	ClientRTM
	Channels
}

var logExample = loggo.GetLogger("bot.worker.example")

func init() {
	Types["example"] = reflect.TypeOf(example{})
	logExample.SetLogLevel(loggo.DEBUG)
}

func (w *example) Setter( p []interface{} ) {
	fmt.Println(p)
	fmt.Println(w)
}

func (w *example) Start() {
	w.Inbox = make(chan slack.RTMEvent, 10)
	go func() {
		for {
			select {
			case msg := <- w.Inbox :
				ev := msg.Data.(*slack.MessageEvent)
				logExample.Tracef("Message: %v\n", ev.Msg)
				logExample.Debugf("Message (%v): %v\n%v\n", ev.Msg.User, ev.Msg.Text, w)
				w.Rtm.SendMessage(w.Rtm.NewOutgoingMessage( reflect.TypeOf(w).String() + " " + ev.Msg.Text, ev.Channel))
			case <- w.Quit :
				return
			}
		}
	}()
}