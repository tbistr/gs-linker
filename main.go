package main

import (
	"context"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/google/go-github/v45/github"
	gh "github.com/tbistr/gs-linker/github"
	"github.com/tbistr/gs-linker/link"
	sl "github.com/tbistr/gs-linker/slack"
)

type envs struct {
	GhAppID          int64  `env:"GH_APP_ID"`
	GhInstallationID int64  `env:"GH_INSTALLATION_ID"`
	SlToken          string `env:"SL_TOKEN"`
	SlSigningSecret  string `env:"SL_SIGNING_SECRET"`
	SlBotUserID      string `env:"SL_BOT_USER_ID"`
}

func main() {
	e := envs{}
	if err := env.Parse(&e); err != nil {
		log.Fatalln(err)
	}
	log.Printf("get envs: %+v\n", e)

	ghClient := gh.New(e.GhAppID, e.GhInstallationID)
	slClient := sl.New(e.SlToken, e.SlSigningSecret, e.SlBotUserID)

	links := link.Links{}
	var (
		onCommented gh.OnCommentedFunc = func(client *gh.Client, thread *gh.Thread, comment *github.IssueComment) error {
			s, err := links.SearchByG(thread)
			if err != nil {
				return err
			}
			return slClient.SendMsg(s, comment.GetBody())
		}

		handleSub sl.HandleSubFunc = func(client *sl.Client, thread *sl.Thread, rawURL string) {
			owner, repo, num, err := ghClient.VerifyURL(rawURL)
			if err != nil {
				log.Println(err)
				return
			}
			gThread := &gh.Thread{
				Owner: owner,
				Repo:  repo,
				Num:   num,
			}
			if err := links.Sub(gThread, thread); err != nil {
				log.Println(err)
				return
			}
		}
		handleUnsub sl.HandleUnsubFunc = func(client *sl.Client, thread *sl.Thread) {
			if err := links.UnSub(thread); err != nil {
				log.Println(err)
				return
			}
		}
		// handleSummary sl.HandleSummaryFunc = func(client *sl.Client, thread *sl.Thread) {}

		onMsgSent sl.OnMsgSentFunc = func(client *sl.Client, thread *sl.Thread, text string) {
			g, err := links.SearchByS(thread)
			if err != nil {
				log.Println(err)
				return
			}
			ghClient.CreateComment(context.Background(), g, text)
		}
	)

	ghClient.RegisterOnCommented(onCommented)
	slClient.RegisterHandleSub(handleSub)
	slClient.RegisterHandleUnsub(handleUnsub)
	// slClient.RegisterHandleSummary(handleSummary)
	slClient.RegisterOnMsgSent(onMsgSent)

	http.HandleFunc("/github/events", ghClient.HandleEvent())
	http.HandleFunc("/slack/events", slClient.HandleEvent())

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
