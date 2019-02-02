package models

type Response struct {
	Resources []Resource `json:"data"`
}

type TaskResponse struct {
	Task Task `json:"data"`
}

type UserResponse struct {
	User User `json:"data"`
}

type ProjectFollowersResponse struct {
	ProjectFollowers ProjectFollowers `json:"data"`
}
