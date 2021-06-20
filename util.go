package indc

import (
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

var (
	// _hundred is 100 in decimal format.
	_hundred = decimal.NewFromInt(100)

	// _one is 1 in decimal format.
	_one = decimal.NewFromInt(1)
)

var (
	// ErrInvalidIndicator is returned when indicator is invalid.
	ErrInvalidIndicator = errors.New("invalid indicator")

	// ErrInvalidLength is returned when incorrect length is provided.
	ErrInvalidLength = errors.New("invalid length")

	// ErrInvalidDataSize is returned when incorrect data size is provided.
	ErrInvalidDataSize = errors.New("invalid data size")

	// ErrInvalidTrend is returned when trend doesn't match any of the
	// available trends.
	ErrInvalidTrend = errors.New("invalid trend")

	// ErrInvalidBand is returned when band doesn't match any of the
	// available bands.
	ErrInvalidBand = errors.New("invalid band")

	// ErrInvalidMA is returned when ma doesn't match any of the
	// availabble ma types.
	ErrInvalidMA = errors.New("invalid moving average")
)

// avg is a helper function that calculates average decimal number of
// given slice.
func avg(dd []decimal.Decimal) decimal.Decimal {
	var sum decimal.Decimal

	for i := range dd {
		sum = sum.Add(dd[i])
	}

	return sum.Div(decimal.NewFromInt(int64(len(dd))))
}

// sqrt is a helper function that calculated the square root of decimal number.
func sqrt(d decimal.Decimal) decimal.Decimal {
	f, _ := d.Float64()
	return decimal.NewFromFloat(math.Sqrt(f))
}

// mdev calculates mean deviation of given slice.
func mdev(dd []decimal.Decimal) decimal.Decimal {
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	res := decimal.Zero
	mean := avg(dd)

	for i := range dd {
		res = res.Add(dd[i].Sub(mean).Abs().Div(length))
	}

	return res
}

// sdev calculates standart deviation of given slice.
func sdev(dd []decimal.Decimal) decimal.Decimal {
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	res := decimal.Zero
	mean := avg(dd)

	for i := range dd {
		res = res.Add(dd[i].Sub(mean).Pow(decimal.NewFromInt(2)).Div(length))
	}

	return sqrt(res)
}

// Trend specifies which trend should be used.
type Trend int

const (
	// TrendUp specifies increasing value trend.
	TrendUp Trend = iota + 1

	// TrendDown specifies decreasing value value.
	TrendDown
)

// Validate checks whether the trend is one of
// supported trend types or not.
func (t Trend) Validate() error {
	switch t {
	case TrendUp, TrendDown:
		return nil
	default:
		return ErrInvalidTrend
	}
}

// MarshalText turns trend into appropriate string
// representation.
func (t Trend) MarshalText() ([]byte, error) {
	var v string

	switch t {
	case TrendUp:
		v = "up"
	case TrendDown:
		v = "down"
	default:
		return nil, ErrInvalidTrend
	}

	return []byte(v), nil
}

// UnmarshalText turns string to appropriate trend value.
func (t *Trend) UnmarshalText(d []byte) error {
	switch string(d) {
	case "up", "u":
		*t = TrendUp
	case "down", "d":
		*t = TrendDown
	default:
		return ErrInvalidTrend
	}

	return nil
}

// Band specifies which band should be used.
type Band int

// Available Bollinger Band indicator types.
const (
	BandUpper Band = iota + 1
	BandLower
	BandWidth
)

// Validate checks whether band is one of supported band types.
func (b Band) Validate() error {
	switch b {
	case BandUpper, BandLower, BandWidth:
		return nil
	default:
		return ErrInvalidBand
	}
}

// MarshalText turns band into appropriate string representation in JSON.
func (b Band) MarshalText() ([]byte, error) {
	var v string

	switch b {
	case BandUpper:
		v = "upper"
	case BandLower:
		v = "lower"
	case BandWidth:
		v = "width"
	default:
		return nil, ErrInvalidBand
	}

	return []byte(v), nil
}

// UnmarshalText turns JSON string to appropriate band value.
func (b *Band) UnmarshalText(d []byte) error {
	switch string(d) {
	case "upper", "u":
		*b = BandUpper
	case "lower", "l":
		*b = BandLower
	case "width", "w":
		*b = BandWidth
	default:
		return ErrInvalidBand
	}

	return nil
}

// MAType is a custom type that validates it to be only of existing
// moving average types.
type MAType int

// Available moving average indicator types.
const (
	MATypeDEMA MAType = iota + 1
	MATypeEMA
	MATypeHMA
	MATypeSMA
	MATypeWMA
)

// Initialize tries to construct new moving average based on the provided
// name.
func (mat MAType) Initialize(length int) (Indicator, error) {
	switch mat {
	case MATypeDEMA:
		return NewDEMA(length)
	case MATypeEMA:
		return NewEMA(length)
	case MATypeHMA:
		return NewHMA(length)
	case MATypeSMA:
		return NewSMA(length)
	case MATypeWMA:
		return NewWMA(length)
	default:
		return nil, ErrInvalidMA
	}
}

// MarshalText turns MAType into appropriate string representation in JSON.
func (mat MAType) MarshalText() ([]byte, error) {
	var v string

	switch mat {
	case MATypeDEMA:
		v = "dema"
	case MATypeEMA:
		v = "ema"
	case MATypeHMA:
		v = "hma"
	case MATypeSMA:
		v = "sma"
	case MATypeWMA:
		v = "wma"
	default:
		return nil, ErrInvalidMA
	}

	return []byte(v), nil
}

// UnmarshalText turns JSON string to appropriate moving average type value.
func (mat *MAType) UnmarshalText(d []byte) error {
	switch string(d) {
	case "dema":
		*mat = MATypeDEMA
	case "ema":
		*mat = MATypeEMA
	case "hma":
		*mat = MATypeHMA
	case "sma":
		*mat = MATypeSMA
	case "wma":
		*mat = MATypeWMA
	default:
		return ErrInvalidMA
	}

	return nil
}

// Indicator is an interface that every indicator should implement.
type Indicator interface {
	// Calc should return calculation results based on provided data
	// points slice.
	Calc([]decimal.Decimal) (decimal.Decimal, error)

	// Count should determine the total amount data points required for
	// the calculation.
	Count() int
}
