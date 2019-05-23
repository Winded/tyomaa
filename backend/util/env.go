package util

import "os"

func EnvOrDefault(name, defaultValue string) string {
	value, found := os.LookupEnv(name)
	if found {
		return value
	} else {
		return defaultValue
	}
}
