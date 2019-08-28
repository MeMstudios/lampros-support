/*
Package containing the controllers for each API implementation (Asana, Twilio, and Gmail).
Twilio controller has the ticker/timer code for timing the texts.
This file is the HTTP controller for our REST API containing most of the business logic.
You will need a credentials.go file (API credentials) in this package described in the README.
*/
package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	. "lampros-support/models"

	"github.com/gorilla/mux"
)

//Make my timers array which should be accessible by this file
var timers []TickerTimer

//GET a JSON response from an API
func getAsanaResponse(request string) []byte {
	bearer := "Bearer " + AsanaAccessToken

	req, err := http.NewRequest("GET", request, nil)

	req.Header.Add("Authorization", bearer)

	//send request using http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData

}

//POST Request with a map of parameters
//returns the response as byte array
func postAsanaRequest(params map[string]string, request string) []byte {
	data := url.Values{}

	for key, value := range params {
		data.Set(key, value)
	}

	bearer := "Bearer " + AsanaAccessToken

	req, err := http.NewRequest("POST", request, strings.NewReader(data.Encode()))

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//send request using http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

//POST Request for twilio
//returns the response as byte array
func postTwilioRequest(params map[string]string, request string) []byte {
	data := url.Values{}

	for key, value := range params {
		data.Set(key, value)
	}

	req, err := http.NewRequest("POST", request, strings.NewReader(data.Encode()))

	req.SetBasicAuth(TwilioSID, TwilioAUTH)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//send request using http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	test := []byte("Test")
	w.Header().Set("X-Hook-Secret", "12301985bugaloo120198")
	w.WriteHeader(http.StatusOK)
	w.Write(test)
	fmt.Println("Test")
}

//Not using right now
// func AddAgentEndpoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var newAgent SupportAgent
// 	response := []byte("")

// 	if err := json.NewDecoder(r.Body).Decode(&newAgent); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		fmt.Printf("invalid request payload: %v\n", err)
// 		return
// 	}
// 	apiKey := r.Header.Get("Api-Key")
// 	if apiKey != OurApiKey {
// 		respondWithError(w, http.StatusBadRequest, "something went wrong")
// 		fmt.Println("invalid API key!")
// 		return
// 	} else {
// 		agents.Agents = append(agents.Agents, newAgent)
// 		err := writeAgentsToJSON("/home/michael/go/src/agents.json")
// 		if err != nil {
// 			log.Fatalf("Error saving agent json: %v", err)
// 		}
// 		response = []byte("Added email and phone to support list.")
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }

//POST endpoint for Asana to send events in webhooks
func WebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Println("///////////////////////Recieved Webhook//////////////////")
	var event WebhookEvent
	//debug := formatRequest(r)
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		fmt.Printf("invalid request payload: %v\n", err)
		//fmt.Println(debug)
		return
	}
	events := event.Events
	fmt.Println("Got events.")
	if len(event.Errors) > 0 {
		log.Fatal(event.Errors)
	}
	//0 events should mean you're trying to create a webhook
	if len(events) == 0 {
		//IF initiating webhook, stop here and send back the x-hook-secret
		hookSecret := r.Header.Get("X-Hook-Secret")
		if hookSecret != "" { //Checking for the correct headers
			fmt.Println("Creating Webhook!")
			w.Header().Set("X-Hook-Secret", hookSecret)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("thanks"))
			return
		} else {
			fmt.Println("Invalid Webhook!")
			w.Header().Set("X-Hook-Secret", "INVALID")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(""))
			return
		}
	}
	//Since it's not invalid and we're not creating a new webhook, move on to handle events.
	supportProjects := readProjectsJSON("/home/michael/go/src/projects.json")
	var agentEmails []string
	var agentNumbers []string
	var projectId string
	var supportEmail string
	for _, e := range events {
		//Figure out the project id and the agents associated with the project.
		if e.Parent.Gid != "" {
			for _, p := range supportProjects.Projects {
				if e.Parent.Gid == p.ProjectId {
					projectId = e.Parent.Gid
					supportEmail = p.SupportEmail
					for _, a := range p.Agents {
						agentEmails = append(agentEmails, a.Email)
						agentNumbers = append(agentNumbers, a.Phone)
					}
				}
			}
			if projectId == "" {
				//If the projectId wasn't found it should be a story added so the parent will be the task Id
				task, err := getTask(e.Parent.Gid)
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, p := range supportProjects.Projects {
					//Loop through the task's projects to compare with ours
					for _, proj := range task.Projects {
						if proj.Gid == p.ProjectId {
							projectId = e.Parent.Gid
							supportEmail = p.SupportEmail
							for _, a := range p.Agents {
								agentEmails = append(agentEmails, a.Email)
								agentNumbers = append(agentNumbers, a.Phone)
							}
						}
					}
				}
			}
			if projectId != "" {
				//if we found a project id start a goroutine to handle all the events individually.
				go handleEvent(e, agentEmails, agentNumbers, projectId, supportEmail)
			}
		}
	}
	fmt.Println("///////////////////////Completed Webhook//////////////////")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("thanks"))
}

//Function called as a goroutine whenever we get a valid webhook payload
func handleEvent(e Event, recips []string, toNumbers []string, supportProjectID string, supportEmail string) {
	var eventType = e.Resource.ResourceType
	fmt.Println("Type: " + eventType)
	fmt.Println("Action: " + e.Action)
	fmt.Println("Created at: " + e.Created)
	fmt.Println("Parent: " + e.Parent.Gid)
	fmt.Println("ID: " + e.Resource.Gid)
	switch eventType {
	case "task":
		parentGID := e.Parent.Gid
		if e.Action == "added" && parentGID == supportProjectID {
			taskId := e.Resource.Gid
			task, err := getTask(taskId)
			if err != nil {
				fmt.Println(err)
				return
			}
			emails := getMessages("me", supportEmail)
			for _, e := range emails {
				readMessage("me", e.Id)
			}
			updateTaskTags(task)
		}
	case "story":
		if e.Action == "added" && e.Parent.Gid != "" {
			taskId := e.Parent.Gid
			storyId := e.Resource.Gid
			story, err := getStory(storyId)
			if err != nil {
				fmt.Println(err)
				return
			}
			urgent, err := taskIsUrgent(taskId)
			if err != nil {
				fmt.Println(err)
				return
			}
			if urgent {
				fmt.Println("URGENT TASK DETECTED")
				if story.StoryType == "added_to_tag" {
					bigTimer := startUrgentTimer(storyId, taskId, 1)
					//Add the timer to the timer array using the task id as key.
					timers = append(timers, bigTimer)
					go func() {
						for t := range bigTimer.Ticker.C {
							for _, n := range toNumbers {
								sendTwilioMessage(n, "You have an urgent support ticket that hasn't been responded to.  \n"+
									"Please reply to the original email. Then leave a comment on the task in Asana to stop the text notifications! "+
									"https://app.asana.com/0/"+supportProjectID+"/"+taskId)
								fmt.Println("Sent semi-urgent text at:", t)
							}
						}
					}()
					go func() {
						<-bigTimer.Timer.C
						bigTimer.Ticker.Stop()
						fmt.Println("big ticker stopped.")
						timers = deleteFromTimers(timers, bigTimer)
						fmt.Println("big timer deleted.")
						timer := startUrgentTimer(storyId, taskId, 2)
						timers = append(timers, timer)
						fmt.Println("Started short timer.")
						go func() {
							for t := range timer.Ticker.C {
								for _, n := range toNumbers {
									sendTwilioMessage(n, "You have an urgent support ticket that hasn't been responded to.  \n"+
										"PLEASE RESPOND OR YOU WILL BE FINED! \n"+
										"If you have already responded to the email ticket, please leave a comment on the Asana task to stop the texts: "+
										"https://app.asana.com/0/"+supportProjectID+"/"+taskId)
									fmt.Println("Sent urgent text at:", t)
								}
							}
						}()
						go func() {
							<-timer.Timer.C
							timer.Ticker.Stop()
							fmt.Println("Short ticker stopped")
							timers = deleteFromTimers(timers, timer)
							fmt.Println("Short timer deleted")
						}()
					}()
					fmt.Println("Urgent Tag Added.")
					sendEmail("You have a new urgent ticket.  Please respond to the client via the orginal email immediately.  \n\n"+
						"Please remove the software support email from the recipient list and cc important parties.  \n\n"+
						"Then leave a comment on the task in Asana to stop the urgent notifications: "+
						"https://app.asana.com/0/"+supportProjectID+"/"+taskId,
						"URGENT REQUEST PLEASE RESPOND", recips)
				}
				if story.StoryType == "comment_added" {
					fmt.Println("Comment Added")
					commenter, err := getUser(story.CreatedBy.Gid)
					if err != nil {
						fmt.Println(err)
						return
					}
					commenterEmailParts := strings.Split(commenter.Email, "@")
					commenterEmailDomain := commenterEmailParts[1]
					if commenterEmailDomain == "lamproslabs.com" {
						//STOP THE TIMERS by finding them in the array with the task id key.
						for i, t := range timers {
							fmt.Println("Timer ", i)
							if t.TaskId == taskId {
								fmt.Println("Stopping timers")
								stopTimer(t)
								timers = deleteFromTimers(timers, t)
							}
						}
						fmt.Println("Comment made by support team.")
					}

				}
			}
		}
	}
}

// formatRequest generates ascii representation of a request
// Used for debugging
// func formatRequest(r *http.Request) string {
// 	// Create return string
// 	var request []string
// 	// Add the request string
// 	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
// 	request = append(request, url)
// 	// Add the host
// 	request = append(request, fmt.Sprintf("Host: %v", r.Host))
// 	// Loop through headers
// 	for name, headers := range r.Header {
// 		name = strings.ToLower(name)
// 		for _, h := range headers {
// 			request = append(request, fmt.Sprintf("%v: %v", name, h))
// 		}
// 	}

// 	// If this is a POST, add post data
// 	if r.Method == "POST" {
// 		// Save a copy of this request for debugging.
// 		requestDump, err := httputil.DumpRequest(r, true)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		request = append(request, "\n")
// 		request = append(request, string(requestDump))
// 	}
// 	// Return the request as a string
// 	return strings.Join(request, "\n")
// }

//Helper function to read the local projects file
func readProjectsJSON(file string) *Projects {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Couldn't read agent json: %v", err)
	}
	defer f.Close()
	a := &Projects{}
	err = json.NewDecoder(f).Decode(a)
	if err != nil {
		log.Fatalf("Error decoding agent json: %v", err)
	}
	return a
}

//Not using right now
// func writeAgentsToJSON(file string) error {
// 	newAgentJSON, err := json.Marshal(agents)
// 	if err != nil {
// 		return err
// 	}
// 	err = ioutil.WriteFile(file, newAgentJSON, 0664)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Saved agents to: " + file)
// 	return nil
// }

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func StartRouter() {
	//fire up the gorilla router
	r := mux.NewRouter()

	//Endpoints
	r.HandleFunc("/receive-webhook/13r98iof2jejqeg309ihe4oq9ug3029givje", WebhookEndpoint).Methods("POST")
	r.HandleFunc("/test", TestEndpoint).Methods("GET")
	//r.HandleFunc("/add-agent", AddAgentEndpoint).Methods("POST")

	if Environment == "prod" {
		//Release the hounds
		fmt.Println("Releasing the hounds securely.")
		if err := http.ListenAndServeTLS(":4443", "/etc/letsencrypt/live/supportapi.lamproslabs.com/fullchain.pem", "/etc/letsencrypt/live/supportapi.lamproslabs.com/privkey.pem", r); err != nil {
			log.Fatal(err)
		}
	} else {
		//Release the hounds
		fmt.Println("Releasing the hounds.")
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Fatal(err)
		}
	}
}
