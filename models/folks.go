package models

type SupportAgent struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Folks struct {
	Agents []SupportAgent `json:"data"`
}
