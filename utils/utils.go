package utils

import (
	"fmt"
	"strconv"
)

func StringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float: %w", err)
	}

	return f, nil
}

func Float64ToString(f float64, precision int) string {
	a := "%." + fmt.Sprintf("%d", precision) + "f"
	return fmt.Sprintf(a, f)
}

func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}
