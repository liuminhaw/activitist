package activities

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/liuminhaw/activitist/gpt"
)

type Activity struct {
	Action    string `json:"action"`
	Name      string `json:"name"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Location  string `json:"location"`
}

func Prompt(message string, key string) (string, error) {
	replyMessage, err := gpt.AnalyzeMessage(message, key)
	if err != nil {
		return "", fmt.Errorf("analyze message: %w", err)
	}

	var activity Activity
	err = json.Unmarshal([]byte(replyMessage), &activity)
	if err != nil {
		return "", fmt.Errorf("unmarshal activity: %w", err)
	}

	switch activity.Action {
	case "create":
		replyMessage = activity.create()
	case "update":
		replyMessage = activity.update()
	case "list":
		replyMessage = activity.list()
	case "delete":
		replyMessage = activity.delete()
	default:
		replyMessage = activity.undefined()
	}

	return replyMessage, nil
}

func (a Activity) create() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("建立活動: %s\n", a.Name))
	b.WriteString(fmt.Sprintf("時間: %s\n", a.StartTime))
	if a.EndTime != "" {
		b.WriteString(fmt.Sprintf("結束時間: %s\n", a.EndTime))
	}
	b.WriteString(fmt.Sprintf("地點: %s", a.Location))

	return b.String()
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
