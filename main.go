package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/liuminhaw/activitist/messages"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	gptApiKey := viper.GetString("gptApi.key")
	// log.Debugf("Line channel secret: %s", channelSecret)
	// log.Debugf("Line channel token: %s", channelToken)

	activityC := messages.Message{
		Line: messages.LineAuth{
			ChannelSecret: channelSecret,
			ChannelToken:  channelToken,
		},
		GptApi: messages.GptApiAuth{
			Key: gptApiKey,
		},
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to activitist!!!"))
	})
	r.Post("/line/message", activityC.Receive)
	http.ListenAndServe(":3000", r)
}
