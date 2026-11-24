package mcpapp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpApp struct {
	ctx    context.Context
	server *mcp.Server
}

func (a *McpApp) Start(ctx context.Context) error {
	httpTransport := &mcp.StreamableServerTransport{}

	chErr := make(chan error, 1)
	
	go func() {
		chErr <- a.server.Run(a.ctx, httpTransport)
	}()

	select {
	case err := <-chErr:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
