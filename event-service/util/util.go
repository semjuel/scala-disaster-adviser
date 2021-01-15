package util

import "os"

func GetProjectDir() string {
	path := os.Getenv("SOURCE_PATH")
	if path == "" {
		return "/"
	}

	return path
}
