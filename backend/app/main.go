package main

import (
	"huinongfinancial/router"
)

func main() {
	router := router.InitRouter()
	router.Run(":8080")
}