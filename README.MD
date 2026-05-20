# Go X-Logger

Go X-Logger is a lightweight and flexible logging library made with Go which allows you to parse large log files and extract meaningful information from them. It supports various log formats and provides an easy-to-use API for logging messages at different levels (info, warning, error, etc.).

It comes with an integrated MCP server based on mcp-golang which allows any AI agent to interrogate the logs.

## Architecture

```mermaid
sequenceDiagram
    participant User
    participant GoXLogger
    participant MCPServer

    User->>GoXLogger: Log message (info, warning, error)
    GoXLogger->>MCPServer: Store log message
    User->>MCPServer: Query logs
    MCPServer->>GoXLogger: Retrieve relevant log messages
    GoXLogger->>MCPServer: Return log messages
    MCPServer->>User: Display log messages
```
