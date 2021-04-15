package engine

import (
	"fmt"
	"strings"
)

func ExtractTaskName(eventType string) (string, error) {
	parts := strings.Split(eventType, ".")
	lenParts := len(parts)
	if lenParts != 5 {
		return "", fmt.Errorf("no valid event")
	}
	return parts[3], nil
}
