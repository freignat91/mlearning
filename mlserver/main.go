package main

import (
	"log"
	"os"

	"github.com/freignat91/mlearning/mlserver/server"
)

// build vars
var (
	Version string
	Build   string
)

func main() {
	server := mlserver.Server{}
	err := server.Start(Version)
	if err != nil {
		log.Printf("Exit on init error: %v\n", err)
		os.Exit(1)
	}
}
