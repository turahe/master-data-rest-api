package search

import (
	"encoding/json"
	"fmt"
)

// mapToStruct converts a map[string]interface{} to a struct
func mapToStruct(data interface{}, target interface{}) error {
	// Convert to JSON bytes first
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Unmarshal into target struct
	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
