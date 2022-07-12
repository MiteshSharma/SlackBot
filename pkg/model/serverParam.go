package model

import (
	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/core/logger"
	"github.com/MiteshSharma/SlackBot/core/metrics"
)

type ServerParam struct {
	Logger  logger.Logger
	Metrics metrics.Metrics
	Config  *config.Config
}

func NewServerParam(logger logger.Logger, metrics metrics.Metrics, config *config.Config) *ServerParam {
	err := &ServerParam{
		Logger:  logger,
		Metrics: metrics,
		Config:  config,
	}
	return err
}
