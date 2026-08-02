package mcpapp

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpApp struct {
	ctx    context.Context
	server *mcp.Server
}

func (a *McpApp) Start() {
	// Use the official SDK's native HTTP transport struct directly
	httpTransport := &mcp.StreamableServerTransport{}

	go func() {
		if err := a.server.Run(a.ctx, httpTransport); err != nil {
			log.Fatal(err)
		}
	}()
}

func NewMcpApp(ctx context.Context) *McpApp {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "my-mcp-server",
			Version: "1.0.0",
		},
		nil,
	)

	app := &McpApp{
		ctx:    ctx,
		server: server,
	}

	app.loadTools()
	return app
}
