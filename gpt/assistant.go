package gpt

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const systemRole = `Please response in traditional chinese. You are an professional activity assistant, please analyze given input and conclude in points below:
- action: should be "create" if aimed to create an activity, "update" if aimed to make change to activity, "delete" if aimed to remove an activity, "list" if aim to view activities, or "undefined" if non of above is matched
- name: name of the activity, should be easy to recognize and non-repetitive
- starttime (in the format of YYYY-mm-dd HH:MM:SS)
- endtime (in the format of YYYY-mm-dd HH:MM:SS)
- location
"action" key is always needed. Based on the provided information, please determine whether this is intended for adding an event, updating an event, deleting an event, listing events, or if it's indeterminate. Fill in the action field with "create", "update", "delete", "list", or "undefined" accordingly.
Provide a distinctive event name for "name" key which is not easily confused with other events in traditional Chinese.
Provide the event's starting time for "starttime" key.
Provide the event's ending time for "endttime" key only if there is information about activity ending time in the given prompt, or this field  should be ignored.
Provide the related location for "location" key.
Except for "action" key, the other fields should not be included if there is no available information to determine their content.
For example, if the message is "Add location to the Thin Eatery event: Thin Eatery," it should be interpreted as "Modify the location of the Thin Eatery event to Thin Eatery."
Here are possible JSON responses:
Adding an event: {"action": "create","name": "Festival Concert","starttime": "2023-09-15 18:00:00","endtime": "2023-09-15 22:00:00","location": "City Center Music Hall"}
Deleting an event: {"action": "delete","name": "Weekend Market"}
Updating an event: {"action": "update","name": "Health Seminar","starttime": "2023-09-20 14:00:00","endtime": "2023-09-20 16:00:00","location": "Health Center"}
Listing events: {"action": "list"}
Indeterminate situation: {"action": "undefined"}`

const defaultTimeFormat = "2006-01-02 15:04:05"

var schema openai.FunctionDefinition = openai.FunctionDefinition{
	Name: "get_action",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"action": {
				Type:        jsonschema.String,
				Description: "Action aimed to produced on given activity",
				Enum:        []string{"create", "update", "delete", "list", "undefined"},
			},
			"name": {
				Type:        jsonschema.String,
				Description: "Name of the activity, should be easy to recognize and non-repetitive",
			},
			"starttime": {
				Type:        jsonschema.String,
				Description: "Start time of the activity, should be in ISO 8601 format",
			},
			"endtime": {
				Type:        jsonschema.String,
				Description: "End time of the activity, should be in ISO 8601 format",
			},
			"location": {
				Type:        jsonschema.String,
				Description: "Location of the activity",
			},
		},
		Required: []string{"action", "name"},
	},
}

type GptAuth struct {
	ApiKey string
}

// func (as *AssistantService) AnalyzeMessage() (string, error) {
func AnalyzeMessage(message string, key string) (string, error) {
	// client := openai.NewClient(as.assistantConf.Key)
	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemRole,
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: fmt.Sprintf("Current time is %s", currentTime(defaultTimeFormat)),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Functions:    []openai.FunctionDefinition{schema},
			FunctionCall: &openai.FunctionCall{Name: "get_action", Arguments: ""},
		},
	)
	if err != nil {
		return "", fmt.Errorf("analyze message: %w", err)
	}
	return resp.Choices[0].Message.FunctionCall.Arguments, nil
	// return resp.Choices[0].Message.Content, nil
}

func currentTime(format string) string {
	current := time.Now()

	return current.Format(format)
}
