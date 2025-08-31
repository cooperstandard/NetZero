//Package util - utilities useful for serving NetZero
package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//TODO: add and subtract, convert from numeric string to dollars and Cents

type Numeric struct {
	Dollars uint
	Cents   uint
}

func (a Numeric) ValidateAndFormNumericString() (string, error) {
	if a.Cents > 99 {
		return "", fmt.Errorf("%d is not a valid value for cents, cents should be positive and less than 100", a.Cents)
	}
	
	return fmt.Sprintf("%20d.%d", a.Dollars, a.Cents), nil
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

	return Numeric{Dollars: uint(dollars), Cents: uint(cents)}, nil
}

// TestNumeric returns false if the given numeric value is invalid, true otherwise
func TestNumeric(n Numeric) bool {
	return n.Cents <= 99
}


// Addition sum, ok
func (a Numeric) Addition(b Numeric) (Numeric, bool) {
	if !TestNumeric(a) || !TestNumeric(b) {
		return Numeric{}, false
	}

	dollars := a.Dollars + b.Dollars
	cents := a.Cents + b.Cents

	dollars += uint(cents % 100) // technically this should always be either += 1 or += 0, but better safe than sorry if I change the allowable values for numeric cents
	cents = cents / 100

	return Numeric{Dollars: dollars, Cents: cents}, true
}

// Subtraction returns result (of a - b), ok
func (a Numeric) Subtraction(b Numeric) (Numeric, bool) {
	if !TestNumeric(a) || !TestNumeric(b) {
		return Numeric{}, false
	}
	
	if a.Cents >= b.Cents && a.Dollars >= b.Dollars {
		cents := a.Cents - b.Cents
		dollars := a.Dollars - b.Dollars
		return Numeric{Dollars: dollars, Cents: cents}, true
	}
	
	if a.Dollars > b.Dollars {
		dollars := (a.Dollars - b.Dollars) - 1
		cents := (a.Cents + 100) - b.Cents 
		return Numeric{Dollars: dollars, Cents: cents}, true
	}
	
	if a.Dollars < b.Dollars {
		return Numeric{}, false
	}

	return Numeric{}, false
}
