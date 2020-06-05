package pipeline

import (
	"errors"
	"fmt"
	"log"
)

type buildStream struct {
	isListening bool
	listener    chan<- string
	receiver    <-chan string
	closed      bool
}

type buildOutput struct {
	stream buildStream
	data   string
	status string
}

type buildMap map[string]*buildOutput

const StreamingCloseMessage = "::JANUS::CLOSE::STREAM::"

var builds = make(buildMap)

func GetOutput(buildID string) (string, error) {
	currBuild, ok := builds[buildID]
	if !ok {
		log.Printf("Build #%s not existing!", buildID)
		return "", fmt.Errorf("Build #%s not existing", buildID)
	}

	return currBuild.data, nil
}

func RegisterListener(buildID string, listenChan chan<- string) error {
	currBuild, ok := builds[buildID]
	if !ok {
		log.Printf("Build #%s not existing!", buildID)
		return fmt.Errorf("Build #%s not existing", buildID)
	}

	if currBuild.stream.closed {
		return errors.New(StreamingCloseMessage)
	}

	currBuild.stream.listener = listenChan
	currBuild.stream.isListening = true

	return nil
}

func StreamLogs(buildID string) {
	currBuild, ok := builds[buildID]
	if !ok {
		log.Printf("Build #%s not existing!", buildID)
		return
	}

	if currBuild.stream.isListening {
		log.Printf("Build #%s already listening to!", buildID)
		return
	}

	if currBuild.stream.closed || currBuild.status == "DONE" {
		log.Printf("Build #%s already closed!", buildID)
		return
	}

	for {
		data := <-currBuild.stream.receiver

		if data == StreamingCloseMessage {
			currBuild.status = "DONE"
			currBuild.stream.closed = true
			currBuild.stream.isListening = false

			if currBuild.stream.isListening {
				log.Println("CLOESED for listener")
				currBuild.stream.listener <- StreamingCloseMessage
			}

			break
		}

		currBuild.data += data

		// write to file
		// write to SSE
		if currBuild.stream.isListening {
			currBuild.stream.listener <- data
		}
	}

	log.Println("TOTAL DATA:::::\n", currBuild.data)
}

func CreateBuild(buildID string) (chan<- string, error) {
	log.Println("CREATING BUILD NUMBER: ", buildID)
	streamChan := make(chan string)

	newBuild := buildOutput{
		status: "STARTED",
		stream: buildStream{
			receiver: streamChan,
		},
	}

	if _, exists := builds[buildID]; exists {
		return nil, fmt.Errorf("Build #%s already exists", buildID)
	}

	builds[buildID] = &newBuild

	return streamChan, nil
}
