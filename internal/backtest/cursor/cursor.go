package cursor

import (
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
)

type Cursor interface{}

type cursor struct {
	Kline *domain.Kline
	Ts    time.Time
	next  *cursor
}

func (c *cursor) Next() bool {
	if c.next == nil {
		return false
	}

	*c = *c.next

	return true
}

func NewFromKlines(klines []*domain.Kline) *cursor {
	if len(klines) == 0 {
		return nil
	}

	rootCur := &cursor{
		Kline: klines[0],
		Ts:    klines[0].Ts,
		next:  nil,
	}

	var c = rootCur

	for _, kline := range klines[1:] {
		currentCur := cursor{
			Ts:    kline.Ts,
			Kline: kline,
		}

		c.next = &currentCur
		c = c.next
	}

	return rootCur
}
