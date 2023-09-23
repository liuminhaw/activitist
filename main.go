package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/liuminhaw/activitist/activities"
	"github.com/liuminhaw/activitist/gpt"
	"github.com/liuminhaw/activitist/messages"
	"github.com/liuminhaw/activitist/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	PSQL    models.PostgresConfig
	Linebot messages.LineService
	Gpt     gpt.GptAuth
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func loadConfig() (config, error) {
	var cfg config

	viper.SetConfigName("activitist")         // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/activitist.d/") // path to look for the config file in
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	err := viper.ReadInConfig()               // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		return cfg, fmt.Errorf("read config: %w", err)
	}

	cfg.PSQL.Host = viper.GetString("db.host")
	cfg.PSQL.Port = viper.GetString("db.port")
	cfg.PSQL.User = viper.GetString("db.user")
	cfg.PSQL.Password = viper.GetString("db.password")
	cfg.PSQL.Database = viper.GetString("db.database")
	cfg.PSQL.SSLMode = viper.GetString("db.sslMode")

	cfg.Linebot.ChannelSecret = viper.GetString("line.channelSecret")
	cfg.Linebot.ChannelToken = viper.GetString("line.channelToken")
	cfg.Gpt.ApiKey = viper.GetString("gptApi.key")

	return cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup service
	lineS := messages.Line{
		LineService: &messages.LineService{
			DB:            db,
			ChannelSecret: cfg.Linebot.ChannelSecret,
			ChannelToken:  cfg.Linebot.ChannelToken,
		},
		Gpt: cfg.Gpt,
		ActivityService: &activities.ActivityService{
			DB: db,
		},
		RegisterService: &activities.RegisterService{
			DB:            db,
			BytesPerToken: 32,
			Duration:      1 * time.Hour,
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
	r.Post("/line/message", lineS.Receive)
	http.ListenAndServe(":3000", r)
}
