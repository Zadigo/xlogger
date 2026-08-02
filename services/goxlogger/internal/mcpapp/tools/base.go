package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type BaseTools struct{}

func (b *BaseTools) GetLogs(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{}, nil
}
