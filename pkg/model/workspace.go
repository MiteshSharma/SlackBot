package model

import (
	"encoding/json"
	"errors"
	"time"
)

// Workspace struct
type Workspace struct {
	WorkspaceID string     `gorm:"primary_key" json:"workspaceId"`
	OwnerID     string     `gorm:"type:varchar(64)" json:"ownerId"`
	Name        string     `gorm:"type:varchar(64)" json:"name"`
	AccessToken string     `json:"accessToken"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
}

// Valid function is to check if policy object is valid
func (w *Workspace) Valid() error {
	if w.WorkspaceID == "" {
		return errors.New("workspace id can not be null or empty")
	}
	if w.AccessToken == "" {
		return errors.New("workspace access token can not be null or empty")
	}
	return nil
}

func (w *Workspace) ToJson() string {
	json, _ := json.Marshal(w)
	return string(json)
}
