package notify

import (
	"context"
)

type Event struct {
	Title   string
	Message string
	Channel string
}

type Notifier interface {
	SendEvent(context.Context, Event) error
}
