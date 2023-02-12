package binance

import (
	"fmt"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`msg=(.*).$`)

func ParseError(err error) error {
	errStr := err.Error()

	parsedErrors := pattern.FindStringSubmatch(errStr)

	l := len(parsedErrors)
	if l == 0 {
		return err
	}

	e := parsedErrors[len(parsedErrors)-1]
	if len(e) == 0 {
		return fmt.Errorf(e)
	}

	e = strings.ToLower(e[:1]) + e[1:]

	return fmt.Errorf(e)
}
