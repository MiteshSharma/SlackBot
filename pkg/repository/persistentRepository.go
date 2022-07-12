package repository

import (
	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/metrics"
	"github.com/MiteshSharma/SlackBot/pkg/repository/sqlRepository"

	"github.com/MiteshSharma/SlackBot/core/sql"
	"github.com/MiteshSharma/SlackBot/logger"
)

type PersistentRepository struct {
	SqlRepository *sql.SqlRepository
	Log           logger.Logger
	Config        *config.Config
	Metrics       metrics.Metrics

	WorkspaceRepository WorkspaceRepository
}

func NewPersistentRepository(log logger.Logger, config *config.Config, metrics metrics.Metrics) PersistentRepository {
	repository := PersistentRepository{
		Log:     log,
		Config:  config,
		Metrics: metrics,
	}

	repository.SqlRepository = sql.NewSqlRepository(log, config.DatabaseConfig)
	repository.WorkspaceRepository = sqlRepository.NewWorkspaceRepository(repository.SqlRepository)

	return repository
}

func (s PersistentRepository) Workspace() WorkspaceRepository {
	return s.WorkspaceRepository
}

func (s PersistentRepository) Close() error {
	return s.SqlRepository.Close()
}
