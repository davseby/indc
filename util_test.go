package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func assertEqualError(t *testing.T, exp, err error) {
	t.Helper()

	if exp != nil {
		if exp == assert.AnError { //nolint:goerr113 // direct check is needed
			assert.Error(t, err)
			return
		}

		assert.Equal(t, exp, err)

		return
	}

	assert.NoError(t, err)
}
func Test_mdev(t *testing.T) {
	cc := map[string]struct {
		Data   []decimal.Decimal
		Result decimal.Decimal
	}{
		"Successful calculation with no values": {
			Data:   []decimal.Decimal{},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation with one value": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
			},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(5),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(2),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res := mdev(c.Data)

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_sdev(t *testing.T) {
	cc := map[string]struct {
		Data   []decimal.Decimal
		Result decimal.Decimal
	}{
		"Successful calculation with no values": {
			Data:   []decimal.Decimal{},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation with one value": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
			},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(600),
				decimal.NewFromInt(470),
				decimal.NewFromInt(170),
				decimal.NewFromInt(430),
				decimal.NewFromInt(300),
			},
			Result: sqrt(decimal.NewFromInt(21704)),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res := sdev(c.Data)

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_Trend_Validate(t *testing.T) {
	cc := map[string]struct {
		Trend Trend
		Err   error
	}{
		"Invalid Trend": {
			Err: ErrInvalidTrend,
		},
		"Successful TrendUp validation": {
			Trend: TrendUp,
		},
		"Successful TrendDown validation": {
			Trend: TrendDown,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			err := c.Trend.Validate()
			assertEqualError(t, c.Err, err)
		})
	}
}

func Test_Trend_MarshalText(t *testing.T) {
	cc := map[string]struct {
		Trend Trend
		Text  string
		Err   error
	}{
		"Invalid Trend": {
			Err: ErrInvalidTrend,
		},
		"Successful TrendUp marshal": {
			Trend: TrendUp,
			Text:  "up",
		},
		"Successful TrendDown marshal": {
			Trend: TrendDown,
			Text:  "down",
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Trend.MarshalText()
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Text, string(res))
		})
	}
}

func Test_Trend_UnmarshalText(t *testing.T) {
	cc := map[string]struct {
		Text   string
		Result Trend
		Err    error
	}{
		"Invalid Trend": {
			Err: ErrInvalidTrend,
		},
		"Successful TrendUp unmarshal (long form)": {
			Text:   "up",
			Result: TrendUp,
		},
		"Successful TrendUp unmarshal (short form)": {
			Text:   "u",
			Result: TrendUp,
		},
		"Successful TrendDown unmarshal  (long form)": {
			Text:   "down",
			Result: TrendDown,
		},
		"Successful TrendDown unmarshal  (short form)": {
			Text:   "d",
			Result: TrendDown,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			var tr Trend
			err := tr.UnmarshalText([]byte(c.Text))
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, tr)
		})
	}
}

func Test_Band_Validate(t *testing.T) {
	cc := map[string]struct {
		Band Band
		Err  error
	}{
		"Invalid Band": {
			Err: ErrInvalidBand,
		},
		"Successful BandUpper validation": {
			Band: BandUpper,
		},
		"Successful BandLower validation": {
			Band: BandLower,
		},
		"Successful BandWidth validation": {
			Band: BandWidth,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			err := c.Band.Validate()
			assertEqualError(t, c.Err, err)
		})
	}
}

func Test_Band_MarshalText(t *testing.T) {
	cc := map[string]struct {
		Band Band
		Text string
		Err  error
	}{
		"Invalid Band": {
			Err: ErrInvalidBand,
		},
		"Successful BandUpper marshal": {
			Band: BandUpper,
			Text: "upper",
		},
		"Successful BandLower marshal": {
			Band: BandLower,
			Text: "lower",
		},
		"Successful BandWidth marshal": {
			Band: BandWidth,
			Text: "width",
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Band.MarshalText()
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Text, string(res))
		})
	}
}

func Test_Band_UnmarshalText(t *testing.T) {
	cc := map[string]struct {
		Text   string
		Result Band
		Err    error
	}{
		"Invalid Band": {
			Err: ErrInvalidBand,
		},
		"Successful BandUpper unmarshal (long form)": {
			Text:   "upper",
			Result: BandUpper,
		},
		"Successful BandUpper unmarshal (short form)": {
			Text:   "u",
			Result: BandUpper,
		},
		"Successful BandLower unmarshal  (long form)": {
			Text:   "lower",
			Result: BandLower,
		},
		"Successful BandLower unmarshal  (short form)": {
			Text:   "l",
			Result: BandLower,
		},
		"Successful BandWidth unmarshal  (long form)": {
			Text:   "width",
			Result: BandWidth,
		},
		"Successful BandWidth unmarshal  (short form)": {
			Text:   "w",
			Result: BandWidth,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			var b Band
			err := b.UnmarshalText([]byte(c.Text))
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, b)
		})
	}
}

func Test_MAType_Initialize(t *testing.T) {
	cc := map[string]struct {
		Type      MAType
		Length    int
		Indicator Indicator
		Err       error
	}{
		"Invalid MAType": {
			Err: ErrInvalidMA,
		},
		"Successful MATypeDEMA initialization": {
			Type:   MATypeDEMA,
			Length: 1,
			Indicator: DEMA{
				valid: true,
				ema: EMA{
					valid: true,
					sma: SMA{
						valid:  true,
						length: 1,
					},
				},
			},
		},
		"Successful MATypeEMA initialization": {
			Type:   MATypeEMA,
			Length: 1,
			Indicator: EMA{
				valid: true,
				sma: SMA{
					valid:  true,
					length: 1,
				},
			},
		},
		"Successful MATypeHMA initialization": {
			Type:   MATypeHMA,
			Length: 1,
			Indicator: HMA{
				valid: true,
				wma: WMA{
					valid:  true,
					length: 1,
				},
			},
		},
		"Successful MATypeSMA initialization": {
			Type:   MATypeSMA,
			Length: 1,
			Indicator: SMA{
				valid:  true,
				length: 1,
			},
		},
		"Successful MATypeWMA initialization": {
			Type:   MATypeWMA,
			Length: 1,
			Indicator: WMA{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			ma, err := c.Type.Initialize(c.Length)
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Indicator, ma)
		})
	}
}

func Test_MAType_MarshalText(t *testing.T) {
	cc := map[string]struct {
		Type MAType
		Text string
		Err  error
	}{
		"Invalid MAType": {
			Type: 70,
			Err:  ErrInvalidMA,
		},
		"Successful MATypeDEMA marshal": {
			Type: MATypeDEMA,
			Text: "dema",
		},
		"Successful MATypeEMA marshal": {
			Type: MATypeEMA,
			Text: "ema",
		},
		"Successful MATypeHMA marshal": {
			Type: MATypeHMA,
			Text: "hma",
		},
		"Successful MATypeSMA marshal": {
			Type: MATypeSMA,
			Text: "sma",
		},
		"Successful MATypeWMA marshal": {
			Type: MATypeWMA,
			Text: "wma",
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Type.MarshalText()
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Text, string(res))
		})
	}
}

func Test_MAType_UnmarshalText(t *testing.T) {
	cc := map[string]struct {
		Text   string
		Result MAType
		Err    error
	}{
		"Invalid MAType": {
			Text: "70",
			Err:  ErrInvalidMA,
		},
		"Successful MATypeDEMA unmarshal": {
			Text:   "dema",
			Result: MATypeDEMA,
		},
		"Successful MATypeEMA unmarshal": {
			Text:   "ema",
			Result: MATypeEMA,
		},
		"Successful MATypeHMA unmarshal": {
			Text:   "hma",
			Result: MATypeHMA,
		},
		"Successful MATypeSMA unmarshal": {
			Text:   "sma",
			Result: MATypeSMA,
		},
		"Successful MATypeWMA unmarshal": {
			Text:   "wma",
			Result: MATypeWMA,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			var mat MAType
			err := mat.UnmarshalText([]byte(c.Text))
			assertEqualError(t, c.Err, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, mat)
		})
	}
}
