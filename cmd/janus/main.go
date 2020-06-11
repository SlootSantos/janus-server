package main

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/api"
	"github.com/SlootSantos/janus-server/pkg/session"
	"github.com/SlootSantos/janus-server/pkg/stacker"
	"github.com/SlootSantos/janus-server/pkg/storage"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	awsSess, err := session.AWSSession()
	if err != nil {
		log.Fatal("could not authenticate against AWS", err)
	}

	storage.Init(awsSess)

	s := stacker.New(awsSess)
	api.Start(s)
}

// func createDependendableJAMCreator() *jam.Creator {
// 	log.Print("START: setting up dependencies")

// 	awsSess, err := session.AWSSession()
// 	if err != nil {
// 		log.Fatal("could not authenticate against AWS", err)
// 	}

// 	q := queue.New(awsSess)
// 	cdn := cdn.New(awsSess, &q)
// 	bucket := bucket.New(awsSess, &q)
// 	repo := repo.New()

// 	storage.Init(awsSess)

// 	log.Print("DONE: setting up dependencies")
// 	return jam.New(bucket, cdn, repo)
// }
