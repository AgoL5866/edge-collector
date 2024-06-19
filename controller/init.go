package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/coolestowl/edge-collector/cmd"
	"github.com/coolestowl/edge-collector/configs"
	v1 "github.com/coolestowl/edge-collector/controller/v1"

	"github.com/gin-gonic/gin"
)

func debug(c *gin.Context) {
	for k, v := range c.Request.Header {
		log.Println("req:", k, strings.Join(v, ";"))
	}

	defer func() {
		for k, v := range c.Writer.Header() {
			log.Println("res:", k, strings.Join(v, ";"))
		}
	}()

	c.Next()
}

func Init(e *gin.Engine) error {
	e.Use(debug)

	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/echo", v1.Echo)

		apiV1.POST("/req", v1.Req)

		apiV1.POST("/routine", func(c *gin.Context) {
			m := &cmd.Monitor{}
			m.Routine(configs.App)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "ok"})
		})
	}

	return nil
}
