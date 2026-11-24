package mcpapp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (a *McpApp) loadTools() {
	handler := &McpTools{}
	a.server.AddTool(&mcp.Tool{Description: "Get the logs present in the folder"}, handler.GetLogs)
}
