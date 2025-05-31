package main

import (
    "huinong-backend/router"
)

func main() {
    router := router.InitRouter()
    router.Run(":8080")
}
