package prompts

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type LogsPrompt struct{}

func (p *LogsPrompt) GetLogs(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{}, nil
}
