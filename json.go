package fix

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func jsonMarshal(target interface{}) ([]byte, error) {
	return json.MarshalIndent(target, "", "  ")
}

func jsonCompare(old, new []byte) error {
	// If equal, target still have the same result.
	if bytes.Equal(new, old) {
		return nil
	}

	return fmt.Errorf("diff: %s", lineDiff(string(old), string(new)))
}