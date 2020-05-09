package jam

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type jamPOSTBody struct {
	Repository string
}

// RoutePrefix is the JAM endpoint
const RoutePrefix = "/jam"

// HandleHTTP handles the JAM http endpoint
func (j *Creator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	methodHandlerMap := map[string]http.HandlerFunc{
		http.MethodGet:    j.handleGET,
		http.MethodPost:   j.handlePOST,
		http.MethodDelete: j.handleDELETE,
	}

	if handler, ok := methodHandlerMap[req.Method]; ok {
		handler(w, req)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (j *Creator) handlePOST(w http.ResponseWriter, req *http.Request) {
	var jb jamPOSTBody

	err := json.NewDecoder(req.Body).Decode(&jb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedList, err := j.build(req.Context(), jb.Repository)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(updatedList)
}

func (j *Creator) handleGET(w http.ResponseWriter, req *http.Request) {
	stacks, _ := j.list(req.Context())
	stackJSON, err := json.Marshal(stacks)
	if err != nil {
		log.Println(err)
	}

	w.Write(stackJSON)
}

func (j *Creator) handleDELETE(w http.ResponseWriter, req *http.Request) {
	urlParams := req.URL.Query()
	idQueries, ok := urlParams["id"]
	if !ok {
		io.WriteString(w, "Invalid Query. Missing \"id\"")
		return
	}

	bucketID := idQueries[0]
	newlist, err := j.delete(req.Context(), bucketID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(newlist)
}
