package notify

import (
	"context"
	"fmt"

	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/slack-go/slack"
)

// Slack contains Token for authentication with slack and Channel name to send notification to
type SlackNotifier struct {
	log     logger.Logger
	channel string
	client  *slack.Client
}

func NewSlackNotifier(log logger.Logger, token string, channel string) *SlackNotifier {
	return &SlackNotifier{
		log:     log,
		channel: channel,
		client:  slack.New(token),
	}
}

// SendEvent sends event notification to slack
func (s *SlackNotifier) SendEvent(ctx context.Context, event Event) error {
	s.log.Debug(fmt.Sprint("Sending message to slack with title: %w", event.Title))
	attachment := getSlackMessage(event)

	targetChannel := event.Channel
	if targetChannel == "" {
		targetChannel = s.channel
	}

	channelID, timestamp, err := s.client.PostMessage(targetChannel, slack.MsgOptionAttachments(attachment), slack.MsgOptionAsUser(true))
	if err != nil {
		return fmt.Errorf("error while posting message to channel %q with error %w at %q", targetChannel, err, timestamp)
	}
	s.log.Debug(fmt.Sprint("Event successfully sent to channel %w", channelID))
	return nil
}

// Format slack message in proper attachment format
func getSlackMessage(event Event) slack.Attachment {
	return slack.Attachment{
		Pretext: fmt.Sprintf("*%s*", event.Title),
		Fields: []slack.AttachmentField{
			{

				Title: "Message",
				Value: event.Message,
				Short: true,
			},
		},
		Footer: "Message from Bot",
	}
}
