package cursor

import (
	"time"

	"github.com/minish144/crypto-trading-bot/domain"
)

type Cursor interface{}

type cursor struct {
	Ts    time.Time
	Kline *domain.Kline
	next  *cursor
}

func (c *cursor) Next() bool {
	if c.next == nil {
		return false
	}

	*c = *c.next

	return true
}

func NewFromKlines(klines []*domain.Kline, start time.Time, interval time.Duration) *cursor {
	rootCur := &cursor{
		next:  nil,
		Ts:    start,
		Kline: nil,
	}

	var c = rootCur

	for i, kline := range klines {
		currentCur := cursor{
			Ts:    start.Add(time.Duration(i+1) * interval),
			Kline: kline,
		}

		c.next = &currentCur
		c = c.next
	}

	return rootCur
}
