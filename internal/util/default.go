package util

import "fmt"

func FormPath(method, route, basePath string) string {
	return fmt.Sprintf("%s %s%s", method, basePath, route)
}
