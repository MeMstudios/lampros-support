package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
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
	//EMAIL RECIPIENTS
	recips := []string{"michael@lamproslabs.com"}
	for _, e := range events {
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
					if !CheckProjectEmail(sender) && senderDomain[1] != "asana.com" && senderDomain[1] != "mail.asana.com" {
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
				if TaskIsUrgent(taskId) {
					fmt.Println("URGENT TASK DETECTED")
				}
			}
		case "webhook":
			hookSecret := r.Header.Get("X-Hook-Signature")
			fmt.Println("Hoook secret: " + hookSecret)
			w.Header().Set("X-Hook-Secret", hookSecret)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
			return
		}
	}
	fmt.Println("///////////////////////Completed Webhook//////////////////")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("thanks"))
	return
}

func checkMAC(message, sentMAC string) bool {
	mac := hmac.New(sha256.New, []byte(SavedHookSecret))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	fmt.Printf("%x\n", expectedMAC)
	fmt.Printf("%x\n", []byte(sentMAC))
	return hmac.Equal([]byte(sentMAC), expectedMAC)
}

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
		if err := http.ListenAndServeTLS(":443", "fullchain.pem", "privkey.pem", r); err != nil {
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
