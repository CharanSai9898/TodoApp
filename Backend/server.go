package main

import (
	"log"

	"todo/src"
)

func main() {
	app := src.SetupApp()

	port := ":4000"
	log.Println("Server Started on Port " + port)
	app.Listen(port)
}
