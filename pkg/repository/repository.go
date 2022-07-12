package repository

import "github.com/MiteshSharma/SlackBot/pkg/model"

type Repository interface {
	Close() error
	Workspace() WorkspaceRepository
}

type WorkspaceRepository interface {
	CreateWorkspace(workspace *model.Workspace) *model.StorageResult
	UpdateWorkspace(workspace *model.Workspace) *model.StorageResult
	GetWorkspace(workspaceID string) *model.StorageResult
}
