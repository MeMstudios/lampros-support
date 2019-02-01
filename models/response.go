package models

type Response struct {
	Resources []Resource `json:"data"`
}

type TaskResponse struct {
	Task Task `json:"data"`
}
