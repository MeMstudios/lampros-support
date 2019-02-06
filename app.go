package main

import (
	control "lampros-support/controllers"
)

func main() {

	// hookResponse := control.CreateWebhook()
	// fmt.Printf("Id: %d \n", hookResponse.Webhook.Id)
	// fmt.Println("Target: " + hookResponse.Webhook.Target)
	control.StartRouter()

	/*
		tasks := control.GetTasks()
		control.UpdateTaskTags(tasks)

		recips := []string{"michael@lamproslabs.com"}
		senders := control.GetSenders()
		for _, s := range senders {
			user, err := control.GetUserByEmail(s)
			if err != nil {
				control.SendEmail("Please add the new user email: "+s+" to the support project: https://app.asana.com/0/"+control.SupportProjectID, "New User Detected for Support.", recips)
			} else {
				fmt.Println("User found: " + user.Gid)
				control.UpdateTaskFollowers(s, "1107883497327024")
			}
		}
	*/
	//control.SendEmail("Thank you for your request, we will be with you shortly.", "Your Support Request Has Been Revieved", senders)

}
