package repository

import (
	"context"

	"go-reading-log-api-next/internal/domain/models"
)

// LogRepository defines the interface for log data access operations
type LogRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Log, error)
	GetByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error)
	GetAll(ctx context.Context) ([]*models.Log, error)
	GetByProjectIDOrdered(ctx context.Context, projectID int64) ([]*models.Log, error)
	Create(ctx context.Context, log *models.Log) (*models.Log, error)
}
