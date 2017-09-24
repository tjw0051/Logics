package models

import "github.com/jinzhu/gorm"

type ConfigModel struct {
	MasterAPIKey string `json:"masterApiKey"`
}

type MessageModel struct {
	gorm.Model
	Platform     string `json:"platform"`
	Project      string `json:"project"`
	Component    string `json:"component"`
	LogType      string `json:"logType"` // crash, error, warning, info
	Severity     string `json:"severity"`
	Message      string `json:"message"`
	MessageGroup string `json:"messageGroup"`
	Data         string `json:"data"`
	UserID       string `json:"userId"`
}

type MessagesModel []MessageModel

type QueryModel struct {
	Messages MessagesModel `json:"messages"`
	Page     int           `json:"page"`
}

type KeyModel struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type KeysModel []KeyModel
