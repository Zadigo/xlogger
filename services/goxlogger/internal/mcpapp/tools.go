package mcpapp

import (
	"context"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetLogsInput struct{}

type GetLogsOutput struct {
	Files []models.File `json:"files"`
}

type McpTools struct{}

func (t *McpTools) GetLogs(ctx context.Context, req *mcp.CallToolRequest, in GetLogsInput) (*mcp.CallToolResult, GetLogsOutput, error) {
	collector := tickerapp.FileCollector{}

	rootDir := ctx.Value("rootDir").(string)
	if files, err := collector.CollectFilesInFolder(rootDir, ""); err != nil {
		return nil, GetLogsOutput{}, err
	} else {
		return nil, GetLogsOutput{Files: files}, nil
	}
}
