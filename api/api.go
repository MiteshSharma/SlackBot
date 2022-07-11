package api

import (
	"context"

	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/MiteshSharma/SlackBot/notify"
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
}

func NewServerAPI(appContext context.Context, router *mux.Router, serverParam *ServerParam, botNotify notify.Notifier) *ServerAPI {
	api := &ServerAPI{
		Context:    appContext,
		MainRouter: router,
		Config:     serverParam.Config,
		Metrics:    serverParam.Metrics,
		Log:        serverParam.Logger,
		Router:     &Router{},
		BotNotify:  botNotify,
	}
	api.setupRoutes()
	return api
}
