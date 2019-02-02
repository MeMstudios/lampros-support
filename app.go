package main

import (
	"fmt"
	control "lampros-support/controllers"
)

func main() {
	/*
		tasks := control.GetTasks()
		control.UpdateTasks(tasks)
	*/
	recips := []string{"michael@lamproslabs.com"}
	senders := control.GetSenders()
	for _, s := range senders {
		user, err := control.GetUserByEmail(s)
		if err != nil {
			control.SendEmail("Please add the user email to the support project: "+s, "New User Detected for Support.", recips)
		} else {
			fmt.Println("User found: " + user.Gid)
		}
	}
	//control.SendEmail("Thank you for your request, we will be with you shortly.", "Your Support Request Has Been Revieved", senders)

}
