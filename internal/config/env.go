package config

import (
	"os"
	"strings"
)

// getCleanedEnv returns an environment variable value with
// surrounding spaces and trailing new line characters removed.
func getCleanedEnv(envKey string) (value string) {
	value = os.Getenv(envKey)
	value = strings.TrimSpace(value)
	value = strings.TrimSuffix(value, "\r\n")
	value = strings.TrimSuffix(value, "\n")
	return value
}
