package repository

import (
	"context"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
)

// ProjectWithLogs represents a project with its associated logs
type ProjectWithLogs struct {
	Project *dto.ProjectResponse
	Logs    []*dto.LogResponse
}

// ProjectRepository defines the interface for project data access operations
type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Project, error)
	GetAll(ctx context.Context) ([]*models.Project, error)
	GetWithLogs(ctx context.Context, id int64) (*ProjectWithLogs, error)
	GetAllWithLogs(ctx context.Context) ([]*ProjectWithLogs, error)
	Create(ctx context.Context, project *models.Project) (*models.Project, error)
}
