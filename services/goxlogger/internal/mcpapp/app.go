package mcpapp

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpApp struct {
	ctx    context.Context
	server *mcp.Server
}

func (a *McpApp) ChiHandle() *mcp.StreamableHTTPHandler {
	a.server = mcp.NewServer(
		&mcp.Implementation{
			Name:       "X-Logger",
			Title:      "X Logger",
			WebsiteURL: "",
			Version:    "1.0.0",
		},
		&mcp.ServerOptions{},
	)

	a.loadTools()

	return mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return a.server
	}, &mcp.StreamableHTTPOptions{})
}

// func (a *McpApp) Start(ctx context.Context) error {
// 	httpTransport := &mcp.StreamableServerTransport{}

// 	chErr := make(chan error, 1)

// 	go func() {
// 		slog.Info("⚡️ Starting MCP server")
// 		chErr <- a.server.Run(a.ctx, httpTransport)
// 	}()

// 	select {
// 	case err := <-chErr:
// 		return err
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	}
// }

// func NewMcpApp(ctx context.Context) *McpApp {
// 	server := mcp.NewServer(
// 		&mcp.Implementation{
// 			Name:       "X-Logger",
// 			Title:      "X Logger",
// 			WebsiteURL: "",
// 			Version:    "1.0.0",
// 		},
// 		&mcp.ServerOptions{},
// 	)

// 	app := &McpApp{
// 		ctx:    ctx,
// 		server: server,
// 	}

// 	app.loadTools()
// 	return app
// }
