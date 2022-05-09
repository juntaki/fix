package fix

import (
	"bytes"
	"fmt"
	"github.com/juntaki/pp"
)

func ppMarshal(target interface{}) ([]byte, error) {
	pp.ColoringEnabled = false
	return []byte(pp.Sprintln(target)), nil
}

func ppCompare(old, new []byte) error {
	// If equal, target still have the same result.
	if bytes.Equal(new, old) {
		return nil
	}

	return fmt.Errorf("diff: %s", lineDiff(string(old), string(new)))
}
