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
	projectRespData := getResponse(parseUrl(AsanaBase + "/projects/" + SupportProjectID + "/tasks"))

	var projectResp Response
	//unmarshal the data to the response object
	json.Unmarshal(projectRespData, &projectResp)
	if len(projectResp.Errors) > 0 {
		logApiErrors(projectResp.Errors)
	}
	var tasks []Task
	for i := 0; i < len(projectResp.Resources); i++ {
		var resp TaskResponse
		//Get the task response data
		taskResponseData := getResponse(parseUrl(AsanaBase + "/tasks/" + projectResp.Resources[i].Gid))
		//Make it an object
		json.Unmarshal(taskResponseData, &resp)
		if len(resp.Errors) > 0 {
			logApiErrors(resp.Errors)
		}
		//append the task responses into the array
		tasks = append(tasks, resp.Task)
		fmt.Println("task: " + tasks[i].Gid)
	}

	return tasks
}

func UpdateTasks(tasks []Task) {
	i := 0
	var params map[string]string
	params = make(map[string]string)
	params["tag"] = UrgentTagGid
	for i < len(tasks) {
		//Check if the task description(email body) or name(email subject) contains urgent (case-insensitive )
		if CaseInsensitiveContains(tasks[i].Notes, "urgent") || CaseInsensitiveContains(tasks[i].Name, "urgent") {
			//Add the urgent tag
			respData := postRequest(params, parseUrl(AsanaBase+"/tasks/"+tasks[i].Gid+"/addTag"))
			var resp Response
			json.Unmarshal(respData, &resp)
			if len(resp.Errors) > 0 {
				logApiErrors(resp.Errors)
			}
		}
		i++
	}
}

func GetUserByEmail(userEmail string) (User, error) {
	var userResp UserResponse
	userRespData := getResponse(parseUrl(AsanaBase + "/users/" + userEmail))
	json.Unmarshal(userRespData, &userResp)
	if len(userResp.Errors) > 0 {
		logApiErrors(userResp.Errors)
		return userResp.User, errors.New("can't find user!")
	}
	return userResp.User, nil
}

func CheckProjectEmail(userEmail string) bool {
	//get the user by email
	u, err := GetUserByEmail(userEmail)
	if err != nil {
		log.Println(err)
		return false
	}
	//get all the followers on a project
	projectResponseData := getResponse(parseUrl(AsanaBase + "/projects/" + SupportProjectID + "?opt_fields=followers"))
	var resp ProjectFollowersResponse
	json.Unmarshal(projectResponseData, &resp)
	if len(resp.Errors) > 0 {
		logApiErrors(resp.Errors)
	}
	for _, f := range resp.ProjectFollowers.Followers {
		if u.Gid == f.Gid {
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
