package util

import "fmt"

//TODO: add and subtract, convert from numeric string to dollars and Cents

type Numeric struct {
	Dollars int64
	Cents int
}

func ValidateAndFormNumericString(val Numeric) (string, error) {
	if val.Cents < 0 || val.Cents > 99 {
		return "", fmt.Errorf("%d is not a valid value for cents, cents should be positive and less than 100", val.Cents)
	}
	if val.Dollars < 0 {

		return "", fmt.Errorf("%d is not a valid value, dollars should be positive", val.Dollars)
	}

	return fmt.Sprintf("%20d.%d", val.Dollars, val.Cents), nil
}

func StringToNumeric(s string) (Numeric, error) {


	return Numeric{}, nil
}


