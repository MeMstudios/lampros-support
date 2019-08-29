/*
Models package contains all struct models corresponding to each API specification for json responses.
Which means there are some unused data members.
response.go is all for Asana responses which contain the Asana models,
since they are sent with a potential array of errors as well.
support.go implements the models for our projects.json.
*/
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
