package gpt

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const systemRole = `You are an professional traditional chinese activity assistant, please analyze given input and conclude in points below:
- action: should be one of "create" if aimed to create an activity, "update" if aimed to make change to activity, "delete" if aimed to remove an activity, "list" if aim to view activities, or "undefined" if non of above is matched
- name: name of the activity, should be easy to recognize and non-repetitive
- starttime (in the format of YYYY-mm-dd HH:MM:SS)
- endtime (in the format of YYYY-mm-dd HH:MM:SS)
- location
Items other than action and name can be omitted if there is no information.Based on the provided information, please determine whether this is intended for adding an event, updating an event, deleting an event, listing events, or if it's indeterminate. Fill in the action field with "create", "update", "delete", "list", or "undefined" accordingly.
For the name field, provide a distinctive event name that is not easily confused with other events in traditional Chinese.
For the starttime field, provide the event's starting time.
For the endtime field, provide the event's ending time, leaving it empty if no specific information is available.
For the location field, provide the event's location.
Except for action key, which are mandatory, the other fields should not be included if there is no available information to determine their content.
For example, if the message is "Add location to the Thin Eatery event: Thin Eatery," it should be interpreted as "Modify the location of the Thin Eatery event to Thin Eatery."
Here are possible JSON responses:
Adding an event: {"action": "create","name": "Festival Concert","starttime": "2023-09-15 18:00","endtime": "2023-09-15 22:00","location": "City Center Music Hall"}
Deleting an event: {"action": "delete","name": "Weekend Market"}
Updating an event: {"action": "update","name": "Health Seminar","starttime": "2023-09-20 14:00","endtime": "2023-09-20 16:00","location": "Health Center"}
Listing events: {"action": "list"}
Indeterminate situation: {"action": "undefined"}`

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

type AssistantConfig struct {
	Key string
}

type AssistantService struct {
	Message       string
	assistantConf *AssistantConfig
}

func NewAssistantService(config *AssistantConfig, msg string) *AssistantService {
	as := AssistantService{
		Message:       msg,
		assistantConf: config,
	}

	return &as
}

func (as *AssistantService) AnalyzeMessage() (string, error) {
	client := openai.NewClient(as.assistantConf.Key)
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
					Role:    openai.ChatMessageRoleUser,
					Content: as.Message,
				},
			},
			// Functions:    []openai.FunctionDefinition{schema},
			// FunctionCall: &openai.FunctionCall{Name: "get_action", Arguments: ""},
		},
	)
	if err != nil {
		return "", fmt.Errorf("analyze message: %w", err)
	}
	// return resp.Choices[0].Message.FunctionCall.Arguments, nil
	return resp.Choices[0].Message.Content, nil
}
