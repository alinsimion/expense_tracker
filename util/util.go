package util

import (
	"fmt"
	"os"
)

func GetFullUrl(path string) string {

	if os.Getenv("ENVIRONMENT") == "DEV" {
		return path
	}
	return fmt.Sprintf("/etracker%s", path)
}
