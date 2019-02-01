package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
