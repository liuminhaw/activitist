package messages

import (
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/liuminhaw/activitist/activities"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Line   LineAuth
	GptApi GptApiAuth
}

type LineAuth struct {
	ChannelSecret string
	ChannelToken  string
}

type GptApiAuth struct {
	Key string
}

func (m Message) Receive(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(m.Line.ChannelSecret, m.Line.ChannelToken)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "POST",
			"path":   "/line/message",
		}).Errorf("new linebot: %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	events, err := bot.ParseRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "POST",
			"path":   "/line/message",
		}).Errorf("parse request: %s", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.WithFields(log.Fields{
					"method": "POST",
					"path":   "/line/message",
					"event":  "textMessage",
				}).Info(message.Text)
				if ok := strings.HasPrefix(message.Text, "@act"); ok {
					prompt := strings.ReplaceAll(message.Text, "@act", "")
					replyMessage, err := activities.Prompt(prompt, m.GptApi.Key)
					if err != nil {
						log.WithFields(log.Fields{
							"method": "POST",
							"path":   "/line/message",
							"event":  "textMessage",
						}).Error(err)
						return
					}
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				}
			}
		}
	}
}
