package types

import "fmt"

type simpleThreshold struct {
	low  float64
	high float64
}

func NewThreshold(low, high float64) *simpleThreshold {
	return &simpleThreshold{
		low:  low,
		high: high,
	}
}

func (st simpleThreshold) LowerBound() float64 {
	return st.low
}

func (st simpleThreshold) UpperBound() float64 {
	return st.high
}

func (st *simpleThreshold) UpdateBounds(low, high float64) {
	st.low, st.high = low, high
}

func (st simpleThreshold) ToString() string {
	return fmt.Sprintf("%+v", st)
}
