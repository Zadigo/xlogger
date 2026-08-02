package mcpapp

import (
	"github.com/Zadigo/goxlogger/internal/mcpapp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// func (a *McpApp) loadPrompts() {
// 	handler := &prompts.LogsPrompt{}
// 	a.server.AddPrompt(&mcp.Prompt{}, func(ctx context.Context, req *mcp.CallPromptRequest) (*mcp.CallPromptResult, error) {
// 		return handler.GetLogs(ctx, req)
// 	})
// }

// func (a *McpApp) loadResources() {
// 	handler := &prompts.LogsPrompt{}
// 	a.server.AddResource(&mcp.Resource{}, func(ctx context.Context, req *mcp.CallPromptRequest) (*mcp.CallPromptResult, error) {
// 		return handler.GetLogs(ctx, req)
// 	})
// }

func (a *McpApp) loadTools() {
	handler := &tools.BaseTools{}

	a.server.AddTool(&mcp.Tool{}, handler.GetLogs)
}
