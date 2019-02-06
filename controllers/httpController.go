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

//GET a JSON response from an API
func getResponse(request string) []byte {
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
func postRequest(params map[string]string, request string) []byte {
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
	var emptyResponse []byte
	var event WebhookEvent
	hookSecret := r.Header.Get("X-Hook-Secret")
	if hookSecret != SavedHookSecret {
		fmt.Println("Sent secret doesn't match the saved one.  GTFO! " + hookSecret)
		w.Header().Set("X-Hook-Secret", "FUCKYOU")
		w.WriteHeader(http.StatusForbidden)
		w.Write(emptyResponse)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	events := event.Events

	//EMAIL RECIPIENTS
	recips := []string{"michael@lamproslabs.com"}
	for _, e := range events {
		switch e.Type {
		case "task":
			{
				task := GetTask(string(e.Resource))
				emails := GetMessages("me")
				for _, e := range emails {
					subject := GetSubject(e.Id)
					sender := GetSender(e.Id)
					user, err := GetUserByEmail(sender)
					if err != nil {
						SendEmail("Please add the new user email: "+sender+" to the support project: https://app.asana.com/0/"+SupportProjectID, "New User Detected for Support.", recips)
					} else {
						fmt.Println("User found: " + user.Gid)
						if subject == task.Name {
							UpdateTaskFollowers(user.Email, task.Gid)
						}
					}
				}
			}
		case "project":
			{

			}
		case "webhook":
			{
				fmt.Println("Hoook secret: " + hookSecret)
				//fmt.Println(ioutil.ReadAll(r.Body))
				w.Header().Set("X-Hook-Secret", hookSecret)
				w.WriteHeader(http.StatusOK)
				w.Write(emptyResponse)
			}
		}
	}

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
