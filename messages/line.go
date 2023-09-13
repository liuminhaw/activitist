package messages

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/liuminhaw/activitist/activities"
	"github.com/liuminhaw/activitist/gpt"
	log "github.com/sirupsen/logrus"
)

type LineService struct {
	Line            LineAuth
	Gpt             gpt.GptAuth
	ActivityService *activities.ActivityService
	RegisterService *activities.RegisterService
}

type LineAuth struct {
	ChannelSecret string
	ChannelToken  string
}

func (ls LineService) Receive(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(ls.Line.ChannelSecret, ls.Line.ChannelToken)
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

				var replyMessage string
				switch {
				case message.Text == "@whoami":
					eventSource := event.Source
					log.Info(fmt.Sprintf("%+v", *eventSource))
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("User id: %s", eventSource.UserID))).Do()
					return
				case strings.HasPrefix(message.Text, "@register"):
					token := strings.TrimSpace(strings.ReplaceAll(message.Text, "@register", ""))
					if token == "" {
						register, err := ls.RegisterService.TokenCreate(event.Source.UserID)
						if err != nil {
							log.WithFields(log.Fields{
								"method": "POST",
								"path":   "/line/message",
								"event":  "textMessage",
							}).Error(err)
							return
						}
						log.WithFields(log.Fields{
							"user":  register.UserID,
							"token": register.Token,
						}).Info("Token created")
						replyMessage = "Please obtain and send generated token"
					} else {
						err := ls.RegisterService.Register(event.Source.UserID, token)
						if err != nil {
							log.WithFields(log.Fields{
								"method": "POST",
								"path":   "/line/message",
								"event":  "textMessage",
							}).Error(err)
							replyMessage = "Failed to register"
						} else {
							replyMessage = "Register successful"
						}
					}
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
					return
				case strings.HasPrefix(message.Text, "@act"):
					prompt := strings.ReplaceAll(message.Text, "@act", "")
					replyMessage, err := ls.ActivityService.Prompt(prompt, ls.Gpt.ApiKey)
					if err != nil {
						log.WithFields(log.Fields{
							"method": "POST",
							"path":   "/line/message",
							"event":  "textMessage",
						}).Error(err)
						return
					}
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
					return
				}
			}
		}
	}
}
