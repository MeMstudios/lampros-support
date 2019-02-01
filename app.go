package main

import (
	control "lampros-support/controllers"
)

func main() {
	tasks := control.GetTasks()
	control.UpdateTasks(tasks)
}
