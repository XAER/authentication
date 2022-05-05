package main

import (
	"authentication/config"
	"authentication/helpers"
	"fmt"
	"log"

	"authentication/routes"

	"github.com/gin-gonic/gin"
)

var (
	gCnf   config.GinConfig
	router *gin.Engine
)

func initGinServer() {
	gCnf.SERVER_NAME = helpers.GetEnv("SERVER_NAME", "GO_AUTHENTICATION")
	gCnf.SERVER_PORT = helpers.GetEnv("SERVER_PORT", "8081")
	gCnf.SERVER_ENV = helpers.GetEnv("SERVER_ENV", "dev")

	fmt.Println("SERVER_NAME: ", gCnf.SERVER_NAME)

	config.SetGinMode(gCnf.SERVER_ENV)
	router = routes.SetupRouter()
}

func main() {
	initGinServer()

	db, err := config.GetConnection()
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	config.App = config.NewServices(router, &gCnf, db, gCnf.SERVER_ENV)

	log.Fatal(router.Run(fmt.Sprintf("0.0.0.0:%s", gCnf.SERVER_PORT)))
}
