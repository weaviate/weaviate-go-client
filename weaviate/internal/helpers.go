package internal

import "encoding/json"

// ToPropertiesMapFunc defines the signature for ToPropertiesMap, allowing for mocking.
type ToPropertiesMapFunc func(data any) (map[string]any, error)

// ToPropertiesMap is a variable holding the actual implementation, which can be replaced for testing.
var ToPropertiesMap ToPropertiesMapFunc = func(data any) (map[string]any, error) {
	if m, ok := data.(map[string]any); ok {
		return m, nil
	}

	// Try JSON marshaling/unmarshaling for struct types
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
