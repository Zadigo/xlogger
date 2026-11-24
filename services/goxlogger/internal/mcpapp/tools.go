package mcpapp

import (
	"context"

	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpTools struct{}

func (t *McpTools) GetLogs(ctx context.Context, ctr *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collector := tickerapp.FileCollector{}

	rootDir := ctx.Value("rootDir").(string)
	if files, err := collector.CollectFilesInFolder(rootDir, ""); err != nil {
		return nil, err
	} else {
		return &mcp.CallToolResult{StructuredContent: files}, nil
	}
}
