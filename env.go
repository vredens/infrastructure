package infrastructure

import "os"

// GetFromEnv returns the first value found in the environment variables.
func GetFromEnv(keys ...string) string {
	for i := range keys {
		if value, ok := os.LookupEnv(keys[i]); ok {
			return value
		}
	}
	return ""
}
