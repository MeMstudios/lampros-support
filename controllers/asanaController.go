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
func getTasks(supportProjectID string) ([]Task, error) {
	projectRespData := getAsanaResponse(parseUrl(AsanaBase + "/projects/" + supportProjectID + "/tasks"))

	var projectResp Response
	//unmarshal the data to the response object
	json.Unmarshal(projectRespData, &projectResp)
	if len(projectResp.Errors) > 0 {
		return nil, errors.New(fmtApiErrors(projectResp.Errors))
	}
	var tasks []Task
	for i := 0; i < len(projectResp.Resources); i++ {
		//Get the task response data
		task, err := getTask(projectResp.Resources[i].Gid)
		if err != nil {
			return tasks, err
		}
		//append the task responses into the array
		tasks = append(tasks, task)
		fmt.Println("task: " + task.Gid)
	}

	return tasks, nil
}

func getTask(taskId string) (Task, error) {
	var resp TaskResponse
	//Get the task response data
	responseData := getAsanaResponse(parseUrl(AsanaBase + "/tasks/" + taskId))
	//Make it an object
	json.Unmarshal(responseData, &resp)
	if len(resp.Errors) > 0 {
		return resp.Task, errors.New(fmtApiErrors(resp.Errors))
	}
	return resp.Task, nil
}

func getStory(storyId string) (Story, error) {
	var resp StoryResponse
	responseData := getAsanaResponse(parseUrl(AsanaBase + "/stories/" + storyId))
	json.Unmarshal(responseData, &resp)
	if len(resp.Errors) > 0 {
		return resp.Story, errors.New(fmtApiErrors(resp.Errors))
	}
	return resp.Story, nil
}

//Return the user object by the userId string
func getUser(userId string) (User, error) {
	var resp UserResponse
	responseData := getAsanaResponse(parseUrl(AsanaBase + "/users/" + userId))
	json.Unmarshal(responseData, &resp)
	if len(resp.Errors) > 0 {
		return resp.User, errors.New(fmtApiErrors(resp.Errors))
	}
	return resp.User, nil
}

//Return the user object and error if we can't find the user
func getUserByEmail(userEmail string) (User, error) {
	var userResp UserResponse
	userRespData := getAsanaResponse(parseUrl(AsanaBase + "/users/" + userEmail))
	json.Unmarshal(userRespData, &userResp)
	if len(userResp.Errors) > 0 {
		return userResp.User, errors.New(fmtApiErrors(userResp.Errors))
	}
	return userResp.User, nil
}

func updateTaskTags(task Task) error {
	params := make(map[string]string)
	params["tag"] = UrgentTagGid

	//Check if the task description(email body) or name(email subject) contains urgent (case-insensitive )
	if caseInsensitiveContains(task.Name, "urgent") ||
		caseInsensitiveContains(task.Notes, "urgent") ||
		caseInsensitiveContains(task.Name, "asap") ||
		caseInsensitiveContains(task.Notes, "asap") ||
		caseInsensitiveContains(task.Name, "important") ||
		caseInsensitiveContains(task.Notes, "important") {
		//Add the urgent tag
		respData := postAsanaRequest(params, parseUrl(AsanaBase+"/tasks/"+task.Gid+"/addTag"))
		var resp Response
		json.Unmarshal(respData, &resp)
		if len(resp.Errors) > 0 {
			return errors.New(fmtApiErrors(resp.Errors))
		}
	}
	return nil
}

func taskIsUrgent(taskGid string) (bool, error) {
	task, err := getTask(taskGid)
	if err != nil {
		return false, err
	}
	for _, t := range task.Tags {
		if t.Gid == UrgentTagGid {
			return false, nil
		}
	}
	return false, nil
}

//Try to add a follower to a task
//Accepts the follower id and task id strings
func updateTaskFollowers(follower, taskId string) (Response, error) {
	params := make(map[string]string)
	params["followers[0]"] = follower
	respData := postAsanaRequest(params, parseUrl(AsanaBase+"/tasks/"+taskId+"/addFollowers"))
	var resp Response
	json.Unmarshal(respData, &resp)
	if len(resp.Errors) > 0 {
		return resp, errors.New(fmtApiErrors(resp.Errors))
	}
	return resp, nil
}

//returns true if a project contains a user email
func checkProjectEmail(userEmail string, supportProjectID string) (bool, error) {
	//get all the followers on a project
	projectResponseData := getAsanaResponse(parseUrl(AsanaBase + "/projects/" + supportProjectID))
	var resp ProjectFollowersResponse
	json.Unmarshal(projectResponseData, &resp)
	if len(resp.Errors) > 0 {
		return false, errors.New(fmtApiErrors(resp.Errors))
	}
	for _, f := range resp.ProjectFollowers.Followers {
		//If there's an error with this then they are a full blown user
		user, err := getUserByEmail(f.Name)
		if err != nil { //You should be able to get the user by the follower's id
			user, err = getUser(f.Gid)
			if err != nil {
				return false, err
			}
			if userEmail == user.Email {
				return true, nil
			}
		} else {
			//If the follower is not a full user in Asana the name will match the email.
			if userEmail == f.Name {
				return true, nil
			}
		}
	}
	return false, nil
}

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

//Returns true if the domain contains asana anywhere (ex: mail.asana.com)
func isAsanaDomain(s string) bool {
	substrings := strings.Split(s, ".")
	for _, s := range substrings {
		if s == "asana" {
			return true
		}
	}
	return false
}

//Parse a url and return string
func parseUrl(u string) string {
	var parsedUrl *url.URL
	parsedUrl, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	return parsedUrl.String()
}

//Format the api erros into a single string
func fmtApiErrors(errs []Error) string {
	var errors []string
	for _, e := range errs {
		err := fmt.Sprint("Error from API: " + e.Message + "\n" + "Get help: " + e.Help)
		errors = append(errors, err)
	}
	return strings.Join(errors, "\n")
}
