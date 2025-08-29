package util

import "fmt"

func FormPath(method, route, basePath string) string {
	return fmt.Sprintf("%s %s%s", method, basePath, route)
}

func ValidateAndFormNumericString(dollars, cents int) (string, error) {
	if cents < 0 || cents > 99 {
		return "", fmt.Errorf("%d is not a valid value for cents, cents should be positive and less than 100", cents)
	}
	if dollars < 0 {

		return "", fmt.Errorf("%d is not a valid value, dollars should be positive", dollars)
	}

	return fmt.Sprintf("%20d.%d", dollars, cents), nil
}
