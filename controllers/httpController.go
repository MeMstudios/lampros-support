package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	if len(events) == 0 {
		//IF initiating webhook, stop here and send back the x-hook-secret
		fmt.Println("Creating Webhook!")
		hookSecret := r.Header.Get("X-Hook-Secret")
		w.Header().Set("X-Hook-Secret", hookSecret)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("thanks"))
		return
	}
	for _, e := range events {
		//Otherwise, start a goroutine to hand all the events individually.
		go handleEvent(e)
	}

	fmt.Println("///////////////////////Completed Webhook//////////////////")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("thanks"))
}

func handleEvent(e Event) {
	fmt.Println("Type: " + e.Type)
	fmt.Println("Action: " + e.Action)
	fmt.Println("Created at: " + e.Created)
	fmt.Printf("Parent: %d\n", e.Parent)
	fmt.Printf("ID: %d\n", e.Resource)
	switch e.Type {
	case "task":
		if e.Action == "added" {
			taskId := strconv.Itoa(e.Resource)
			task := GetTask(string(taskId))
			emails := GetMessages("me")
			for _, e := range emails {
				subject := GetSubject(e.Id)
				sender := GetSender(e.Id)
				senderDomain := strings.Split(sender, "@")
				senderArr := []string{sender} //Needed for SendEmail
				if !CheckProjectEmail(sender) && !IsAsanaDomain(senderDomain[1]) {
					SendEmail("Please add the new user email: "+sender+" to the support project: https://app.asana.com/0/"+SupportProjectID+"\nThen add them to the request: https://app.asana.com/0/"+SupportProjectID+"/"+taskId, "New User Detected for Support.", recips)
					SendEmail("Thank you for your for your request to Lampros Support. \nWe do not recognize your email.  You will need to be added to Asana to recieve support notifications. \nWe will confirm your email and add you to your support project. \nWe will contact you directly if we need more information. \n\n Thank you, \n\n -Lampros Labs Team \n", "New Software Support Request", senderArr)
				} else {
					fmt.Println("Follower found: " + sender)
					if subject == task.Name {
						fmt.Println("Trying to add to task")
						UpdateTaskFollowers(sender, task.Gid)
					}
				}
			}
			UpdateTaskTags(task)
		}
	case "story":
		if e.Action == "added" && e.Parent != 0 {
			taskId := strconv.Itoa(e.Parent)
			storyId := strconv.Itoa(e.Resource)
			if TaskIsUrgent(taskId) {
				fmt.Println("URGENT TASK DETECTED")
				story := GetStory(storyId)
				if story.StoryType == "added_to_tag" {
					bigTimer := StartUrgentTimer(e.Resource, e.Parent, 1)
					//Add the timer to the timer array using the task id as key.
					timers = append(timers, bigTimer)
					go func() {
						for t := range bigTimer.Ticker.C {
							for _, n := range toNumbers {
								SendTwilioMessage(n, "You have an urgent support ticket that hasn't been responded to.  Please check your email and respond! https://app.asana.com/0/"+SupportProjectID+"/"+taskId)
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
									SendTwilioMessage(n, "You have an urgent support ticket that hasn't been responded to.  PLEASE RESPOND OR YOU WILL BE FINED! https://app.asana.com/0/"+SupportProjectID+"/"+taskId)
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
					SendEmail("You have a new urgent ticket please respond immediately: https://app.asana.com/0/"+SupportProjectID+"/"+taskId, "URGENT REQUEST PLEASE RESPOND", recips)
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
	r.HandleFunc("/receive-webhook/support-test", WebhookEndpoint).Methods("POST")
	r.HandleFunc("/test", TestEndpoint).Methods("GET")

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
