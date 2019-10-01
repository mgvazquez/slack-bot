package workers

import (
	"github.com/juju/loggo"
	"github.com/nlopes/slack"
	"math/rand"
	"reflect"
	"regexp"
	"time"
)

type troll struct {
	ClientRTM
	Channels
	Advices []t_advice
}

type t_advice struct {
	regexp *regexp.Regexp
	rtas   []t_rta
}

type t_rta struct {
	msg string
	img string
}

var logTroll = loggo.GetLogger("bot.worker.troll")

func init() {
	Types["troll"] = reflect.TypeOf(troll{})
	logTroll.SetLogLevel(loggo.DEBUG)
}

func (w *troll) Setter( p []interface{} ) {
	for _, a := range p {
		adv := a.(map[string]interface {})
		logTroll.Tracef("Advise to process: %v\n", adv)
		// Se agrega el "(?i)" para hacer case insensitive la RegEx por default
		re, err := regexp.Compile("(?i)"+adv["regex"].(string))
		if err != nil {
			logTroll.Errorf("Fail to build regex: %v\n", adv["regex"].(string))
			continue
		} else {
			logTroll.Debugf("RegEx: \"%v\"\n", adv["regex"].(string))
			advice := t_advice{}
			advice.regexp = re
			for _, r := range adv["rtas"].([]interface{}) {
				_r := r.(map[string]interface{})
				rta := t_rta{"",""}
				if _r["msg"] != nil {
					logTroll.Debugf(" - Reply (txt): %v\n", _r["msg"].(string))
					rta.msg = _r["msg"].(string)
				}
				if _r["img"] != nil {
					logTroll.Debugf(" - Reply (img): %v\n", _r["img"].(string))
					rta.img = _r["img"].(string)
				}
				advice.rtas = append(advice.rtas, rta)
			}
			w.Advices = append(w.Advices, advice)
		}
	}
}

func (w *troll) Start() {
	w.Inbox = make(chan slack.RTMEvent, 10)
	rand.Seed(time.Now().UnixNano())
	// Busca el string '<user>' para despues poder reemplazarlo por el usuario
	// que envio el mensaje
	userRegEX := regexp.MustCompile("<user>")
	msgParams := slack.NewPostMessageParameters()
	msgParams.EscapeText = false
	msgParams.AsUser = true
	msgParams.Markdown = true
	go func() {
		for {
			select {
			case msg := <- w.Inbox :
				ev := msg.Data.(*slack.MessageEvent)
				reponses := []t_rta{}
				logTroll.Tracef("Incomming message: %v\n", ev.Msg)
				logTroll.Debugf("%v: %v\n", ev.Msg.User, ev.Msg.Text)
				// Arma un array de las respuestas validas para todos los matches
				for _, a := range w.Advices {
					t := a.regexp.MatchString(ev.Msg.Text)
					if t {
						logTroll.Debugf("Match message with regex: %v\n", a.regexp.String())
						for _, r := range a.rtas {
							reponses = append(reponses, r)
						}
					}
				}
				if len(reponses) >= 1 {
					attachment := slack.Attachment{}
					message := ""
					m := t_rta{}
					m = reponses[rand.Intn(len(reponses))]
					m.msg = userRegEX.ReplaceAllString(m.msg, "<@"+ev.Msg.User+">")
					logTroll.Tracef("Responses list: %v\n", reponses)
					logTroll.Tracef("Select response: %v\n", m)
					if m.img != "" {
						attachment.ImageURL = m.img
						attachment.Text = m.msg
					} else {
						message = m.msg
					}
					chanResponse, _, err := w.Rtm.PostMessage(
						ev.Channel,
						slack.MsgOptionPostMessageParameters(msgParams),
						slack.MsgOptionAttachments(attachment),
						slack.MsgOptionText(message,false))
					if err != nil {
						logTroll.Errorf("Failed to send message: %v", err)
						logTroll.Debugf("Channel response: %v", chanResponse)
					} else {
						logTroll.Infof("Message sended: [TEXT] %v / [ATTACH] %v\n", message,attachment)
					}
				}
			case <- w.Quit :
				return
			}
		}
	}()
}