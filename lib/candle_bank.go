package lib

import (
	"github.com/shopspring/decimal"
)

type CandleBias int32

const (
	Bullish CandleBias = 1
	Bearish CandleBias = -1
	Invalid CandleBias = 0
)

type Candle struct {
	Open  decimal.Decimal
	Close decimal.Decimal
	High  decimal.Decimal
	Low   decimal.Decimal
	Time  int64
	Bias  CandleBias
}

var _candleBasket []Candle

func (candle *Candle) SetBias() {

	if candle.Close.GreaterThan(candle.Open) {
		candle.Bias = Bullish
	} else if candle.Open.GreaterThan(candle.Close) {
		candle.Bias = Bearish
	} else {
		candle.Bias = Invalid
	}
}

// func doji(candle *Candle) bool {
// 	openCloseDifference := candle.Open.Sub(candle.Close).Abs()
// 	multiplier, _ := decimal.NewFromString("0.1")
// 	benchmark := candle.High.Abs().Sub(candle.Low).Mul(multiplier)

// 	return openCloseDifference.LessThanOrEqual(benchmark)
// }

func FillCandleBank(candle Candle) {
	// check current count of the stack, if not 3 the add new candle directly
	if len(_candleBasket) < 3 {
		_candleBasket = append(_candleBasket, candle)
	} else if len(_candleBasket) == 3 {
		// remove the last item and add to stack

		_candleBasket = append(_candleBasket[:0], _candleBasket[1:]...)
		_candleBasket = append(_candleBasket, candle)
	}
}

func ValidateCandleBank() bool {
	return len(_candleBasket) == 3
}

func WithdrawCandleBank() (Candle, Candle, Candle) {
	return _candleBasket[2], _candleBasket[1], _candleBasket[0]
}
