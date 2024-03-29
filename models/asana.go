package models

type Resource struct {
	Id           int    `json:"id"`
	Gid          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
	Error        string `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Help    string `json:"help"`
}

type Task struct {
	Id             int        `json:"id"`
	Gid            string     `json:"gid"`
	Assignee       string     `json:"assignee"`
	AssigneeStatus string     `json:"assignee_status"`
	Completed      bool       `json:"completed"`
	CompletedAt    string     `json:"completed_at"`
	Created        string     `json:"created_at"` //2019-01-31T16:22:27.471
	DueTime        string     `json:"due_at"`     //2019-01-31T16:22:27.471Z
	DueDate        string     `json:"due_on"`     //2019-02-01
	Followers      []Resource `json:"followers"`
	Modified       string     `json:"modified_at"`
	Name           string     `json:"name"`
	Notes          string     `json:"notes"`
	Projects       []Resource `json:"projects"`
	ResourceType   string     `json:"resource_type"`
	StartDate      string     `json:"start_on"`
	Tags           []Resource `json:"tags"`
	Workspace      Resource   `json:"workspace"`
}

type User struct {
	Id           int        `json:"id"`
	Gid          string     `json:"gid"`
	ResourceType string     `json:"resource_type"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Workspaces   []Resource `json:"workspaces"`
}

type ProjectFollowers struct {
	Id        int        `json:"id"`
	Gid       string     `json:"gid"`
	Followers []Resource `json:"followers"`
}

type Webhook struct {
	Id             int      `json:"id"`
	ResourceType   string   `json:"resource_type"`
	Resource       Resource `json:"resource"`
	Target         string   `json:"target"`
	Active         bool     `json:"active"`
	Created        string   `json:"created_at"`
	LastSuccess    string   `json:"last_success_at"`
	LastFailure    string   `json:"last_failure_at"`
	FailureContent string   `json:"last_failure_content"`
}

type Event struct {
	Action   string   `json:"action"`
	Created  string   `json:"created_at"`
	Parent   Resource `json:"parent"`
	Resource Resource `json:"resource"`
	User     Resource `json:"user"`
}

type Story struct {
	Id           int      `json:"id"`
	Gid          string   `json:"gid"`
	Created      string   `json:"created_at"`
	CreatedBy    Resource `json:"created_by"`
	StoryType    string   `json:"resource_subtype"`
	ResourceType string   `json:"resource_type"`
	Source       string   `json:"source"`
	Target       Resource `json:"target"`
	Text         string   `json:"text"`
	Type         string   `json:"type"`
}
