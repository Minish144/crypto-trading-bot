package utils

import (
	"fmt"
	"strings"
)

func TransforSlashseparatedSymbol(s string) (string, error) {
	coins, err := GetCoinsFromSymbol(s)
	if err != nil {
		return "", err
	}

	return strings.Join(coins, ""), nil
}

func GetCoinsFromSymbol(s string) ([]string, error) {
	coins := strings.Split(s, "/")
	if len(coins) != 2 {
		return nil, fmt.Errorf("failed to split symbol by slash: got not 2 coins")
	}

	return coins, nil
}

func ConvertSymbol(s string) (string, error) {
	coins, err := GetCoinsFromSymbol(s)
	if err != nil {
		return "", err
	}

	return strings.Join(coins, ""), nil
}

func GetBaseCoin(s string) (string, error) {
	coins, err := GetCoinsFromSymbol(s)
	if err != nil {
		return "", err
	}

	return coins[1], nil
}

func GetQuoteCoin(s string) (string, error) {
	coins, err := GetCoinsFromSymbol(s)
	if err != nil {
		return "", err
	}

	return coins[0], nil
}
