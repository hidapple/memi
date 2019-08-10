package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hidapple/memi/bot"
	"github.com/hidapple/memi/kibela"
	"golang.org/x/xerrors"
)

var help = fmt.Sprintf(`
%s
ping:
Test the reachability of the bot.
> @memi ping

link:
Append markdown link to the Kibela note.
> @memi link $URL [$TITLE $TEXT...]
%s`, "```", "```")

func main() {
	b := bot.NewBot(os.Getenv("SLACK_TOKEN"), help)
	b.AddCommands(
		&bot.Command{
			Name: "ping",
			Do: func(args ...string) (string, error) {
				return "にゃん:two_hearts:", nil
			},
		},
		&bot.Command{
			Name: "link",
			Do: func(args ...string) (string, error) {
				if len(args) == 0 {
					return "こうやって使ってね:point_right: `@memi link $URL [$TITLE $TEXT...]`", nil
				}
				var url, title string
				url = args[0]
				if len(args) == 1 {
					title = url
				} else {
					title = strings.Join(args[1:], " ")
				}
				k := kibela.New(os.Getenv("KIBELA_TOKEN"), os.Getenv("KIBELA_TEAM"))
				note, err := k.AddLink(os.Getenv("KIBELA_LINK_NOTE"), url, title)
				if err != nil {
					return "", xerrors.Errorf("link command failed: %s", err)
				}
				msg := "更新したよ〜:cat:\n"
				msg += "```\n"
				msg += note.Content
				msg += "```"
				return msg, nil
			},
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := b.OpenRTM(ctx); err != nil {
		log.Fatalf("memi is stopped by err: %s", err)
	}
}
