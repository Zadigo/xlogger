package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpTools struct{}

func (t *McpTools) GetLogs(ctx context.Context, ctr *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{}, nil
}
