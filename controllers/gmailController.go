package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

//  Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "/home/michael/go/src/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

//  Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

//  Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

//  Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func startGmailClient() *gmail.Service {
	b, err := ioutil.ReadFile("/home/michael/go/src/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return srv
}

func getSender(id string) string {
	api := startGmailClient()
	user := "me"
	meta, err := api.Users.Messages.Get(user, id).Format("metadata").Do()
	if err != nil {
		log.Fatalf("Failed to get metadata")
	}
	var from string
	for _, h := range meta.Payload.Headers {
		if h.Name == "From" {
			fromSplit := strings.Split(h.Value, "<")
			email := strings.Split(fromSplit[1], ">")
			from = email[0]
			fmt.Println("From: " + from)
		}
	}
	return from
}

//  returns all the unread messages in the users inbox sent to a certain email address
func getMessages(user string, supportEmailAddress string) []*gmail.Message {
	api := startGmailClient()
	mesList, err := api.Users.Messages.List(user).Q("to:" + supportEmailAddress + " is:unread").Do()
	if err != nil {
		log.Fatalf("Failed to get messages: %v", err)
	}
	if len(mesList.Messages) == 0 {
		fmt.Println("No messages to " + supportEmailAddress)
	}
	return mesList.Messages
}

//  removes a message from the "unread" list in gmail by id
func readMessage(user string, id string) {
	api := startGmailClient()
	var req gmail.ModifyMessageRequest
	req.RemoveLabelIds = append(req.AddLabelIds, "UNREAD")
	res, err := api.Users.Messages.Modify(user, id, &req).Do()
	if err != nil {
		fmt.Printf("Error from Gmail: %v\n", err)
	} else {
		fmt.Println("Message read: " + res.Id)
	}
}

func getSubject(id string) string {
	api := startGmailClient()
	user := "me"
	meta, err := api.Users.Messages.Get(user, id).Format("metadata").Do()
	if err != nil {
		log.Fatalf("Failed to get metadata")
	}
	subj := ""
	for _, h := range meta.Payload.Headers {
		if h.Name == "Subject" {
			subj = h.Value
		}
	}
	return subj
}

func sendEmail(body, subj string, recip []string) bool {
	auth := smtp.PlainAuth("", EmailAddress, Password, "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email
	to := recip
	toHeader := "To: "
	for _, r := range to {
		toHeader += r + ","
		fmt.Println("Sending message to: " + r)
	}
	msg := []byte(toHeader + "\r\n" +
		"Subject: " + subj + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, EmailAddress, to, msg)
	if err != nil {
		log.Printf("Error sending email: %e", err)
		return false
	}
	return true
}
