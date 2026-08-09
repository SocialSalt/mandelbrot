package utils

import "math/big"

var FP128_PRECISION uint = 113

func FP128Quotient(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).SetPrec(FP128_PRECISION).Quo(x, y)
}

func FP128Product(x *big.Float, y *big.Float) *big.Float {
	return new(big.Float).SetPrec(FP128_PRECISION).Mul(x, y)
}

func NewBigFloat(val float64) *big.Float {
	return new(big.Float).SetPrec(FP128_PRECISION).SetFloat64(val)
}
