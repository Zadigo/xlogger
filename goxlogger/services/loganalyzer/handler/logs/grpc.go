package logs

import (
	"context"

	genprotologs "github.com/Zadigo/goxlogger/services/common/genproto/logs"
	"github.com/Zadigo/goxlogger/services/loganalyzer/models"
)

type LogsGrpcHandler struct {
	logsService models.LogServiceInterface
	genprotologs.UnimplementedLogsServiceServer
}

func (LogsGrpcHandler) GetLogs(ctx context.Context, req *genprotologs.GetLogsRequest) (*genprotologs.GetLogsResponse, error) {
	// Call the service layer to get logs
	return &genprotologs.GetLogsResponse{}, nil
}

func NewLogsGrpcHandler(logsService models.LogServiceInterface) *LogsGrpcHandler {
	return &LogsGrpcHandler{
		logsService: logsService,
	}
}
