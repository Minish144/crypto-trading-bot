package tinkoff

import (
	"context"
	"log"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

type TinkoffAPI struct {
	Client *sdk.RestClient
}

const registerTimeout = 30 * time.Second

func New(token string) *TinkoffAPI {
	return &TinkoffAPI{Client: sdk.NewRestClient(token)}
}

func NewSandbox(token string, initialBalance float64) *TinkoffAPI {
	client := sdk.NewSandboxRestClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), registerTimeout)
	defer cancel()

	account, err := client.Register(ctx, sdk.AccountTinkoff)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.SetCurrencyBalance(ctx, account.ID, sdk.RUB, initialBalance)
	if err != nil {
		log.Fatalln(err)
	}

	return &TinkoffAPI{Client: client.RestClient}
}
