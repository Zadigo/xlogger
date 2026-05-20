package mcp

// type McpApp struct {
// 	ctx    context.Context
// 	server *mcp.Server
// }

// func (a *McpApp) Start() {
// 	// Create an HTTP transport
// 	httpTransport := http.NewHTTPTransport("/mcp")
// 	httpTransport.WithAddr(":8080")

// 	go func() {
// 		if err := a.server.Run(a.ctx, httpTransport); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// }

// func NewMcpApp(ctx context.Context) *McpApp {
// 	app := &McpApp{ctx: ctx}
// 	app.loadRoutes()
// 	return app
// }
