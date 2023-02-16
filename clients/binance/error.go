package binance

import (
	"fmt"
	"regexp"
)

var pattern = regexp.MustCompile(`msg=(.*)$`)

func ParseError(err error) error {
	errStr := err.Error()

	parsedErrors := pattern.FindStringSubmatch(errStr)

	l := len(parsedErrors)
	if l == 0 {
		return err
	}

	return fmt.Errorf(parsedErrors[l-1])
}
