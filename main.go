package main

import (
	"go-jwt-challenge/database"
	"go-jwt-challenge/router"
)

func main() {

	database.StartDB()
	r := router.StartApp()
	r.Run(":8080")

}
