package api

import (
	"log"
	"net/http"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/pipeline"
	"github.com/SlootSantos/janus-server/pkg/repo"
)

// Start sets up the HTTP endpoints for Janus
func Start(j *jam.Creator) {
	log.Println("START: setting up API routes")

	http.HandleFunc(repo.RoutePrefix, auth.WithCredentials(repo.HandleHTTP))
	http.HandleFunc(jam.RoutePrefix, auth.WithCredentials(j.ServeHTTP))
	http.HandleFunc(auth.LoginCheck, auth.HandleLoginCheck)
	http.HandleFunc(auth.Callback, auth.HandleCallback)
	http.HandleFunc(pipeline.Hook, pipeline.HandleHook)
	http.HandleFunc(auth.Login, auth.HandleLogin)
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ALIVE"))
	})
	log.Println("DONE: setting up API routes")
	log.Println("LISTEN: :8888")

	http.ListenAndServe(":8888", nil)
}
