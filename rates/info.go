package rates

import "solution/ent/schema"

type ExchangeInfo struct {
	data []exchangeInfo
}

type exchangeInfo struct {
	From schema.Currency
	To   schema.Currency
	Rate float64
}

func (c ExchangeInfo) GetRate(from, to schema.Currency) float64 {
	for _, info := range c.data {
		if info.From == from && info.To == to {
			return info.Rate
		}

		if info.From == to && info.To == from {
			return 1 / info.Rate
		}
	}

	return 0
}
