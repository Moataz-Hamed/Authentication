package main

import (
	"Auth/controllers"
	"Auth/initializers"
	"Auth/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
	initializers.SyncDatabase()

}

func main() {
	r := gin.Default()

	r.POST("/login", controllers.Login)
	r.POST("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/signup", controllers.SignUp)
	r.Run()

}
