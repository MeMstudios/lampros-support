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
