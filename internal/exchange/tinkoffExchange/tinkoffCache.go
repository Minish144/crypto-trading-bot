package tinkoffExchange

type TinkoffCache struct {
	instruments map[string]*Instrument // ticker - instrument
}

func NewCache() *TinkoffCache {
	return &TinkoffCache{instruments: make(map[string]*Instrument)}
}

func (c *TinkoffCache) AddInstrument(i *Instrument) {
	c.instruments[i.Ticker] = i
}

func (c *TinkoffCache) AddInstruments(is []*Instrument) {
	for _, i := range is {
		c.AddInstrument(i)
	}
}
