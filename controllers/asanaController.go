package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	. "lampros-support/models"
	"log"
	"net/url"
	"strings"
)

//Get task details from an asana project
func GetTasks() []Task {
	projectRespData := getAsanaResponse(parseUrl(AsanaBase + "/projects/" + SupportProjectID + "/tasks"))

	var projectResp Response
	//unmarshal the data to the response object
	json.Unmarshal(projectRespData, &projectResp)
	if len(projectResp.Errors) > 0 {
		logApiErrors(projectResp.Errors)
	}
	var tasks []Task
	for i := 0; i < len(projectResp.Resources); i++ {
		//Get the task response data
		task := GetTask(projectResp.Resources[i].Gid)
		//append the task responses into the array
		tasks = append(tasks, task)
		fmt.Println("task: " + task.Gid)
	}

	return tasks
}

func GetTask(taskId string) Task {
	var resp TaskResponse
	//Get the task response data
	responseData := getAsanaResponse(parseUrl(AsanaBase + "/tasks/" + taskId))
	//Make it an object
	json.Unmarshal(responseData, &resp)
	if len(resp.Errors) > 0 {
		logApiErrors(resp.Errors)
	}
	return resp.Task
}

func GetStory(storyId string) Story {
	var resp StoryResponse
	responseData := getAsanaResponse(parseUrl(AsanaBase + "/stories/" + storyId))
	json.Unmarshal(responseData, &resp)
	if len(resp.Errors) > 0 {
		logApiErrors(resp.Errors)
	}
	return resp.Story

}

func GetUser(userId string) User {
	var resp UserResponse
	responseData := getAsanaResponse(parseUrl(AsanaBase + "/users/" + userId))
	json.Unmarshal(responseData, &resp)
	if len(resp.Errors) > 0 {
		logApiErrors(resp.Errors)
	}
	return resp.User
}

func UpdateTaskTags(task Task) {
	var params map[string]string
	params = make(map[string]string)
	params["tag"] = UrgentTagGid

	//Check if the task description(email body) or name(email subject) contains urgent (case-insensitive )
	if CaseInsensitiveContains(task.Notes, "urgent") || CaseInsensitiveContains(task.Name, "urgent") {
		//Add the urgent tag
		respData := postAsanaRequest(params, parseUrl(AsanaBase+"/tasks/"+task.Gid+"/addTag"))
		var resp Response
		json.Unmarshal(respData, &resp)
		if len(resp.Errors) > 0 {
			logApiErrors(resp.Errors)
		}
	}
}

func TaskIsUrgent(taskGid string) bool {
	task := GetTask(taskGid)
	for _, t := range task.Tags {
		if t.Gid == UrgentTagGid {
			return true
		}
	}
	return false
}

func GetUserByEmail(userEmail string) (User, error) {
	var userResp UserResponse
	userRespData := getAsanaResponse(parseUrl(AsanaBase + "/users/" + userEmail))
	json.Unmarshal(userRespData, &userResp)
	if len(userResp.Errors) > 0 {
		logApiErrors(userResp.Errors)
		return userResp.User, errors.New("can't find user!")
	}
	return userResp.User, nil
}

func UpdateTaskFollowers(follower, taskId string) Response {
	params := make(map[string]string)
	params["followers[0]"] = follower
	respData := postAsanaRequest(params, parseUrl(AsanaBase+"/tasks/"+taskId+"/addFollowers"))
	var resp Response
	json.Unmarshal(respData, &resp)
	if len(resp.Errors) > 0 {
		logApiErrors(resp.Errors)
	}
	return resp
}

func CheckProjectEmail(userEmail string) bool {
	//get all the followers on a project
	projectResponseData := getAsanaResponse(parseUrl(AsanaBase + "/projects/" + SupportProjectID))
	var resp ProjectFollowersResponse
	json.Unmarshal(projectResponseData, &resp)
	if len(resp.Errors) > 0 {
		logApiErrors(resp.Errors)
	}
	for _, f := range resp.ProjectFollowers.Followers {
		if userEmail == f.Name {
			return true
		}
	}
	return false
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func parseUrl(u string) string {
	var Url *url.URL
	Url, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	return Url.String()
}

func logApiErrors(errs []Error) {
	for _, e := range errs {
		fmt.Println("Error from API: " + e.Message + "\n" + "Get help: " + e.Help)
	}
}
