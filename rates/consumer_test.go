package rates

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"solution/ent/schema"
	"testing"
)

func Test(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	i, ok := toInfo(logger, "GBPEUR", 1.05)
	assert.True(t, ok)

	assert.Equal(t, exchangeInfo{
		From: schema.Currency("GBP"),
		To:   schema.Currency("EUR"),
		Rate: 1.05,
	}, i)
}

func Test2(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	m := &Manager{
		reader: nil,
		logger: logger,
	}

	m.convertRate(map[string]float64{
		"EURUSD": 1.08,
		"GBPUSD": 1.22,
		"USDRUB": 79.14,
		"GBPEUR": 1.14,
		"GBPRUB": 98.15,
		"EURRUB": 89.10,
	})

	assert.NotNil(t, m.Rate())
}
