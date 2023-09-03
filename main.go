package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/liuminhaw/activitist/gpt"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	viper.SetConfigName("activitist")         // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/activitist.d/") // path to look for the config file in
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	err := viper.ReadInConfig()               // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		log.WithFields(log.Fields{"event": "readConfig"}).Fatal(err)
	}

	channelSecret := viper.GetString("line.channelSecret")
	channelToken := viper.GetString("line.channelToken")
	// log.Debugf("Line channel secret: %s", channelSecret)
	// log.Debugf("Line channel token: %s", channelToken)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to activitist!!!"))
	})
	r.Post("/line/message", func(w http.ResponseWriter, r *http.Request) {
		bot, err := linebot.New(channelSecret, channelToken)
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
						gptService := gpt.NewAssistantService(&gpt.AssistantConfig{
							Key: viper.GetString("gptApi.key"),
						}, prompt)
						replyMessage, err := gptService.AnalyzeMessage()
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
	})
	http.ListenAndServe(":3000", r)
}
