package app

import (
	"context"

	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/MiteshSharma/SlackBot/metrics"
	"github.com/MiteshSharma/SlackBot/notify"
	"github.com/MiteshSharma/SlackBot/pkg/model"
	"github.com/MiteshSharma/SlackBot/pkg/repository"
)

type App struct {
	Context    context.Context
	Config     *config.Config
	Metrics    metrics.Metrics
	Log        logger.Logger
	BotNotify  notify.Notifier
	Repository repository.Repository
}

func NewApp(context context.Context, serverParam *model.ServerParam, repository repository.Repository, botNotify notify.Notifier) *App {
	app := &App{
		Context:    context,
		Config:     serverParam.Config,
		Metrics:    serverParam.Metrics,
		Log:        serverParam.Logger,
		BotNotify:  botNotify,
		Repository: repository,
	}
	return app
}
