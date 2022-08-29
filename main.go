package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/google/go-github/v45/github"
	gh "github.com/tbistr/gs-linker/github"
	"github.com/tbistr/gs-linker/link"
	sl "github.com/tbistr/gs-linker/slack"
)

type envs struct {
	GhAppID          int64 `env:"GH_APP_ID"`
	GhInstallationID int64 `env:"GH_INSTALLATION_ID"`

	SlToken         string `env:"SL_TOKEN"`
	SlSigningSecret string `env:"SL_SIGNING_SECRET"`
	SlBotUserID     string `env:"SL_BOT_USER_ID"`

	DBName string `env:"MYSQL_DATABASE"`
	DBUser string `env:"MYSQL_USER"`
	DBPass string `env:"MYSQL_PASSWORD"`
	DBAddr string `env:"MYSQL_ADDR"`
}

func main() {
	e := envs{}
	if err := env.Parse(&e); err != nil {
		log.Fatalln(err)
	}
	log.Printf("get envs: %+v\n", e)

	ghClient := gh.New(e.GhAppID, e.GhInstallationID)
	slClient := sl.New(e.SlToken, e.SlSigningSecret, e.SlBotUserID)

	linkerConf := &link.DbConfig{
		UserName: e.DBUser,
		Pass:     e.DBPass,
		Addr:     e.DBAddr,
		DbName:   e.DBName,
	}
	linker := link.New(linkerConf)

	// Ideally, the behavier should be implemented here and each liblary (gh, sl, link) should only provide functions.
	// For example, response messages (ex. "already subscribed!!") should be written in below handlers.
	// And more, functional logs (ex. "link stored in DB successfully") should be thrown from liblary.
	var (
		onCommented gh.OnCommentedFunc = func(client *gh.Client, thread *gh.Thread, comment *github.IssueComment) error {
			if s := linker.SearchByG(thread); s != nil {
				return slClient.SendMsg(s, comment.GetBody())
			}
			return fmt.Errorf("link not found")
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
			if err := linker.Sub(gThread, thread); err != nil {
				log.Println(err)
				return
			}
		}
		handleUnsub sl.HandleUnsubFunc = func(client *sl.Client, thread *sl.Thread) {
			if err := linker.UnSub(thread); err != nil {
				log.Println(err)
				return
			}
		}
		// handleSummary sl.HandleSummaryFunc = func(client *sl.Client, thread *sl.Thread) {}

		onMsgSent sl.OnMsgSentFunc = func(client *sl.Client, thread *sl.Thread, text string) {
			if g := linker.SearchByS(thread); g != nil {
				ghClient.CreateComment(context.Background(), g, text)
				return
			}
			fmt.Printf("link not found\n")
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
