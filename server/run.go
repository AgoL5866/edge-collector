package server

import (
	"log"
	"net/http"

	"github.com/coolestowl/edge-collector/controller"
	"github.com/coolestowl/edge-collector/env"

	"github.com/gin-gonic/gin"
)

var (
	eng *gin.Engine
)

func Init() {
	for _, f := range []func() error{
		env.Init,

		MonitorInit,

		func() error {
			app := gin.Default()
			app.RemoteIPHeaders = append(app.RemoteIPHeaders, "Cf-Connecting-Ip")

			defer func() {
				gin.SetMode(gin.ReleaseMode)
				eng = app
			}()

			if err := controller.Init(app); err != nil {
				return err
			}

			return nil
		},
	} {
		if err := f(); err != nil {
			log.Fatalln("[init]", err)
		}
	}
}

func Engine() *gin.Engine {
	return eng
}

func Serve(w http.ResponseWriter, r *http.Request) {
	eng.ServeHTTP(w, r)
}
