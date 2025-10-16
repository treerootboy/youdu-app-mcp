# Youdu Multi-Interface Service

A comprehensive service for Youdu IM that provides CLI, MCP (Model Context Protocol), and API interfaces through a unified adapter layer.

## Architecture

```
   CLI       MCP       API (planned)
    │         │         │
    └─────────┼─────────┘
              │
           Adapter
              │
           Youdu-SDK
```

## Features

- **Unified Adapter Layer**: All Youdu SDK operations are wrapped in a simplified adapter layer
- **Automatic Tool/Command Generation**: CLI commands and MCP tools are automatically generated from adapter methods using reflection
- **Type-Safe**: Full type safety with Go structs and JSON schema annotations
- **Configuration Management**: Flexible configuration via files and environment variables

## Installation

### Prerequisites

- Go 1.23 or higher
- Access to a Youdu IM server

### Build

```bash
# Build MCP server
go build -o bin/youdu-mcp ./cmd/youdu-mcp

# Build CLI
go build -o bin/youdu-cli ./cmd/youdu-cli
```

## Configuration

Create a `config.yaml` file in your project root or at `~/.youdu/config.yaml`:

```yaml
youdu:
  addr: "http://your-youdu-server:7080"
  buin: 123456789
  app_id: "your-app-id"
  aes_key: "your-aes-key"
```

Or use environment variables:

```bash
export YOUDU_ADDR="http://your-youdu-server:7080"
export YOUDU_BUIN=123456789
export YOUDU_APP_ID="your-app-id"
export YOUDU_AES_KEY="your-aes-key"
```

## Usage

### CLI

The CLI provides commands organized by functionality:

```bash
# List all commands
./bin/youdu-cli --help

# Department operations
./bin/youdu-cli dept get-list --dept-id=0
./bin/youdu-cli dept get-user-list --dept-id=1
./bin/youdu-cli dept create --name="Engineering" --parent-id=0

# User operations
./bin/youdu-cli user get --user-id="user123"
./bin/youdu-cli user create --user-id="newuser" --name="New User" --dept-id=1

# Message operations
./bin/youdu-cli message send-text-message --to-user="user123" --content="Hello!"

# Group operations
./bin/youdu-cli group get-list --user-id="user123"
./bin/youdu-cli group create --name="Project Team"

# Session operations
./bin/youdu-cli session create --title="Team Chat" --creator="user123" --type="group"
```

### MCP Server

The MCP server provides all adapter methods as MCP tools that can be called by Claude Desktop or other MCP clients.

#### Running the MCP Server

```bash
./bin/youdu-mcp
```

#### Claude Desktop Integration

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "youdu": {
      "command": "/path/to/youdu-app-mcp/bin/youdu-mcp"
    }
  }
}
```

#### Available MCP Tools

All tools follow snake_case naming:

- **Department**: `get_dept_list`, `get_dept_user_list`, `get_dept_alias_list`, `create_dept`, `update_dept`, `delete_dept`
- **User**: `get_user`, `create_user`, `update_user`, `delete_user`
- **Message**: `send_text_message`, `send_image_message`, `send_file_message`, `send_link_message`, `send_sys_message`
- **Group**: `get_group_list`, `get_group_info`, `create_group`, `update_group`, `delete_group`, `add_group_member`, `del_group_member`
- **Session**: `create_session`, `get_session`, `update_session`, `send_text_session_message`, `send_image_session_message`, `send_file_session_message`

## Project Structure

```
youdu-app-mcp/
├── cmd/
│   ├── youdu-cli/          # CLI entry point
│   └── youdu-mcp/          # MCP server entry point
├── internal/
│   ├── adapter/            # Adapter layer (core business logic)
│   │   ├── adapter.go      # Base adapter
│   │   ├── dept.go         # Department methods
│   │   ├── user.go         # User methods
│   │   ├── message.go      # Message methods
│   │   ├── group.go        # Group methods
│   │   └── session.go      # Session methods
│   ├── cli/                # CLI implementation
│   │   ├── root.go         # Root command
│   │   └── generator.go    # Auto-generate commands
│   ├── mcp/                # MCP server implementation
│   │   └── server.go       # Auto-register tools
│   └── config/             # Configuration management
│       └── config.go       # Viper configuration
├── bin/                    # Built binaries
├── config.yaml.example     # Example configuration
├── go.mod
├── go.sum
└── README.md
```

## Development

### Adding New Methods

To add a new Youdu API method:

1. Add the method to the appropriate adapter file (e.g., `internal/adapter/dept.go`)
2. Follow the pattern:
   ```go
   type MethodNameInput struct {
       Field string `json:"field" jsonschema:"description=Field description,required"`
   }

   type MethodNameOutput struct {
       Result string `json:"result" jsonschema:"description=Result description"`
   }

   func (a *Adapter) MethodName(ctx context.Context, input MethodNameInput) (*MethodNameOutput, error) {
       // Implementation
   }
   ```
3. The method will automatically be available as:
   - CLI command: `youdu-cli category method-name --field=value`
   - MCP tool: `method_name`

### Key Design Principles

1. **Single Source of Truth**: All APIs are defined once in the adapter layer
2. **Auto-Generation**: CLI commands and MCP tools are automatically generated using reflection
3. **Type Safety**: Input/Output structs with JSON schema annotations
4. **Simplicity**: Adapter methods have simple, intuitive names and parameters

## Dependencies

- [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official MCP SDK
- [github.com/addcnos/youdu/v2](https://github.com/addcnos/youdu) - Youdu IM SDK
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [github.com/spf13/viper](https://github.com/spf13/viper) - Configuration management

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
