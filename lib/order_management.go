package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

type OrderIDS struct {
	MarketOrderID     int64
	StoplossOrderID   int64
	TakeprofitOrderID int64
}
type Order struct {
	OrderID       OrderIDS
	Symbol        string
	Count         int
	Entry         string
	Stoploss      string
	Takeprofit    string
	EntryTime     string
	Profit_n_Loss string
}

var Orders []Order

const MULTIPLIER = 3

func CreateOrder(result RuleEngineResult, bias CandleBias, client *futures.Client, session *TradeSession) {
	if result == NotViable {
		fmt.Println("-- no viable trade")
		return
	}

	if bias == Bullish {
		goLong(client, session)
	} else if bias == Bearish {
		goShort(client, session)
	}
}

func OrderManager(client *futures.Client, session *TradeSession, flag chan bool) {
	listenKey, err := client.NewStartUserStreamService().Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	OnUserData := func(event *futures.WsUserDataEvent) {
		client.NewSetServerTimeService().Do(context.Background())
		if event.Event == futures.UserDataEventTypeListenKeyExpired {
			// restart it
			OrderManager(client, session, flag)
		}

		if event.OrderTradeUpdate.Symbol == session.SYMBOL {
			if event.OrderTradeUpdate.Type == futures.OrderTypeMarket {
				// check if the order status is filled
				if event.OrderTradeUpdate.Status == futures.OrderStatusTypeFilled {
					// check if it matches any order tp/sl
					for _, order := range Orders {
						if order.OrderID.StoplossOrderID == event.OrderTradeUpdate.ID {
							// remove tp order
							_order, err := client.NewGetOrderService().
								Symbol(session.SYMBOL).
								OrderID(order.OrderID.TakeprofitOrderID).
								Do(context.Background())

							if err != nil {
								log.Fatal(err)
							}

							if _order.Status == futures.OrderStatusTypeCanceled || _order.Status == futures.OrderStatusTypeFilled {
								return
							}

							_, err = client.NewCancelOrderService().
								OrderID(order.OrderID.TakeprofitOrderID).
								Symbol(session.SYMBOL).
								Do(context.Background())
							if err != nil {
								log.Fatal(err)
							}

							fmt.Printf("[info] %v takeprofit order cancelled\n", event.OrderTradeUpdate.Symbol)

						} else if order.OrderID.TakeprofitOrderID == event.OrderTradeUpdate.ID {
							// remove sl order
							_order, err := client.NewGetOrderService().
								Symbol(session.SYMBOL).
								OrderID(order.OrderID.StoplossOrderID).
								Do(context.Background())

							if err != nil {
								log.Fatal(err)
							}

							if _order.Status == futures.OrderStatusTypeCanceled || _order.Status == futures.OrderStatusTypeFilled {
								return
							}

							_, err = client.NewCancelOrderService().
								OrderID(order.OrderID.StoplossOrderID).
								Symbol(session.SYMBOL).
								Do(context.Background())
							if err != nil {
								log.Fatal(err)
							}

							fmt.Printf("[info] %v stoploss order cancelled\n", event.OrderTradeUpdate.Symbol)
						} else {
							fmt.Printf("[info] %v market order initiated | status: %v\n", event.OrderTradeUpdate.Symbol, event.OrderTradeUpdate.Status)
						}
					}
				}
			} else if event.OrderTradeUpdate.Type == futures.OrderTypeStopMarket {
				// is it new or what
				if event.OrderTradeUpdate.Status == futures.OrderStatusTypeNew {
					fmt.Printf("[info] %v stoploss order initiated\n", event.OrderTradeUpdate.Symbol)

				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypePartiallyFilled {
					fmt.Println("[info] partially filled!!")
				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypeCanceled {
					// cancel tp and market order
					for _, order := range Orders {
						if order.OrderID.StoplossOrderID == event.OrderTradeUpdate.ID {

							_order, err := client.NewGetOrderService().Symbol(session.SYMBOL).OrderID(order.OrderID.TakeprofitOrderID).Do(context.Background())

							if err != nil {
								log.Fatal(err)
							}

							if _order.Status == futures.OrderStatusTypeCanceled || _order.Status == futures.OrderStatusTypeFilled {
								return
							}
							_, err = client.NewCancelOrderService().
								OrderID(order.OrderID.TakeprofitOrderID).
								Symbol(session.SYMBOL).
								Do(context.Background())
							if err != nil {
								log.Fatal(err)
							}

							fmt.Printf("[info] %v stoploss order cancelled\n", event.OrderTradeUpdate.Symbol)
						}
					}

				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypeExpired {
					fmt.Printf("[info] %v stoploss order expired\n", event.OrderTradeUpdate.Symbol)

				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypeRejected {
					fmt.Printf("[info] %v stoploss order rejected\n", event.OrderTradeUpdate.Symbol)

				}
			} else if event.OrderTradeUpdate.Type == futures.OrderTypeTakeProfitMarket {
				// is it new or what
				if event.OrderTradeUpdate.Status == futures.OrderStatusTypeNew {
					fmt.Printf("[info] %v takeprofit order initiated\n", event.OrderTradeUpdate.Symbol)

				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypePartiallyFilled {
					fmt.Println("[info] partially filled!!")
				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypeCanceled {
					// cancel sl and market order
					for _, order := range Orders {
						if order.OrderID.TakeprofitOrderID == event.OrderTradeUpdate.ID {

							_order, err := client.NewGetOrderService().
								Symbol(session.SYMBOL).
								OrderID(order.OrderID.StoplossOrderID).
								Do(context.Background())

							if err != nil {
								log.Fatal(err)
							}

							if _order.Status == futures.OrderStatusTypeCanceled || _order.Status == futures.OrderStatusTypeFilled {
								return
							}

							_, err = client.NewCancelOrderService().
								OrderID(order.OrderID.StoplossOrderID).
								Symbol(session.SYMBOL).
								Do(context.Background())
							if err != nil {
								log.Fatal(err)
							}

							fmt.Printf("[info] %v takeprofit order cancelled\n", event.OrderTradeUpdate.Symbol)
						}
					}

				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypeExpired {
					fmt.Printf("[info] %v takeprofit order expired\n", event.OrderTradeUpdate.Symbol)

				} else if event.OrderTradeUpdate.Status == futures.OrderStatusTypeRejected {
					fmt.Printf("[info] %v takeprofit order rejected\n", event.OrderTradeUpdate.Symbol)

				}
			}
		}
	}
	OnError := func(err error) {
		log.Fatal(err)
		flag <- true
	}

	doneC, _, err := futures.WsUserDataServe(listenKey, OnUserData, OnError)

	if err != nil {
		log.Fatal(err)
		flag <- true
	}

	<-doneC
}

func goLong(client *futures.Client, session *TradeSession) {
	// use the client future for Futures
	client.NewSetServerTimeService().Do(context.Background())

	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var quantity_precision decimal.Decimal
	var price_precision decimal.Decimal

	for _, _symbol := range info.Symbols {
		if _symbol.Symbol == session.SYMBOL {
			quantity_precision = decimal.NewFromInt32(int32(_symbol.QuantityPrecision))
			price_precision = decimal.NewFromInt32(int32(_symbol.PricePrecision))
		}
	}

	// variables
	var risk_amount decimal.Decimal
	var stoploss_distance decimal.Decimal

	stoploss_distance, _ = decimal.NewFromString(session.STOPLOSS_DISTANCE)

	// get amount to be risked
	risk_amount, _ = decimal.NewFromString(session.RISK_AMOUNT)

	// calculate position size / quantity
	position_size := risk_amount.Abs().DivRound(stoploss_distance, int32(quantity_precision.IntPart()))

	// convert position size to base currency value
	var base_currency_position_size decimal.Decimal

	// get current symbol price
	ticker_list, err := client.NewListPricesService().Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for _, ticker := range ticker_list {
		if ticker.Symbol == session.SYMBOL {
			_ticker, _ := decimal.NewFromString(ticker.Price)
			base_currency_position_size = position_size.DivRound(_ticker, int32(quantity_precision.IntPart()))
		}
	}

	// we can now use the value of the position size in the base currency to enter a market order

	order, err := client.NewCreateOrderService().Symbol(session.SYMBOL).
		Side(futures.SideTypeBuy).
		Type(futures.OrderTypeMarket).
		Quantity(base_currency_position_size.String()).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		WorkingType(futures.WorkingTypeMarkPrice).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	// once order is in then activate stoploss and take profit | determine the entry price

	entry_price, _ := decimal.NewFromString(order.AvgPrice)
	rr_buffer := entry_price.Mul(stoploss_distance)

	stoploss_price := entry_price.Sub(rr_buffer).Round(int32(price_precision.IntPart()))
	takeprofit_price := entry_price.Add(rr_buffer.Mul(decimal.NewFromInt32(MULTIPLIER))).Round(int32(price_precision.IntPart()))

	// create stoploss and takeprofit orders
	stoploss_order, err := client.NewCreateOrderService().Symbol(session.SYMBOL).
		Side(futures.SideTypeSell).Type(futures.OrderTypeStopMarket).
		StopPrice(stoploss_price.String()).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Quantity(base_currency_position_size.String()).
		TimeInForce(futures.TimeInForceTypeGTC).
		ReduceOnly(true).
		WorkingType(futures.WorkingTypeMarkPrice).
		Do(context.Background())

	if err != nil {
		log.Fatal(err.Error())
	}

	takeprofit_order, err := client.NewCreateOrderService().Symbol(session.SYMBOL).
		Side(futures.SideTypeSell).Type(futures.OrderTypeTakeProfitMarket).
		StopPrice(takeprofit_price.String()).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Quantity(base_currency_position_size.String()).
		TimeInForce(futures.TimeInForceTypeGTC).
		ReduceOnly(true).
		WorkingType(futures.WorkingTypeMarkPrice).
		Do(context.Background())

	if err != nil {
		log.Fatal(err.Error())
	}
	new_order := Order{
		OrderID: OrderIDS{
			MarketOrderID:     order.OrderID,
			StoplossOrderID:   stoploss_order.OrderID,
			TakeprofitOrderID: takeprofit_order.OrderID,
		},
		Symbol:        session.SYMBOL,
		Count:         len(Orders) + 1,
		Entry:         entry_price.String(),
		Stoploss:      stoploss_price.String(),
		Takeprofit:    takeprofit_price.String(),
		EntryTime:     "***",
		Profit_n_Loss: "0.0%",
	}
	Orders = append(Orders, new_order)

	fmt.Printf("[info] new order - symbol: %v entry: %v stoploss: %v takeprofit: %v \n", new_order.Symbol, new_order.Entry, new_order.Stoploss, new_order.Takeprofit)
}

func goShort(client *futures.Client, session *TradeSession) {
	// use the client future for Futures
	client.NewSetServerTimeService().Do(context.Background())

	info, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var quantity_precision decimal.Decimal
	var price_precision decimal.Decimal

	for _, _symbol := range info.Symbols {
		if _symbol.Symbol == session.SYMBOL {
			quantity_precision = decimal.NewFromInt32(int32(_symbol.QuantityPrecision))
			price_precision = decimal.NewFromInt32(int32(_symbol.PricePrecision))
		}
	}

	// variables
	var risk_amount decimal.Decimal
	var stoploss_distance decimal.Decimal

	stoploss_distance, _ = decimal.NewFromString(session.STOPLOSS_DISTANCE)

	// get amount to be risked
	risk_amount, _ = decimal.NewFromString(session.RISK_AMOUNT)

	// calculate position size / quantity
	position_size := risk_amount.Abs().DivRound(stoploss_distance, int32(quantity_precision.IntPart()))

	// convert position size to base currency value
	var base_currency_position_size decimal.Decimal

	// get current symbol price
	ticker_list, err := client.NewListPricesService().Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for _, ticker := range ticker_list {
		if ticker.Symbol == session.SYMBOL {
			_ticker, _ := decimal.NewFromString(ticker.Price)
			base_currency_position_size = position_size.DivRound(_ticker, int32(quantity_precision.IntPart()))
		}
	}

	// we can now use the value of the position size in the base currency to enter a market order

	order, err := client.NewCreateOrderService().Symbol(session.SYMBOL).
		Side(futures.SideTypeSell).
		Type(futures.OrderTypeMarket).
		Quantity(base_currency_position_size.String()).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		WorkingType(futures.WorkingTypeMarkPrice).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	// once order is in then activate stoploss and take profit | determine the entry price

	entry_price, _ := decimal.NewFromString(order.AvgPrice)
	rr_buffer := entry_price.Mul(stoploss_distance)

	stoploss_price := entry_price.Add(rr_buffer).Round(int32(price_precision.IntPart()))
	takeprofit_price := entry_price.Sub(rr_buffer.Mul(decimal.NewFromInt32(MULTIPLIER))).Round(int32(price_precision.IntPart()))

	// create stoploss and takeprofit orders
	stoploss_order, err := client.NewCreateOrderService().Symbol(session.SYMBOL).
		Side(futures.SideTypeBuy).Type(futures.OrderTypeStopMarket).
		StopPrice(stoploss_price.String()).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Quantity(base_currency_position_size.String()).
		TimeInForce(futures.TimeInForceTypeGTC).
		ReduceOnly(true).
		WorkingType(futures.WorkingTypeMarkPrice).
		Do(context.Background())

	if err != nil {
		log.Fatal(err.Error())
	}

	takeprofit_order, err := client.NewCreateOrderService().Symbol(session.SYMBOL).
		Side(futures.SideTypeBuy).Type(futures.OrderTypeTakeProfitMarket).
		StopPrice(takeprofit_price.String()).
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Quantity(base_currency_position_size.String()).
		TimeInForce(futures.TimeInForceTypeGTC).
		ReduceOnly(true).
		WorkingType(futures.WorkingTypeMarkPrice).
		Do(context.Background())

	if err != nil {
		log.Fatal(err.Error())
	}

	new_order := Order{
		OrderID: OrderIDS{
			MarketOrderID:     order.OrderID,
			StoplossOrderID:   stoploss_order.OrderID,
			TakeprofitOrderID: takeprofit_order.OrderID,
		},
		Symbol:        session.SYMBOL,
		Count:         len(Orders) + 1,
		Entry:         entry_price.String(),
		Stoploss:      stoploss_price.String(),
		Takeprofit:    takeprofit_price.String(),
		EntryTime:     "***",
		Profit_n_Loss: "0.0%",
	}
	Orders = append(Orders, new_order)

	fmt.Printf("[info] new order - symbol: %v entry: %v stoploss: %v takeprofit: %v \n", new_order.Symbol, new_order.Entry, new_order.Stoploss, new_order.Takeprofit)
}
