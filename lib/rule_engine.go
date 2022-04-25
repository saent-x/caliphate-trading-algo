package lib

type RuleEngineResult int32

const (
	Viable    RuleEngineResult = 1
	NotViable RuleEngineResult = 0
)

func ReviewCandleBank() (RuleEngineResult, CandleBias) {
	if _RULE1_containsInvalids() {
		return NotViable, Invalid
	} else {
		bias := _RULE2_bias()

		if bias == Invalid {
			return NotViable, Invalid
		}

		if !_RULE3_invalidHighorLow(bias) {
			return NotViable, Invalid
		}

		return Viable, bias
	}
}

func _RULE1_containsInvalids() bool {
	T, T1, T2 := WithdrawCandleBank()

	return T.Bias == Invalid || T1.Bias == Invalid || T2.Bias == Invalid
}

func _RULE2_bias() CandleBias {
	T, T1, _ := WithdrawCandleBank()

	if T.Bias == T1.Bias {
		return T.Bias
	} else {
		return Invalid
	}
}

func _RULE3_invalidHighorLow(bias CandleBias) bool {
	T, T1, T2 := WithdrawCandleBank()

	switch bias {
	case Bullish:
		if T.Close.GreaterThan(T1.Close) && T.Close.GreaterThan(T2.Close) {
			return true
		} else {
			return false
		}
	case Bearish:
		if T.Close.LessThan(T1.Close) && T.Close.LessThan(T2.Close) {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}
