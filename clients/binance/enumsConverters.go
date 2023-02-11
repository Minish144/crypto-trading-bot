package binance

import (
	"github.com/Minish144/crypto-trading-bot/models"
	gobinance "github.com/adshao/go-binance/v2"
)

var (
	otToModels = map[gobinance.OrderType]models.OrderType{
		gobinance.OrderTypeLimit:           models.OrderTypeLimit,
		gobinance.OrderTypeMarket:          models.OrderTypeMarket,
		gobinance.OrderTypeLimitMaker:      models.OrderTypeLimitMaker,
		gobinance.OrderTypeStopLoss:        models.OrderTypeStopLoss,
		gobinance.OrderTypeStopLossLimit:   models.OrderTypeStopLossLimit,
		gobinance.OrderTypeTakeProfit:      models.OrderTypeTakeProfit,
		gobinance.OrderTypeTakeProfitLimit: models.OrderTypeTakeProfit,
	}

	otFromModels = map[models.OrderType]gobinance.OrderType{
		models.OrderTypeLimit:           gobinance.OrderTypeLimit,
		models.OrderTypeMarket:          gobinance.OrderTypeMarket,
		models.OrderTypeLimitMaker:      gobinance.OrderTypeLimitMaker,
		models.OrderTypeStopLoss:        gobinance.OrderTypeStopLoss,
		models.OrderTypeStopLossLimit:   gobinance.OrderTypeStopLossLimit,
		models.OrderTypeTakeProfit:      gobinance.OrderTypeTakeProfit,
		models.OrderTypeTakeProfitLimit: gobinance.OrderTypeTakeProfitLimit,
	}
)

var (
	tiftToModels = map[gobinance.TimeInForceType]models.TimeInForceType{
		gobinance.TimeInForceTypeGTC: models.TimeInForceTypeGTC,
		gobinance.TimeInForceTypeIOC: models.TimeInForceTypeIOC,
		gobinance.TimeInForceTypeFOK: models.TimeInForceTypeFOK,
	}

	tiftFromModels = map[models.TimeInForceType]gobinance.TimeInForceType{
		models.TimeInForceTypeGTC: gobinance.TimeInForceTypeGTC,
		models.TimeInForceTypeIOC: gobinance.TimeInForceTypeIOC,
		models.TimeInForceTypeFOK: gobinance.TimeInForceTypeFOK,
	}
)

var (
	ostToModels = map[gobinance.OrderStatusType]models.OrderStatusType{
		gobinance.OrderStatusTypeNew:             models.OrderStatusTypeNew,
		gobinance.OrderStatusTypePartiallyFilled: models.OrderStatusTypePartiallyFilled,
		gobinance.OrderStatusTypeFilled:          models.OrderStatusTypeFilled,
		gobinance.OrderStatusTypeCanceled:        models.OrderStatusTypeCanceled,
		gobinance.OrderStatusTypePendingCancel:   models.OrderStatusTypePendingCancel,
		gobinance.OrderStatusTypeRejected:        models.OrderStatusTypeRejected,
		gobinance.OrderStatusTypeExpired:         models.OrderStatusTypeExpired,
	}

	ostFromModels = map[models.OrderStatusType]gobinance.OrderStatusType{
		models.OrderStatusTypeNew:             gobinance.OrderStatusTypeNew,
		models.OrderStatusTypePartiallyFilled: gobinance.OrderStatusTypePartiallyFilled,
		models.OrderStatusTypeFilled:          gobinance.OrderStatusTypeFilled,
		models.OrderStatusTypeCanceled:        gobinance.OrderStatusTypeCanceled,
		models.OrderStatusTypePendingCancel:   gobinance.OrderStatusTypePendingCancel,
		models.OrderStatusTypeRejected:        gobinance.OrderStatusTypeRejected,
		models.OrderStatusTypeExpired:         gobinance.OrderStatusTypeExpired,
	}
)

var (
	stToModels = map[gobinance.SideType]models.SideType{
		gobinance.SideTypeBuy:  models.SideTypeBuy,
		gobinance.SideTypeSell: models.SideTypeSell,
	}

	stFromModels = map[models.SideType]gobinance.SideType{
		models.SideTypeBuy:  gobinance.SideTypeBuy,
		models.SideTypeSell: gobinance.SideTypeSell,
	}
)

var (
	intervalToModels = map[gobinance.RateLimitInterval]models.Interval{
		gobinance.RateLimitIntervalSecond: models.IntervalSecond,
		gobinance.RateLimitIntervalMinute: models.IntervalMinute,
		gobinance.RateLimitIntervalDay:    models.IntervalDay,
	}

	intervalFromModels = map[models.Interval]gobinance.RateLimitInterval{
		models.IntervalSecond: gobinance.RateLimitIntervalSecond,
		models.IntervalMinute: gobinance.RateLimitIntervalMinute,
		models.IntervalDay:    gobinance.RateLimitIntervalDay,
	}
)
