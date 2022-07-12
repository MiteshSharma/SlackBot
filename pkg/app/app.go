package app

import (
	"context"

	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/MiteshSharma/SlackBot/metrics"
	"github.com/MiteshSharma/SlackBot/notify"
	"github.com/MiteshSharma/SlackBot/pkg/model"
)

type App struct {
	Context   context.Context
	Config    *config.Config
	Metrics   metrics.Metrics
	Log       logger.Logger
	BotNotify notify.Notifier
}

func NewApp(context context.Context, serverParam *model.ServerParam, botNotify notify.Notifier) *App {
	app := &App{
		Context:   context,
		Config:    serverParam.Config,
		Metrics:   serverParam.Metrics,
		Log:       serverParam.Logger,
		BotNotify: botNotify,
	}
	return app
}
