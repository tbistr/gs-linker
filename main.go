package main

import (
	"log"
	"net/http"

	"github.com/tbistr/gs-linker/gs"
)

func main() {
	linker := gs.New()

	http.HandleFunc("/github/events", linker.HandleGithubEvent())
	http.HandleFunc("/slack/events", linker.HandleSlackEvent())

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
