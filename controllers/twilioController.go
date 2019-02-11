package controllers

import (
	"encoding/json"
	"fmt"
	. "lampros-support/models"
	"strconv"
)

func SendTwilioMessage(toNumber, message string) TwilioMessageResponse {
	params := make(map[string]string)
	params["To"] = toNumber
	params["From"] = TwilioNumber
	params["Body"] = message
	respData := postTwilioRequest(params, parseUrl(TwilioBase+"/Messages.json"))
	var resp TwilioMessageResponse
	json.Unmarshal(respData, &resp)
	if resp.ErrorCode > 0 {
		fmt.Println("Twilio API Error! Code: " + strconv.Itoa(resp.ErrorCode) + " Message: " + resp.ErrorMessage)
	}
	return resp
}
