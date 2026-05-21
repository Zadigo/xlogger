package models

import "context"

type LogServiceInterface interface {
	GetLogs(ctx context.Context) ([]string, error)
}
