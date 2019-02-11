package models

type TwilioMessageResponse struct {
	Id           string  `json:"sid"`
	Created      string  `json:"date_created"` //Mon, 11 Feb 2019 21:28:25 +0000
	Updated      string  `json:"date_updated"`
	Sent         string  `json:"date_sent"`
	AccountId    string  `json:"account_sid"`
	To           string  `json:"to"`
	From         string  `json:"from"`
	ServiceId    string  `json:"messaging_service_sid"`
	Body         string  `json:"body"`
	Status       string  `json:"status"`
	NumSegments  string  `json:"num_segments"`
	NumMedia     string  `json:"num_media"`
	Direction    string  `json:"direction"`
	ApiVersion   string  `json:"api_version"`
	Price        float32 `json:"price"`
	PriceUnit    string  `json:"price_unit"`
	ErrorCode    int     `json:"error_code"`
	ErrorMessage string  `json:"error_message"`
	URI          string  `json:"uri"`
	SubURI       []Media `json:"subresource_uris"`
}

type Media struct {
	URI string `json:"media"`
}
