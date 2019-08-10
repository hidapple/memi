package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
	"golang.org/x/xerrors"
)

// SlackBot respresents Slack Bot App
type SlackBot struct {
	client   *slack.Client
	rtm      *slack.RTM
	id       string
	help     string
	commands []*Command
}

// NewBot creates new SlackBot
func NewBot(token, help string) *SlackBot {
	return &SlackBot{
		client: slack.New(token),
		help:   help,
	}
}

// Command defines bot command
type Command struct {
	Name string
	Do   func(args ...string) (string, error)
}

// AddCommands implements commands in SlackBot
func (s *SlackBot) AddCommands(commands ...*Command) {
	s.commands = append(s.commands, commands...)
}

// OpenRTM opens RTM connection and listens incomming events
func (s *SlackBot) OpenRTM(ctx context.Context) error {
	s.rtm = s.client.NewRTM()

	go s.rtm.ManageConnection()
	for {
		select {
		case <-ctx.Done():
			s.rtm.Disconnect()
			return ctx.Err()
		case msg, ok := <-s.rtm.IncomingEvents:
			if !ok {
				return nil
			}
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				s.id = ev.Info.User.ID

			case *slack.MessageEvent:
				if err := s.handleMessage(ev); err != nil {
					log.Printf("failed to handle message: %s", err)
				}

			case *slack.RTMError:
				return xerrors.Errorf("RTM Error. code=%d, msg=%s", ev.Code, ev.Msg)
			}
		}
	}
}

func (s *SlackBot) handleMessage(ev *slack.MessageEvent) error {
	// Input is supposed to be `$bot_name $command [argument(s)]`
	msg := strings.Split(ev.Msg.Text, " ")

	// Ignore if input is not bot mention
	if !strings.HasPrefix(msg[0], fmt.Sprintf("<@%s>", s.id)) {
		return nil
	}

	// Print help when command is not given
	if len(msg) < 2 {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(s.help, ev.Channel))
		return nil
	}
	return s.exec(ev, msg[1], msg[2:]...)
}

func (s *SlackBot) exec(ev *slack.MessageEvent, cmdName string, args ...string) error {
	cmd := s.searchCommand(cmdName)
	if cmd == nil {
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(s.help, ev.Channel))
		return nil
	}
	msg, err := cmd.Do(args...)
	if err != nil {
		return err
	}
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(msg, ev.Channel))
	return nil
}

func (s *SlackBot) searchCommand(cmdName string) *Command {
	for _, cmd := range s.commands {
		if cmd.Name == cmdName {
			return cmd
		}
	}
	return nil
}
