package sqlRepository

import (
	"net/http"

	"github.com/MiteshSharma/SlackBot/core/sql"
	"github.com/MiteshSharma/SlackBot/pkg/model"
)

type WorkspaceRepository struct {
	*sql.SqlRepository
}

func NewWorkspaceRepository(sqlRepository *sql.SqlRepository) WorkspaceRepository {
	workspaceRepository := WorkspaceRepository{sqlRepository}

	if !workspaceRepository.DB.HasTable(&model.Workspace{}) {
		workspaceRepository.DB.CreateTable(&model.Workspace{})
	}
	return workspaceRepository
}

// CreateWorkspace func is used to create workspace object in db
func (wr WorkspaceRepository) CreateWorkspace(workspace *model.Workspace) *model.StorageResult {
	if err := wr.DB.Create(&workspace).Error; err != nil {
		return model.NewStorageResult(nil, model.NewAppError(err.Error(), http.StatusInternalServerError))
	}
	return model.NewStorageResult(workspace, nil)
}

func (wr WorkspaceRepository) UpdateWorkspace(workspace *model.Workspace) *model.StorageResult {
	if err := wr.DB.Save(&workspace).Error; err != nil {
		return model.NewStorageResult(nil, model.NewAppError(err.Error(), http.StatusInternalServerError))
	}
	return model.NewStorageResult(workspace, nil)
}

func (wr WorkspaceRepository) GetWorkspace(workspaceID string) *model.StorageResult {
	var workspace model.Workspace
	if err := wr.DB.First(&workspace, "workspace_id = ?", workspaceID).Error; err != nil {
		return model.NewStorageResult(nil, model.NewAppError(err.Error(), http.StatusInternalServerError))
	}
	return model.NewStorageResult(&workspace, nil)
}
