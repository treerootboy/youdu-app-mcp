package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/yourusername/youdu-app-mcp/internal/adapter"
)

// generateCommands uses reflection to generate CLI commands from adapter methods
func generateCommands() error {
	// Create a temporary adapter instance to inspect its methods
	adapterType := reflect.TypeOf((*adapter.Adapter)(nil))

	// Create command groups
	deptCmd := &cobra.Command{
		Use:   "dept",
		Short: "Department management commands",
	}

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	messageCmd := &cobra.Command{
		Use:   "message",
		Short: "Message sending commands",
	}

	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Group management commands",
	}

	sessionCmd := &cobra.Command{
		Use:   "session",
		Short: "Session management commands",
	}

	// Iterate through all methods
	for i := 0; i < adapterType.NumMethod(); i++ {
		method := adapterType.Method(i)

		// Skip non-exported methods and special methods
		if !method.IsExported() || method.Name == "Close" || method.Name == "Context" {
			continue
		}

		// Get method type information
		methodType := method.Type

		// Verify method signature
		if methodType.NumIn() != 3 || methodType.NumOut() != 2 {
			continue
		}

		// Get input and output types
		inputType := methodType.In(2)
		outputType := methodType.Out(0)

		// Determine which command group this belongs to
		var parentCmd *cobra.Command
		methodNameLower := strings.ToLower(method.Name)
		if strings.Contains(methodNameLower, "dept") {
			parentCmd = deptCmd
		} else if strings.Contains(methodNameLower, "user") {
			parentCmd = userCmd
		} else if strings.Contains(methodNameLower, "message") {
			parentCmd = messageCmd
		} else if strings.Contains(methodNameLower, "group") {
			parentCmd = groupCmd
		} else if strings.Contains(methodNameLower, "session") {
			parentCmd = sessionCmd
		} else {
			// Add to root for uncategorized commands
			parentCmd = rootCmd
		}

		// Generate command
		cmd := generateCommand(method.Name, inputType, outputType)
		parentCmd.AddCommand(cmd)
	}

	// Add command groups to root
	rootCmd.AddCommand(deptCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(messageCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(sessionCmd)

	return nil
}

// generateCommand creates a cobra command for a specific adapter method
func generateCommand(methodName string, inputType, outputType reflect.Type) *cobra.Command {
	// Convert method name to kebab-case for command name
	cmdName := toKebabCase(methodName)

	// Generate description
	description := generateDescription(methodName)

	// Create flags map to store input values
	inputValues := make(map[string]interface{})

	cmd := &cobra.Command{
		Use:   cmdName,
		Short: description,
		Long:  fmt.Sprintf("%s\n\nThis command calls the %s method.", description, methodName),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create input struct
			input := reflect.New(inputType).Interface()

			// Populate input from flags
			if err := populateInputFromFlags(cmd, inputValues, input); err != nil {
				return fmt.Errorf("failed to populate input: %w", err)
			}

			// Call the adapter method
			adapterValue := reflect.ValueOf(youduAdapter)
			method := adapterValue.MethodByName(methodName)

			if !method.IsValid() {
				return fmt.Errorf("method %s not found", methodName)
			}

			// Call method
			results := method.Call([]reflect.Value{
				reflect.ValueOf(context.Background()),
				reflect.ValueOf(input).Elem(),
			})

			// Check for error
			if !results[1].IsNil() {
				err := results[1].Interface().(error)
				return err
			}

			// Print output as JSON
			output := results[0].Interface()
			outputJSON, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal output: %w", err)
			}

			fmt.Println(string(outputJSON))
			return nil
		},
	}

	// Add flags for each input field
	addFlagsForStruct(cmd, inputType, inputValues)

	return cmd
}

// addFlagsForStruct adds flags for each field in a struct
func addFlagsForStruct(cmd *cobra.Command, structType reflect.Type, inputValues map[string]interface{}) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Get JSON tag for flag name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse JSON tag (handle omitempty, etc.)
		flagName := strings.Split(jsonTag, ",")[0]

		// Get jsonschema tag for description
		description := field.Tag.Get("jsonschema")
		if description == "" {
			description = fmt.Sprintf("%s field", field.Name)
		} else {
			// Extract description from jsonschema tag
			if idx := strings.Index(description, "description="); idx != -1 {
				desc := description[idx+12:]
				if endIdx := strings.Index(desc, ","); endIdx != -1 {
					desc = desc[:endIdx]
				}
				description = desc
			}
		}

		// Add flag based on field type
		switch field.Type.Kind() {
		case reflect.String:
			var val string
			cmd.Flags().StringVar(&val, flagName, "", description)
			inputValues[flagName] = &val
		case reflect.Int:
			var val int
			cmd.Flags().IntVar(&val, flagName, 0, description)
			inputValues[flagName] = &val
		case reflect.Bool:
			var val bool
			cmd.Flags().BoolVar(&val, flagName, false, description)
			inputValues[flagName] = &val
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				var val []string
				cmd.Flags().StringSliceVar(&val, flagName, []string{}, description)
				inputValues[flagName] = &val
			}
		}
	}
}

// populateInputFromFlags populates the input struct from flag values
func populateInputFromFlags(cmd *cobra.Command, inputValues map[string]interface{}, input interface{}) error {
	inputVal := reflect.ValueOf(input).Elem()
	inputType := inputVal.Type()

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		flagName := strings.Split(jsonTag, ",")[0]

		if valuePtr, ok := inputValues[flagName]; ok {
			fieldVal := inputVal.Field(i)
			if fieldVal.CanSet() {
				valueReflect := reflect.ValueOf(valuePtr).Elem()
				fieldVal.Set(valueReflect)
			}
		}
	}

	return nil
}

// toKebabCase converts a string from PascalCase to kebab-case
func toKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

// generateDescription generates a human-readable description from a method name
func generateDescription(methodName string) string {
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
