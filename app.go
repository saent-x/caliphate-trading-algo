package main

import (
	// "context"
	"fmt"
	// "log"
	"os"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/rivo/tview"

	"github.com/joho/godotenv"
	"github.com/saent-x/caliphate/lib"
	// "github.com/saent-x/caliphate/utilities"
	// "github.com/shopspring/decimal"
)

type Slide func(nextSlide func()) (title string, content tview.Primitive)

func main() {
	// load env file
	godotenv.Load()

	var STOP_LOSS_DISTANCE_FOR_SESSION string = "0.0025"
	var RISK_AMOUNT string = "10"
	var ASSET string = "ETHUSDT"

	futures.UseTestnet = true

	_ = binance.NewClient(os.Getenv("BINANCE-CLIENT-API-KEY"), os.Getenv("BINANCE-CLIENT-SECRET"))
	futuresClient := binance.NewFuturesClient(os.Getenv("BINANCE-FUTURES-CLIENT-API-KEY"), os.Getenv("BINANCE-FUTURES-CLIENT-SECRET"))

	caliphate := lib.TradeSession{}
	caliphate_session := caliphate.CreateSession(RISK_AMOUNT, STOP_LOSS_DISTANCE_FOR_SESSION, ASSET)

	fmt.Println("** caliphate in session **")
	wait_flag0 := make(chan bool)
	wait_flag1 := make(chan bool)

	go lib.OrderManager(futuresClient, caliphate_session, wait_flag0)

	// lib.CreateOrder(lib.Viable, lib.Bearish, futuresClient, caliphate_session)
	// lib.CreateOrder(lib.Viable, lib.Bearish, futuresClient, caliphate_session)

	go lib.WatchCrypto(futuresClient, caliphate_session, wait_flag1)

	<-wait_flag0
	<-wait_flag1

	fmt.Println("** caliphate closed safely **")
}
