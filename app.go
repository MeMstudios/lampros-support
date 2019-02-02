package main

import (
	control "lampros-support/controllers"
)

func main() {
	/*
		tasks := control.GetTasks()
		control.UpdateTasks(tasks)
	*/
	senders := control.GetSenders()
	recips := []string{"michaelmurphystudios@gmail.com"}
	control.SendEmail("Fuck you again", "Test Email", recips)

}
