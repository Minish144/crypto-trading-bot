package binance

import (
	"context"
	"errors"

	"github.com/Minish144/crypto-trading-bot/entity"
	gobinance "github.com/adshao/go-binance/v2"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/web"
)

var _ broker.Dealer = (*BinanceDealer)(nil)

var ErrNotImplemented = errors.New("not implemented")

type BinanceDealer struct {
	asset entity.Asset
	*gobinance.Client
}

func NewBinanceDealer(c *Config, symbol string, test bool) *BinanceDealer {
	client := &BinanceDealer{
		asset: entity.Asset{
			Symbol: symbol,
			Base:   c.Base,
		},
		Client: gobinance.NewClient(c.Key, c.Secret),
	}

	if symbol == "" {
		client.asset.Symbol = entity.AssetAnySymbol
	}

	if test {
		gobinance.UseTestnet = true
		client.Client = gobinance.NewClient(c.TestKey, c.TestSecret)
	}

	return client
}

func (d *BinanceDealer) GetBalance(context.Context) (*broker.AccountBalance, *web.Response, error)
func (d *BinanceDealer) PlaceOrder(context.Context, broker.Order) (*broker.Order, *web.Response, error)
func (d *BinanceDealer) CancelOrders(context.Context) (*web.Response, error)
func (d *BinanceDealer) ListPositions(context.Context, *web.ListOpts) ([]broker.Position, *web.Response, error)
func (d *BinanceDealer) ListRoundTurns(context.Context, *web.ListOpts) ([]broker.RoundTurn, *web.Response, error)
