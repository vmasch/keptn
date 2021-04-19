package utils

import (
	"fmt"
	"strings"
)

func ExtractStageName(eventType string) (string, error) {
	parts := strings.Split(eventType, ".")
	lenParts := len(parts)
	if lenParts != 6 {
		return "", fmt.Errorf("no valid event")
	}
	return parts[3], nil
}

func ExtractSequenceName(eventType string) (string, error) {
	parts := strings.Split(eventType, ".")
	lenParts := len(parts)
	if lenParts != 6 {
		return "", fmt.Errorf("no valid event")
	}
	return parts[4], nil
}

func ExtractTaskName(eventType string) (string, error) {
	parts := strings.Split(eventType, ".")
	lenParts := len(parts)
	if lenParts != 5 {
		return "", fmt.Errorf("no valid event")
	}
	return parts[3], nil
}
