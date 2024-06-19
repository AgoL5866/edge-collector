package api

import (
	"net/http"

	"github.com/coolestowl/edge-collector/server"
)

func init() {
	server.Init()
}

func Entry(w http.ResponseWriter, r *http.Request) {
	server.Serve(w, r)
}
