package api

import (
	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/gorilla/mux"

	"github.com/MiteshSharma/SlackBot/metrics"
)

type ServerAPI struct {
	MainRouter *mux.Router
	Config     *config.Config
	Metrics    metrics.Metrics
	Log        logger.Logger
	Router     *Router
}

func NewServerAPI(router *mux.Router, serverParam *ServerParam) *ServerAPI {
	api := &ServerAPI{
		MainRouter: router,
		Config:     serverParam.Config,
		Metrics:    serverParam.Metrics,
		Log:        serverParam.Logger,
		Router:     &Router{},
	}
	api.setupRoutes()
	return api
}
