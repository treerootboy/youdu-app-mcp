package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yourusername/youdu-app-mcp/internal/adapter"
	"github.com/yourusername/youdu-app-mcp/internal/config"
)

// Server represents the MCP server
type Server struct {
	server  *mcp.Server
	adapter *adapter.Adapter
}

// New creates a new MCP server
func New(cfg *config.Config) (*Server, error) {
	adapter, err := adapter.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "youdu-mcp",
		Version: "1.0.0",
	}, nil)

	s := &Server{
		server:  server,
		adapter: adapter,
	}

	// Register all adapter methods as MCP tools
	if err := s.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return s, nil
}

// Run starts the MCP server
func (s *Server) Run(ctx context.Context) error {
	transport := &mcp.StdioTransport{}
	return s.server.Run(ctx, transport)
}

// registerTools uses reflection to register all adapter methods as MCP tools
func (s *Server) registerTools() error {
	adapterType := reflect.TypeOf(s.adapter)
	adapterValue := reflect.ValueOf(s.adapter)

	// Iterate through all methods of the adapter
	for i := 0; i < adapterType.NumMethod(); i++ {
		method := adapterType.Method(i)

		// Skip non-exported methods and special methods
		if !method.IsExported() || method.Name == "Close" || method.Name == "Context" {
			continue
		}

		// Convert method name to snake_case for MCP tool name
		toolName := toSnakeCase(method.Name)

		// Get method type information
		methodType := method.Type

		// Verify method signature: func(ctx context.Context, input InputType) (*OutputType, error)
		if methodType.NumIn() != 3 || methodType.NumOut() != 2 {
			continue
		}

		// Check first parameter is context.Context
		if !methodType.In(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			continue
		}

		// Check last return type is error
		if !methodType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			continue
		}

		// Get input and output types
		inputType := methodType.In(2)
		outputType := methodType.Out(0)

		// Create tool description from method name
		description := generateDescription(method.Name)

		// Register the tool
		if err := s.registerTool(toolName, description, method, adapterValue, inputType, outputType); err != nil {
			return fmt.Errorf("failed to register tool %s: %w", toolName, err)
		}
	}

	return nil
}

// registerTool registers a single tool with the MCP server
func (s *Server) registerTool(name, description string, method reflect.Method, adapterValue reflect.Value, inputType, outputType reflect.Type) error {
	// Create tool definition
	tool := &mcp.Tool{
		Name:        name,
		Description: description,
	}

	// Create handler function
	handler := func(ctx context.Context, req *mcp.CallToolRequest, rawInput json.RawMessage) (*mcp.CallToolResult, any, error) {
		// Create new input instance
		input := reflect.New(inputType).Interface()

		// Unmarshal input
		if err := json.Unmarshal(rawInput, input); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal input: %w", err)
		}

		// Call the adapter method
		results := method.Func.Call([]reflect.Value{
			adapterValue,
			reflect.ValueOf(ctx),
			reflect.ValueOf(input).Elem(),
		})

		// Check for error
		if !results[1].IsNil() {
			err := results[1].Interface().(error)
			return nil, nil, err
		}

		// Return output
		output := results[0].Interface()
		return &mcp.CallToolResult{}, output, nil
	}

	// Add tool to server
	mcp.AddTool(s.server, tool, handler)

	return nil
}

// toSnakeCase converts a string from PascalCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// generateDescription generates a human-readable description from a method name
func generateDescription(methodName string) string {
	// Convert PascalCase to space-separated words
	var words []string
	var currentWord strings.Builder

	for i, r := range methodName {
		if i > 0 && unicode.IsUpper(r) {
			words = append(words, currentWord.String())
			currentWord.Reset()
		}
		currentWord.WriteRune(unicode.ToLower(r))
	}
	words = append(words, currentWord.String())

	return strings.Join(words, " ")
}
