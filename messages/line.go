package messages

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/liuminhaw/activitist/activities"
	"github.com/liuminhaw/activitist/gpt"
	log "github.com/sirupsen/logrus"
)

type Line struct {
	LineService     *LineService
	Gpt             gpt.GptAuth
	ActivityService *activities.ActivityService
	RegisterService *activities.RegisterService
}

type LineService struct {
	ChannelSecret string
	ChannelToken  string
	DB            *sql.DB
}

type User struct {
	ID     int
	LineID string
}

func (l Line) Receive(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(l.LineService.ChannelSecret, l.LineService.ChannelToken)
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

	var user User

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.WithFields(log.Fields{
					"method": "POST",
					"path":   "/line/message",
					"event":  "textMessage",
				}).Info(message.Text)

				user.LineID = event.Source.UserID

				var replyMessage string
				switch {
				case message.Text == ":whoami":
					eventSource := event.Source
					log.Info(fmt.Sprintf("%+v", *eventSource))
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("User id: %s", user.LineID))).Do()
					return
				case strings.HasPrefix(message.Text, ":register"):
					token := strings.TrimSpace(strings.ReplaceAll(message.Text, ":register", ""))
					if token == "" {
						register, err := l.RegisterService.TokenCreate(user.LineID)
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
						err := l.RegisterService.Register(event.Source.UserID, token)
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
				case strings.HasPrefix(message.Text, ":act"):
					prompt := strings.ReplaceAll(message.Text, ":act", "")
					// TODO: Check if user is registered
					id, err := l.checkIdentity(user.LineID, "user")
					if err != nil {
						replyMessage := "使用者尚未註冊"
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
						return
					}
					user.ID = id
					err = l.ActivityService.Prompt(prompt, l.Gpt.ApiKey)
					if err != nil {
						log.WithFields(log.Fields{
							"method": "POST",
							"path":   "/line/message",
							"event":  "textMessage",
						}).Error(err)
						return
					}

					replyMessage, err := l.ActivityService.Create(user.ID)
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

func (l Line) checkIdentity(lineID, source string) (int, error) {
	var id int
	switch source {
	case "group":
		return 0, errors.New("group check not implemented")
	case "user":
		row := l.LineService.DB.QueryRow(`
				SELECT id FROM users WHERE line_id = $1
			`, lineID)
		err := row.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("line check identity: %w", err)
		}
	}

	return id, nil
}
