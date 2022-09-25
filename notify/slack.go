package notify

import (
	"context"
	"fmt"

	"github.com/MiteshSharma/SlackBot/core/logger"
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

func (s *SlackNotifier) SendMessage(ctx context.Context, message string, targetChannel string) error {
	_, timestamp, err := s.client.PostMessage(targetChannel, slack.MsgOptionText(message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		return fmt.Errorf("error while posting message to channel %q with error %w at %q", targetChannel, err, timestamp)
	}
	return nil
}

func (s *SlackNotifier) SendAttachmentMessage(ctx context.Context, attachment slack.Attachment, targetChannel string) error {
	_, timestamp, err := s.client.PostMessage(targetChannel, slack.MsgOptionAttachments(attachment), slack.MsgOptionAsUser(true))
	if err != nil {
		return fmt.Errorf("error while posting attachmenet message to channel %q with error %w at %q", targetChannel, err, timestamp)
	}
	return nil
}

func (s *SlackNotifier) SendDialog(ctx context.Context, dialog slack.Dialog, triggerId string) error {
	err := s.client.OpenDialogContext(ctx, triggerId, dialog)
	if err != nil {
		return fmt.Errorf("error while ppsting dialog to triggerId %q with error %w", triggerId, err)
	}
	return nil
}

func (s *SlackNotifier) PublishView(userID string, view slack.HomeTabViewRequest, hash string) error {
	_, err := s.client.PublishView(userID, view, hash)
	if err != nil {
		return fmt.Errorf("error while publishing home view with error %w", err)
	}
	return err
}
