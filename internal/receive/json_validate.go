package receive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/jsonschema-go/jsonschema"
)

func ValidateJSON(jsonSchemaPath string, jsonData string) error {
	schemaBytes, err := os.ReadFile(jsonSchemaPath)
	if err != nil {
		return fmt.Errorf("Error reading JSON Schema file: %w", err)
	}

	var schema jsonschema.Schema
	if err := json.Unmarshal(schemaBytes, &schema); err != nil {
		return fmt.Errorf("error parsing JSON Schema: %w", err)
	}

	resolved, err := schema.Resolve(&jsonschema.ResolveOptions{})
	if err != nil {
		return fmt.Errorf("Error resolving JSON Schema: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return fmt.Errorf("Error parsing JSON data: %w", err)
	}

	return resolved.Validate(data)
}
