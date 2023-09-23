package activities

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/liuminhaw/activitist/gpt"
)

type ActivityService struct {
	Activity *Activity
	DB       *sql.DB
}

type Activity struct {
	Action    string `json:"action"`
	Name      string `json:"name"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Location  string `json:"location"`
}

func (service *ActivityService) Prompt(message string, key string) error {
	replyMessage, err := gpt.AnalyzeMessage(message, key)
	if err != nil {
		return fmt.Errorf("analyze message: %w", err)
	}
	fmt.Printf("Prompt reply message: %s\n", replyMessage)

	var activity Activity
	err = json.Unmarshal([]byte(replyMessage), &activity)
	if err != nil {
		return fmt.Errorf("unmarshal activity: %w", err)
	}

	service.Activity = &activity

	return nil
}

func (service *ActivityService) Create(id int) (string, error) {
	var b bytes.Buffer
	var retID int

	// Insert into database
	row := service.DB.QueryRow(`
		INSERT INTO individual_activities (activity, location, user_id) 
		VALUES ($1, $2, $3) RETURNING id
	`, service.Activity.Name, service.Activity.Location, id)
	err := row.Scan(&retID)
	if err != nil {
		return "", fmt.Errorf("create activity: %w", err)
	}

	b.WriteString(fmt.Sprintf("建立活動: %s\n", service.Activity.Name))
	b.WriteString(fmt.Sprintf("時間: %s\n", service.Activity.StartTime))
	if service.Activity.EndTime != "" {
		b.WriteString(fmt.Sprintf("結束時間: %s\n", service.Activity.EndTime))
	}
	b.WriteString(fmt.Sprintf("地點: %s", service.Activity.Location))

	return b.String(), nil
}

func (a Activity) update() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("更新活動: %s\n", a.Name))
	b.WriteString(fmt.Sprintf("時間: %s\n", a.StartTime))
	if a.EndTime != "" {
		b.WriteString(fmt.Sprintf("結束時間: %s\n", a.EndTime))
	}
	b.WriteString(fmt.Sprintf("地點: %s", a.Location))

	return b.String()
}

func (a Activity) delete() string {
	return fmt.Sprintf("刪除活動: %s\n", a.Name)
}

func (a Activity) list() string {
	return fmt.Sprintln("列出近期活動")
}

func (a Activity) undefined() string {
	return fmt.Sprintln("無法辨識的動作")
}
