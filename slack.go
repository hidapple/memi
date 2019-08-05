package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hidapple/memi/kibela"
	"github.com/nlopes/slack"
)

// SlackBot respresents Slack Bot App
type SlackBot struct {
	client *slack.Client
	rtm    *slack.RTM
	botID  string
}

// OpenRTM opens RTM connection and listens incomming events
func (s *SlackBot) OpenRTM() {
	s.rtm = s.client.NewRTM()

	go s.rtm.ManageConnection()
	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			s.botID = ev.Info.User.ID

		case *slack.MessageEvent:
			if err := s.handleMessage(ev); err != nil {
				log.Printf("failed to handle message: %s", err)
			}

		case *slack.RTMError:
			log.Printf("RTM Error. code=%d, msg=%s", ev.Code, ev.Msg)
		}
	}
}

func (s *SlackBot) handleMessage(ev *slack.MessageEvent) error {
	// Input is supposed to be `$bot_name $command [argument(s)]`
	msg := strings.Split(ev.Msg.Text, " ")

	// Ignore if input is not bot mention
	if !strings.HasPrefix(msg[0], fmt.Sprintf("<@%s>", s.botID)) {
		return nil
	}

	// Print help when command is not given
	if len(msg) < 2 {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(s.printHelp(), ev.Channel))
		return nil
	}

	switch msg[1] {
	case "ping":
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("にゃん:two_hearts:", ev.Channel))

	case "link":
		if len(msg) < 3 {
			s.rtm.SendMessage(s.rtm.NewOutgoingMessage("使い方: `link $URL [$TITLE]`", ev.Channel))
			return nil
		}

		var url, title string
		url = msg[2]
		if len(msg) > 3 {
			title = strings.Join(msg[3:], " ")
		} else {
			title = msg[2]
		}
		k := kibela.New(os.Getenv("KIBELA_TOKEN"), os.Getenv("KIBELA_TEAM"))
		note, err := k.AddLink(os.Getenv("KIBELA_LINK_NOTE"), url, title)
		if err != nil {
			return err
		}
		msg := "更新したよ〜:cat:\n"
		msg += "```\n"
		msg += note.Content
		msg += "```"
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(msg, ev.Channel))

	default:
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(s.printHelp(), ev.Channel))
	}
	return nil
}

func (s *SlackBot) printHelp() string {
	return fmt.Sprintf(`
%s
- ping: Test the reachability of the bot.
        usage) $BOT ping

- link: Append markdown link to the Kibela note.
        usage) $BOT link $URL [$TITLE]
%s`, "```", "```")
}
