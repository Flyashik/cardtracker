package models

type Notification struct {
	Id        int    `json:"notification_id"`
	ModelTag  string `json:"model_tag"`
	Source    string `json:"notification_source"`
	Sender    string `json:"sender"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
}
