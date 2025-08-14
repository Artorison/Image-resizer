package helpers

import "fmt"

func GenerateCacheKey(width, height int, url string) string {
	return fmt.Sprintf("%dx%d|%s", width, height, url)
}
