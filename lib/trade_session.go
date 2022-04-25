package lib

type TradeSession struct {
	RISK_AMOUNT string
	STOPLOSS_DISTANCE string
	SYMBOL string
}


func (tradeSession *TradeSession) CreateSession(risk_amount string, stoploss_distance string, symbol string) *TradeSession {
	_tradeSession := new(TradeSession)
	
	_tradeSession.RISK_AMOUNT = risk_amount
	_tradeSession.STOPLOSS_DISTANCE = stoploss_distance
	_tradeSession.SYMBOL = symbol

	return _tradeSession;
}
