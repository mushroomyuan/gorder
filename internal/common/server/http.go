package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		panic("empty http address")
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)
	apiRouter.Group("/api")
	apiRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"response": "pong"})
	})
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}
