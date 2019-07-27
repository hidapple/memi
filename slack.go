package main

import (
	"fmt"
	"log"
	"strings"

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
				log.Printf("Failed to handle message: %s", err)
			}

		case *slack.RTMError:
			log.Printf("RTM Error. code=%d, msg=%s", ev.Code, ev.Msg)
		}
	}
}

func (s *SlackBot) handleMessage(ev *slack.MessageEvent) error {
	msg := strings.Split(ev.Msg.Text, " ")
	if !strings.HasPrefix(msg[0], fmt.Sprintf("<@%s>", s.botID)) {
		return nil
	}

	switch msg[1] {
	case "ping":
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("にゃん:two_hearts:", ev.Channel))

	default:
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(s.printHelp(), ev.Channel))
	}
	return nil
}

// TODO: Implement help message
func (s *SlackBot) printHelp() string {
	return `Help message...`
}
