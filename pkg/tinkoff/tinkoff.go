package tinkoff

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	sdk "github.com/tinkoff/invest-api-go-sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type TinkoffAPI struct {
	Ctx context.Context

	InstrumentsClient      sdk.InstrumentsServiceClient
	MarketDataClient       sdk.MarketDataServiceClient
	OperationsClient       sdk.OperationsServiceClient
	UsersClient            sdk.UsersServiceClient
	OrdersClient           sdk.OrdersServiceClient
	StopOrdersClient       sdk.StopOrdersServiceClient
	SandboxClient          sdk.SandboxServiceClient
	OrdersStreamClient     sdk.OrdersStreamServiceClient
	OperationsStreamClient sdk.OperationsStreamServiceClient
}

const addressProd = "invest-public-api.tinkoff.ru:443"

const registerTimeout = 30 * time.Second

func New(ctx context.Context, token string) *TinkoffAPI {
	conn, err := grpc.Dial(
		addressProd,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("grpc dial: %v", err)
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	api := &TinkoffAPI{}

	api.InstrumentsClient = sdk.NewInstrumentsServiceClient(conn)
	api.MarketDataClient = sdk.NewMarketDataServiceClient(conn)
	api.OperationsClient = sdk.NewOperationsServiceClient(conn)
	api.UsersClient = sdk.NewUsersServiceClient(conn)
	api.OrdersClient = sdk.NewOrdersServiceClient(conn)
	api.StopOrdersClient = sdk.NewStopOrdersServiceClient(conn)
	api.SandboxClient = sdk.NewSandboxServiceClient(conn)
	api.OrdersStreamClient = sdk.NewOrdersStreamServiceClient(conn)
	api.OperationsStreamClient = sdk.NewOperationsStreamServiceClient(conn)

	return api
}

func NewSandbox(token string, initialBalance float64) *TinkoffAPI {
	return &TinkoffAPI{}
}
