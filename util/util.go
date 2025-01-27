package util

import (
	"os"
)

func GetFullUrl(path string) string {

	if os.Getenv("ENVIRONMENT") == "DEV" {
		return path
	}
	// return fmt.Sprintf("/etracker%s", path)
	return path
}

func StringSliceContains(slice []string, needle string) bool {
	for _, elem := range slice {
		if elem == needle {
			return true
		}
	}
	return false
}
