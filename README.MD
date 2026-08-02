# Go X-Logger

Go X-Logger is a lightweight and flexible logging library made with Go which allows you to parse large log files and extract meaningful information from them. It supports various log formats and provides an easy-to-use API for logging messages at different levels (info, warning, error, etc.).

It comes with an integrated MCP server based on mcp-golang which allows any AI agent to interrogate the logs.

## Architecture

```mermaid
sequenceDiagram

autonumber
  
actor U as User
participant F@{type: "entity"} as File
participant T as Ticker
participant C@{type: "entity"} as Claude
participant R@{type: "database"} as Redis
participant M as MCPServer

U->>F: Log message (info, warning, error)

par message analysis
F<<-->>T: Analyze file
T->>R: Store log message
end 

U->>C: Query logs
C->>M: Retrieve relevant log messages
M->>R: Get log messages
R->>C: Return messages
```
