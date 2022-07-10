package notify

import (
	"context"
)

type Event struct {
	Message      string
	Channel      string
	IsAttachment bool
}

type Notifier interface {
	SendEvent(context.Context, Event) error
}
