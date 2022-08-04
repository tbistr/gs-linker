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

		onMentioned sl.OnMentionedFunc = func(client *sl.Client, thread *sl.Thread, text string) error {
			// TODO: parse msg.
			gThread := gh.Thread{
				Owner: "tbistr",
				Repo:  "gs-linker",
				Num:   8,
			}
			return links.Sub(&gThread, thread)
		}

		onMsgSent sl.OnMsgSentFunc = func(client *sl.Client, thread *sl.Thread, text string) error {
			g, err := links.SearchByS(thread)
			if err != nil {
				return err
			}
			ghClient.CreateComment(context.Background(), g, text)
			return nil
		}
	)

	ghClient.RegisterOnCommented(onCommented)
	slClient.RegisterOnMentioned(onMentioned)
	slClient.RegisterOnMsgSent(onMsgSent)

	http.HandleFunc("/github/events", ghClient.HandleEvent())
	http.HandleFunc("/slack/events", slClient.HandleEvent())

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
