package controllers

import (
	"encoding/json"
	"fmt"
	. "lampros-support/models"
	"log"
	"net/url"
	"strings"
)

//Get task details from an asana project
func GetTasks() []Task {
	var Url *url.URL
	Url, err := url.Parse("https://app.asana.com/api/1.0/projects/" + SupportProjectID + "/tasks")
	if err != nil {
		log.Fatal(err)
	}

	projectResponseData := getResponse(Url.String())

	var projectResponseObject Response
	//unmarshal the data to the response object
	json.Unmarshal(projectResponseData, &projectResponseObject)

	var tasks []Task
	for i := 0; i < len(projectResponseObject.Resources); i++ {
		//Build the task URL
		Url, err := url.Parse("https://app.asana.com/api/1.0/tasks/" + projectResponseObject.Resources[i].Gid)
		if err != nil {
			log.Fatal(err)
		}
		var resp TaskResponse
		//Get the task response data
		taskResponseData := getResponse(Url.String())
		//Make it an object
		json.Unmarshal(taskResponseData, &resp)
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
			fmt.Println(tasks[i].Gid + " Contains Urgent")
			Url, err := url.Parse("https://app.asana.com/api/1.0/tasks/" + tasks[i].Gid + "/addTag")
			if err != nil {
				log.Fatal(err)
			}
			respData := postRequest(params, Url.String())
			var resp Response
			json.Unmarshal(respData, &resp)
			if len(resp.Resources) > 0 {
				log.Fatal(resp.Resources[0].Name)
			}
		}
		i++
	}
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
