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

// var agents *Agents
// var recips []string
// var toNumbers []string

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

//POST Request with a map of parameters.
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

//POST Request for twilio.
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
	fuck := []byte("Fuck")
	w.Header().Set("X-Hook-Secret", "12301985bugaloo120198")
	w.WriteHeader(http.StatusOK)
	w.Write(fuck)
	fmt.Println("Fuck")
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
	//TODO: figure out HMAC decode
	// hookSecret := r.Header.Get("X-Hook-Signature")
	// if !checkMAC("message", hookSecret) {
	// 	fmt.Println("Hoook secret: " + hookSecret)
	// 	fmt.Println("Sent secret doesn't match the saved one.  GTFO! " + hookSecret)
	// 	w.Header().Set("X-Hook-Secret", "FUCKYOU")
	// 	w.WriteHeader(http.StatusForbidden)
	// 	w.Write(emptyResponse)
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		fmt.Printf("invalid request payload: %v\n", err)
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
		if hookSecret != "" { //hopefully enough to stop the HAXXXORZ
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
		if e.Parent != 0 {
			for _, p := range supportProjects.Projects {
				if e.Parent == p.ProjectId {
					projectId = strconv.Itoa(e.Parent)
					supportEmail = p.SupportEmail
					for _, a := range p.Agents {
						agentEmails = append(agentEmails, a.Email)
						agentNumbers = append(agentNumbers, a.Phone)
					}
				}
			}
			if projectId == "" {
				//If the projectId wasn't found it should be a story added so the parent will be the task Id
				taskId := strconv.Itoa(e.Parent)
				task := GetTask(taskId)
				for _, p := range supportProjects.Projects {
					//Loop through the task's projects to compare with ours
					for _, proj := range task.Projects {
						if proj.Id == p.ProjectId {
							projectId = strconv.Itoa(e.Parent)
							supportEmail = p.SupportEmail
							for _, a := range p.Agents {
								agentEmails = append(agentEmails, a.Email)
								agentNumbers = append(agentNumbers, a.Phone)
							}
						}
					}
				}
			}
			//start a goroutine to hand all the events individually.
			go handleEvent(e, agentEmails, agentNumbers, projectId, supportEmail)
		}
	}
	fmt.Println("///////////////////////Completed Webhook//////////////////")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("thanks"))
}

func handleEvent(e Event, recips []string, toNumbers []string, supportProjectID string, supportEmail string) {
	fmt.Println("Type: " + e.Type)
	fmt.Println("Action: " + e.Action)
	fmt.Println("Created at: " + e.Created)
	fmt.Printf("Parent: %d\n", e.Parent)
	fmt.Printf("ID: %d\n", e.Resource)
	switch e.Type {
	case "task":
		parentGID := strconv.Itoa(e.Parent)
		if e.Action == "added" && parentGID == supportProjectID {
			taskId := strconv.Itoa(e.Resource)
			task := GetTask(taskId)
			emails := GetMessages("me", supportEmail)
			for _, e := range emails {
				subject := GetSubject(e.Id)
				sender := GetSender(e.Id)
				senderDomain := strings.Split(sender, "@")
				senderArr := []string{sender} //Needed for SendEmail
				if !IsAsanaDomain(senderDomain[1]) {
					if !CheckProjectEmail(sender, supportProjectID) {
						SendEmail("Please add the new user email: "+sender+" to the support project: https://app.asana.com/0/"+supportProjectID+"\nThen add them to the request: https://app.asana.com/0/"+supportProjectID+"/"+taskId, "New User Detected for Support.", recips)
						SendEmail("Thank you for your for your request to Lampros Support. \nWe do not recognize your email.  You will need to be added to Asana to recieve support notifications. \nWe will confirm your email and add you to your support project. \nWe will contact you directly if we need more information. \n\n Thank you, \n\n -Lampros Labs Team \n", "New Software Support Request", senderArr)
					} else {
						fmt.Println("Follower found: " + sender)
						if subject == task.Name {
							fmt.Println("Trying to add to task")
							UpdateTaskFollowers(sender, task.Gid)
						}
					}
					SendEmail("You have a new support ticket please leave a comment on the asana ticket to respond and/or assign the task to yourself: https://app.asana.com/0/"+supportProjectID+"/"+taskId, "New Software Support Ticket: "+task.Name, recips)
				}
				ReadMessage("me", e.Id)
			}
			UpdateTaskTags(task)
		}
	case "story":
		if e.Action == "added" && e.Parent != 0 {
			taskId := strconv.Itoa(e.Parent)
			storyId := strconv.Itoa(e.Resource)
			story := GetStory(storyId)
			if TaskIsUrgent(taskId) {
				fmt.Println("URGENT TASK DETECTED")
				if story.StoryType == "added_to_tag" {
					bigTimer := StartUrgentTimer(e.Resource, e.Parent, 1)
					//Add the timer to the timer array using the task id as key.
					timers = append(timers, bigTimer)
					go func() {
						for t := range bigTimer.Ticker.C {
							for _, n := range toNumbers {
								SendTwilioMessage(n, "You have an urgent support ticket that hasn't been responded to.  Please check your email and respond! https://app.asana.com/0/"+supportProjectID+"/"+taskId)
								fmt.Println("Sent semi-urgent text at:", t)
							}
						}
					}()
					go func() {
						<-bigTimer.Timer.C
						bigTimer.Ticker.Stop()
						fmt.Println("big ticker stopped.")
						timers = DeleteFromTimers(timers, bigTimer)
						fmt.Println("big timer deleted.")
						timer := StartUrgentTimer(e.Resource, e.Parent, 2)
						timers = append(timers, timer)
						fmt.Println("Started short timer.")
						go func() {
							for t := range timer.Ticker.C {
								for _, n := range toNumbers {
									SendTwilioMessage(n, "You have an urgent support ticket that hasn't been responded to.  PLEASE RESPOND OR YOU WILL BE FINED! https://app.asana.com/0/"+supportProjectID+"/"+taskId)
									fmt.Println("Sent urgent text at:", t)
								}
							}
						}()
						go func() {
							<-timer.Timer.C
							timer.Ticker.Stop()
							fmt.Println("Short ticker stopped")
							timers = DeleteFromTimers(timers, timer)
							fmt.Println("Short timer deleted")
						}()
					}()
					fmt.Println("Urgent Tag Added.")
					SendEmail("You have a new urgent ticket please respond immediately: https://app.asana.com/0/"+supportProjectID+"/"+taskId, "URGENT REQUEST PLEASE RESPOND", recips)
				}
				if story.StoryType == "comment_added" {
					fmt.Println("Comment Added")
					commenter := GetUser(story.CreatedBy.Gid)
					commenterEmailParts := strings.Split(commenter.Email, "@")
					commenterEmailDomain := commenterEmailParts[1]
					if commenterEmailDomain == "lamproslabs.com" {
						//STOP THE TIMERS by finding them in the array with the task id key.
						for i, t := range timers {
							fmt.Println("Timer ", i)
							if t.TaskId == e.Parent {
								fmt.Println("Stopping timers")
								StopTimer(t)
								timers = DeleteFromTimers(timers, t)
							}
						}
						fmt.Println("Comment made by support team.")
					}
				}
			}
		}
	}
}

// func checkMAC(message, sentMAC string) bool {
// 	mac := hmac.New(sha256.New, []byte(SavedHookSecret))
// 	mac.Write([]byte(message))
// 	expectedMAC := mac.Sum(nil)
// 	fmt.Printf("%x\n", expectedMAC)
// 	fmt.Printf("%x\n", []byte(sentMAC))
// 	return hmac.Equal([]byte(sentMAC), expectedMAC)
// }

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
	//r.HandleFunc("/test", TestEndpoint).Methods("GET")
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
