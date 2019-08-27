package models

type SupportAgent struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type SupportProject struct {
	SupportEmail string         `json:"email"`
	ProjectId    string         `json:"id"`
	Agents       []SupportAgent `json:"agents"`
}

type Projects struct {
	Projects []SupportProject `json:"projects"`
}
