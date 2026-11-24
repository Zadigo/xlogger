package mcpapp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (a *McpApp) loadTools() {
	handler := &McpTools{}
	tool := &mcp.Tool{
		Title:       "Get Logs",
		Name:        "get_logs",
		Description: "Get the logs present in the folder",
	}
	mcp.AddTool(a.server, tool, handler.GetLogs)
}
