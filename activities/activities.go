package activities

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/liuminhaw/activitist/gpt"
)

type ActivityService struct {
	Activity *Activity
	DB       *sql.DB
}

type Activity struct {
	Action    string       `json:"action"`
	Name      string       `json:"name"`
	StartTime activityTime `json:"starttime"`
	EndTime   activityTime `json:"endtime,omitempty"`
	Location  string       `json:"location,omitempty"`
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

func (service *ActivityService) Action(id int) (string, error) {
	var reply string
	var err error
	switch service.Activity.Action {
	case "create":
		reply, err = service.create(id)
	case "list":
		reply, err = service.list(id)
	}
	if err != nil {
		return "", fmt.Errorf("action: %w", err)
	}

	return reply, nil
}

func (service *ActivityService) create(id int) (string, error) {
	var b bytes.Buffer
	var retID int

	// Insert into database
	row := service.DB.QueryRow(`
		INSERT INTO individual_activities (activity, location, user_id, starttime, endtime) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, service.Activity.Name, service.Activity.Location, id,
		service.Activity.StartTime.NewNullTime(), service.Activity.EndTime.NewNullTime(),
	)
	err := row.Scan(&retID)
	if err != nil {
		return "", fmt.Errorf("create activity: %w", err)
	}

	b.WriteString(fmt.Sprintf("建立活動: %s\n", service.Activity.Name))
	b.WriteString(fmt.Sprintf("時間: %s\n", service.Activity.StartTime.String()))
	if !time.Time(service.Activity.EndTime).IsZero() {
		b.WriteString(fmt.Sprintf("結束時間: %s\n", service.Activity.EndTime.String()))
	}
	b.WriteString(fmt.Sprintf("地點: %s", service.Activity.Location))

	return b.String(), nil
}

func (service *ActivityService) list(id int) (string, error) {
	var b bytes.Buffer

	rows, err := service.DB.Query(`
		select activity, starttime 
		FROM individual_activities 
		INNER JOIN users 
		    ON user_id = users.id 
		    WHERE user_id = $1 AND starttime >= $2 
		ORDER BY starttime ASC;
	`, id, time.Now())
	if err != nil {
		return "", fmt.Errorf("list activities: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			activity  string
			startTime time.Time
		)
		if err := rows.Scan(&activity, &startTime); err != nil {
			return "", fmt.Errorf("list activity: %w", err)
		}

		actStartTime := activityTime(startTime)
		b.WriteString(fmt.Sprintf("%s\n", activity))
		b.WriteString(fmt.Sprintf("時間: %s\n", actStartTime.String()))
		b.WriteString("---\n")
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("list activities: %w", err)
	}

	return b.String(), nil
}

func (a Activity) update() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("更新活動: %s\n", a.Name))
	b.WriteString(fmt.Sprintf("時間: %s\n", a.StartTime.String()))
	if !time.Time(a.EndTime).IsZero() {
		b.WriteString(fmt.Sprintf("結束時間: %s\n", a.EndTime.String()))
	}
	b.WriteString(fmt.Sprintf("地點: %s", a.Location))

	return b.String()
}

func (a Activity) delete() string {
	return fmt.Sprintf("刪除活動: %s\n", a.Name)
}

func (a Activity) undefined() string {
	return fmt.Sprintln("無法辨識的動作")
}
