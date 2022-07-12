package api

import (
	"context"

	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/MiteshSharma/SlackBot/notify"
	"github.com/MiteshSharma/SlackBot/pkg/app"
	"github.com/MiteshSharma/SlackBot/pkg/model"
	"github.com/MiteshSharma/SlackBot/pkg/repository"
	"github.com/gorilla/mux"

	"github.com/MiteshSharma/SlackBot/metrics"
)

type ServerAPI struct {
	Context    context.Context
	MainRouter *mux.Router
	Config     *config.Config
	Metrics    metrics.Metrics
	Log        logger.Logger
	Router     *Router
	BotNotify  notify.Notifier
	App        *app.App
}

func NewServerAPI(appContext context.Context, router *mux.Router, repository repository.Repository,
	serverParam *model.ServerParam, botNotify notify.Notifier) *ServerAPI {

	api := &ServerAPI{
		Context:    appContext,
		MainRouter: router,
		Config:     serverParam.Config,
		Metrics:    serverParam.Metrics,
		Log:        serverParam.Logger,
		Router:     &Router{},
		BotNotify:  botNotify,
		App:        app.NewApp(appContext, serverParam, repository, botNotify),
	}

	api.setupRoutes()
	return api
}
