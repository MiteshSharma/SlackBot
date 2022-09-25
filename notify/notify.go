package notify

import (
	"context"

	"github.com/slack-go/slack"
)

type Event struct {
	Message      string
	Channel      string
	IsAttachment bool
}

type Notifier interface {
	SendMessage(ctx context.Context, message string, targetChannel string) error
	SendAttachmentMessage(ctx context.Context, attachment slack.Attachment, targetChannel string) error
	SendDialog(ctx context.Context, dialog slack.Dialog, triggerId string) error
	PublishView(userID string, view slack.HomeTabViewRequest, hash string) error
}
