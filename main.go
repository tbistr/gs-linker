package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/k0kubun/pp/v3"
	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

type envs struct {
	GhAppID          int64  `env:"GH_APP_ID"`
	GhInstallationID int64  `env:"GH_INSTALLATION_ID"`
	SlToken          string `env:"SL_TOKEN"`
	SlSigningSecret  string `env:"SL_SIGNING_SECRET"`
}

func main() {
	e := envs{}
	if err := env.Parse(&e); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Get envs: %+v\n", e)

	ghClient := gh.New(e.GhAppID, e.GhInstallationID)
	slClient := sl.New(e.SlToken, e.SlSigningSecret)

	var (
		// onIssueCommented gh.OnIssueCommentedFunc = func(client *gh.Client, owner, repo string, num int, comment *github.IssueComment) error {
		// 	return nil
		// }
		// onPrCommented gh.OnPrCommentedFunc = func(client *gh.Client, owner, repo string, num int, comment *github.PullRequestComment) error {
		// 	return nil
		// }

		onMentioned sl.OnMentionedFunc = func(client *sl.Client, channel, threadTS, text string) error {
			fmt.Println("Mentioned:")
			pp.Println(channel)
			pp.Println(threadTS)
			pp.Println(text)
			return nil
		}
		onMsgSent sl.OnMsgSentFunc = func(client *sl.Client, channel, threadTS, text string) error {
			fmt.Println("MsgSent:")
			pp.Println(threadTS)
			pp.Println(text)
			return nil
		}
	)
	slClient.RegisterOnMentioned(onMentioned)
	slClient.RegisterOnMsgSent(onMsgSent)

	http.HandleFunc("/github/events", ghClient.HandleEvent())
	http.HandleFunc("/slack/events", slClient.HandleEvent())

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
