package main

import (
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/redis"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
)

func main() {
    println("Hello, World!")
    router := gin.Default()

    store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
    router.Use(sessions.Sessions("mysession", store))

    router.GET("/set", func(c *gin.Context) {
        session := sessions.Default(c)
        session.Set("username", "john_doe")
        session.Save()
        c.JSON(http.StatusOK, gin.H{"message": "Session set"})
    })

    router.GET("/get", func(c *gin.Context) {
        session := sessions.Default(c)
        username := session.Get("username")
        c.JSON(http.StatusOK, gin.H{"username": username})
    })

    router.Run(":8080")
}
