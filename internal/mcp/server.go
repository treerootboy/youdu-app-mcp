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
	// Create input schema from the input type
	inputSchema := generateInputSchema(inputType)

	// Create tool definition
	tool := &mcp.Tool{
		Name:        name,
		Description: description,
		InputSchema: inputSchema,
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

// generateInputSchema generates a JSON schema for the input type
func generateInputSchema(inputType reflect.Type) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}

	properties := schema["properties"].(map[string]interface{})
	required := []string{}

	// Iterate through struct fields
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Extract field name from JSON tag
		fieldName := strings.Split(jsonTag, ",")[0]

		// Get jsonschema tag for description and required
		schemaTag := field.Tag.Get("jsonschema")
		fieldSchema := map[string]interface{}{
			"type": getJSONType(field.Type),
		}

		// Parse jsonschema tag
		if schemaTag != "" {
			parts := strings.Split(schemaTag, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "description=") {
					fieldSchema["description"] = strings.TrimPrefix(part, "description=")
				} else if part == "required" {
					required = append(required, fieldName)
				} else if strings.HasPrefix(part, "default=") {
					// Skip default values - MCP SDK has strict type checking
					// Users can provide defaults in their requests
					continue
				}
			}
		}

		properties[fieldName] = fieldSchema
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// getJSONType returns the JSON schema type for a Go type
func getJSONType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "string"
	}
}
