package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/pipeline"
	"github.com/SlootSantos/janus-server/pkg/repo"
	"github.com/buildkite/terminal-to-html"
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
		w.Write([]byte("ALIVE :)"))
	})

	http.HandleFunc("/sse/build", func(w http.ResponseWriter, req *http.Request) {
		urlParams := req.URL.Query()
		idQueries, ok := urlParams["id"]
		if !ok {
			io.WriteString(w, "Invalid Query. Missing \"id\"")
			return
		}

		buildID := idQueries[0]

		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("X-Accel-Buffering", "no")
		flusher := w.(http.Flusher)

		initalBuildOutput, _ := pipeline.GetOutput(buildID)
		initialHTML := terminal.Render([]byte(initalBuildOutput))
		initialBuildEscaped := strings.Replace(string(initialHTML), "\n", "<br/>", -1)
		fmt.Fprintf(w, "data: %s\n\n", initialBuildEscaped)
		flusher.Flush()

		listenChan := make(chan string)
		err := pipeline.RegisterListener(buildID, listenChan)

		if err != nil {
			if err.Error() == pipeline.StreamingCloseMessage {
				return
			}
		}

		for {
			time.Sleep(time.Millisecond * 200)
			d := <-listenChan

			if d == pipeline.StreamingCloseMessage {
				log.Println("CLOESED for listener")

				break
			}

			go func(data string) {
				htmlData := terminal.Render([]byte(data))
				htmlEscaped := strings.Replace(string(htmlData), "\n", "<br/>", -1)

				fmt.Fprintf(w, "data: %s\n\n", htmlEscaped)
				flusher.Flush()
			}(d)
		}

	})

	log.Println("DONE: setting up API routes")
	log.Println("LISTEN: :8888")

	if os.Getenv("IS_ENTERPRISE") != "" {
		log.Println("\n----\n Running properitary Janus Version including Payment options!\n----")
	}

	http.ListenAndServe(":8888", nil)
}
