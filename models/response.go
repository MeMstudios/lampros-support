package models

type Response struct {
	Resources []Resource `json:"data"`
	Errors    []Error    `json:"errors"`
}

type TaskResponse struct {
	Task   Task    `json:"data"`
	Errors []Error `json:"errors"`
}

type UserResponse struct {
	User   User    `json:"data"`
	Errors []Error `json:"errors"`
}

type ProjectFollowersResponse struct {
	ProjectFollowers ProjectFollowers `json:"data"`
	Errors           []Error          `json:"errors"`
}

type StoryResponse struct {
	Story  Story   `json:"data"`
	Errors []Error `json:"errors"`
}

type WebhookResponse struct {
	Webhook Webhook `json:"data"`
	Errors  []Error `json:"errors"`
}

type WebhookEvent struct {
	Events []Event `json:"events"`
	Errors []Error `json:"errors"`
}
