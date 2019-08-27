package main

/*
Main package starts the REST API
*/
import (
	control "lampros-support/controllers"
)

func main() {

	control.StartRouter()

}
