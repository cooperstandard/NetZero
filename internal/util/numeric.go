package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//TODO: add and subtract, convert from numeric string to dollars and Cents

type Numeric struct {
	Dollars int
	Cents   int
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
	parts := strings.Split(s, ".")
	dollars, err := strconv.Atoi(parts[0])
	if err != nil {
		return Numeric{}, errors.New("invalid dollar portion of the numeric string")
	}

	cents, err := strconv.Atoi(parts[1])

	if err != nil {
		return Numeric{}, errors.New("invalid cents portion of the numeric string")
	}

	return Numeric{Dollars: dollars, Cents: cents}, nil
}

func NumericAddition(a, b Numeric) Numeric {
	return Numeric{}
}

func NumericSubtraction(a, b Numeric) Numeric {
	return Numeric{}
}
