package main

import (
	control "lampros-support/controllers"
)

func main() {

	//control.StartRouter()
	control.SendTwilioMessage("+18592402898", "fuck you from the app.")
}
