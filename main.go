package main

import (
	"github.com/gin-gonic/gin"
	"workWithDB/handlers"
)

func main() {
	router := gin.Default()
	router.POST("/create_user", handlers.CreateUser)
	router.DELETE("/delete_user/:id", handlers.DeleteUser)
	router.POST("/sign_in", handlers.Entry)
	//
	router.GET("/user", handlers.Something)
	router.Run(":8001")

}
