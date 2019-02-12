package main

/*
At this time it does not check for if a COMMENT was added by a customer that could indicate the task is urgent.
The urgent tag will only be set automatically if someone emails to the support email with a subject or body containing any form of the word 'urgent'
However, the urgent response texts/emails will start to get sent if the urgent tag is set manually.
*/
import (
	control "lampros-support/controllers"
)

func main() {

	control.StartRouter()

}
