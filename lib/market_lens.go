package lib

import (
	"log"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

func WatchCrypto(client *futures.Client, session *TradeSession, flag chan bool) {

	OnCandleUpdate := func(event *futures.WsKlineEvent) {
		if event.Kline.IsFinal {

			open, _ := decimal.NewFromString(event.Kline.Open)
			close, _ := decimal.NewFromString(event.Kline.Close)
			high, _ := decimal.NewFromString(event.Kline.High)
			low, _ := decimal.NewFromString(event.Kline.Low)

			candle := Candle{
				Open:  open,
				Close: close,
				High:  high,
				Low:   low,
				Time:  event.Time,
			}
			candle.SetBias()

			FillCandleBank(candle)

			if ValidateCandleBank() {
				reviewResult, bias := ReviewCandleBank()
				CreateOrder(reviewResult, bias, client, session)
			}
		}
	}
	OnError := func(err error) {
		log.Fatal(err)
		flag <- true
	}
	doneC, _, err := futures.WsKlineServe(session.SYMBOL, "30m", OnCandleUpdate, OnError)

	if err != nil {
		log.Fatal(err)
		flag <- true
	}

	<-doneC

}

func WatchForex(){
	
}
