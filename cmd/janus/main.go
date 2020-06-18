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
