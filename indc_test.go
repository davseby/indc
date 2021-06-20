package indc

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewAroon(t *testing.T) {
	cc := map[string]struct {
		Trend  Trend
		Length int
		Result Aroon
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new Aroon": {
			Trend:  TrendDown,
			Length: 5,
			Result: Aroon{
				valid:  true,
				trend:  TrendDown,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewAroon(c.Trend, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Aroon_validate(t *testing.T) {
	cc := map[string]struct {
		Aroon Aroon
		Error error
	}{
		"Invalid trend": {
			Aroon: Aroon{
				valid:  false,
				trend:  70,
				length: 5,
			},
			Error: ErrInvalidTrend,
		},
		"Invalid length": {
			Aroon: Aroon{
				trend: TrendDown,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			Aroon: Aroon{
				trend:  TrendUp,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Aroon.validate())
			if c.Error == nil {
				assert.True(t, c.Aroon.valid)
			}
		})
	}
}

func Test_Aroon_Calc(t *testing.T) {
	cc := map[string]struct {
		Aroon  Aroon
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			Aroon: Aroon{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Aroon: Aroon{
				valid:  true,
				trend:  TrendDown,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with TrendUp": {
			Aroon: Aroon{
				valid:  true,
				trend:  TrendUp,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: decimal.NewFromInt(40),
		},
		"Successful calculation with TrendDown": {
			Aroon: Aroon{
				valid:  true,
				trend:  TrendDown,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: _hundred,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Aroon.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_Aroon_Count(t *testing.T) {
	assert.Equal(t, 5, Aroon{
		length: 5,
	}.Count())
}

func Test_NewBB(t *testing.T) {
	cc := map[string]struct {
		Percent bool
		Band    Band
		StdDev  decimal.Decimal
		Length  int
		Result  BB
		Error   error
	}{
		"NewSMA returns an error": {
			Error: assert.AnError,
		},
		"Validate returns an error": {
			Length: 1,
			Error:  ErrInvalidBand,
		},
		"Successfully created new BB": {
			Percent: true,
			Band:    BandUpper,
			StdDev:  decimal.RequireFromString("2.5"),
			Length:  5,
			Result: BB{
				valid:   true,
				percent: true,
				band:    BandUpper,
				stdDev:  decimal.RequireFromString("2.5"),
				sma: SMA{
					length: 5,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewBB(c.Percent, c.Band, c.StdDev, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_BB_validate(t *testing.T) {
	cc := map[string]struct {
		BB    BB
		Error error
	}{
		"Invalid band": {
			BB: BB{
				band:   70,
				stdDev: decimal.Decimal{},
				sma: SMA{
					length: 5,
				},
			},
			Error: ErrInvalidBand,
		},
		"Invalid BB band width configuration": {
			BB: BB{
				percent: true,
				band:    BandWidth,
				stdDev:  decimal.Decimal{},
				sma: SMA{
					length: 1,
				},
			},
			Error: errors.New("invalid bb configuration"),
		},
		"Successfully validated": {
			BB: BB{
				band:   BandUpper,
				stdDev: decimal.Decimal{},
				sma: SMA{
					length: 1,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.BB.validate())
			if c.Error == nil {
				assert.True(t, c.BB.valid)
			}
		})
	}
}

func Test_BB_Calc(t *testing.T) {
	cc := map[string]struct {
		BB     BB
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			BB: BB{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			BB: BB{
				valid: true,
				band:  BandUpper,
				sma: SMA{
					valid:  true,
					length: 5,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with BandUpper": {
			BB: BB{
				valid:  true,
				band:   BandUpper,
				stdDev: decimal.RequireFromString("1"),
				sma: SMA{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(35),
				decimal.NewFromInt(40),
				decimal.NewFromInt(38),
				decimal.NewFromInt(32),
			},
			Result: decimal.RequireFromString("38.68781778"),
		},
		"Successful calculation with BandUpper using percent": {
			BB: BB{
				valid:   true,
				percent: true,
				band:    BandUpper,
				stdDev:  decimal.RequireFromString("1"),
				sma: SMA{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(35),
				decimal.NewFromInt(40),
				decimal.NewFromInt(38),
				decimal.NewFromInt(32),
			},
			Result: decimal.RequireFromString("10.53662224"),
		},
		"Successful calculation with BandLower": {
			BB: BB{
				valid:  true,
				band:   BandLower,
				stdDev: decimal.RequireFromString("1"),
				sma: SMA{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(35),
				decimal.NewFromInt(40),
				decimal.NewFromInt(38),
				decimal.NewFromInt(32),
			},
			Result: decimal.RequireFromString("31.31218222"),
		},
		"Successful calculation with BandLower using percent": {
			BB: BB{
				valid:   true,
				percent: true,
				band:    BandLower,
				stdDev:  decimal.RequireFromString("1"),
				sma: SMA{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(35),
				decimal.NewFromInt(40),
				decimal.NewFromInt(38),
				decimal.NewFromInt(32),
			},
			Result: decimal.RequireFromString("-10.53662224"),
		},
		"Successful calculation with BandWidth": {
			BB: BB{
				valid:  true,
				band:   BandWidth,
				stdDev: decimal.RequireFromString("2"),
				sma: SMA{
					length: 20,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.RequireFromString("63.98"),
				decimal.RequireFromString("64.17"),
				decimal.RequireFromString("64.71"),
				decimal.RequireFromString("64.75"),
				decimal.RequireFromString("63.94"),
				decimal.RequireFromString("63.82"),
				decimal.RequireFromString("63.19"),
				decimal.RequireFromString("62.84"),
				decimal.RequireFromString("62.25"),
				decimal.RequireFromString("63.20"),
				decimal.RequireFromString("63.02"),
				decimal.RequireFromString("63.35"),
				decimal.RequireFromString("64.21"),
				decimal.RequireFromString("64.91"),
				decimal.RequireFromString("64.05"),
				decimal.RequireFromString("63.28"),
				decimal.RequireFromString("62.78"),
				decimal.RequireFromString("62.36"),
				decimal.RequireFromString("63.19"),
				decimal.RequireFromString("64.69"),
			},
			Result: decimal.RequireFromString("4.91959301"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.BB.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.Round(8).String(), res.Round(8).String())
		})
	}
}

func Test_BB_Count(t *testing.T) {
	assert.Equal(t, 1, BB{sma: SMA{length: 1}}.Count())
}

func Test_NewCCI(t *testing.T) {
	cc := map[string]struct {
		Type   MAType
		Length int
		Factor decimal.Decimal
		Result CCI
		Error  error
	}{
		"NewSMA returns an error": {
			Error: assert.AnError,
		},
		"Invalid provided moving average type": {
			Length: 1,
			Error:  errors.New("invalid moving average"),
		},
		"Invalid factor": {
			Type:   MATypeSMA,
			Length: 1,
			Factor: decimal.RequireFromString("-1"),
			Error:  errors.New("invalid factor"),
		},
		"Successfully created new CCI with default factor": {
			Type:   MATypeSMA,
			Length: 10,
			Factor: decimal.Zero,
			Result: CCI{
				valid: true,
				ma: SMA{
					length: 10,
					valid:  true,
				},
				factor: decimal.RequireFromString("0.015"),
			},
		},
		"Successfully created new CCI": {
			Type:   MATypeSMA,
			Length: 10,
			Factor: _hundred,
			Result: CCI{
				valid: true,
				ma: SMA{
					length: 10,
					valid:  true,
				},
				factor: _hundred,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewCCI(c.Type, c.Length, c.Factor)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_CCI_validate(t *testing.T) {
	cc := map[string]struct {
		CCI   CCI
		Error error
	}{
		"Invalid factor": {
			CCI: CCI{
				valid: false,
				ma: SMA{
					length: 1,
				},
				factor: decimal.NewFromInt(-1),
			},
			Error: errors.New("invalid factor"),
		},
		"Successfully validated": {
			CCI: CCI{
				valid: false,
				ma: SMA{
					length: 1,
				},
				factor: decimal.RequireFromString("1"),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.CCI.validate())
			if c.Error == nil {
				assert.True(t, c.CCI.valid)
			}
		})
	}
}

func Test_CCI_Calc(t *testing.T) {
	cc := map[string]struct {
		CCI    CCI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			CCI:   CCI{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			CCI: CCI{
				valid: true,
				ma: SMA{
					length: 31,
				},
				factor: decimal.RequireFromString("0.015"),
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Invalid SMA calc": {
			CCI: CCI{
				valid:  true,
				ma:     SMA{},
				factor: decimal.RequireFromString("0.015"),
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful handled division by 0": {
			CCI: CCI{
				valid: true,
				ma: SMA{
					length: 1,
					valid:  true,
				},
				factor: decimal.RequireFromString("0.015"),
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(3),
			},
			Result: decimal.Zero,
		},
		"Successful calculation": {
			CCI: CCI{
				valid: true,
				ma: SMA{
					length: 3,
					valid:  true,
				},
				factor: decimal.RequireFromString("0.015"),
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(3),
				decimal.NewFromInt(6),
				decimal.NewFromInt(9),
			},
			Result: decimal.NewFromInt(100),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.CCI.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_CCI_Count(t *testing.T) {
	assert.Equal(t, 10, CCI{
		ma: SMA{
			length: 10,
		},
	}.Count())
}

func Test_NewDEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result DEMA
		Error  error
	}{
		"NewEMA returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new DEMA": {
			Length: 1,
			Result: DEMA{
				valid: true,
				ema: EMA{
					sma: SMA{
						length: 1,
						valid:  true,
					},
					valid: true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewDEMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_DEMA_Calc(t *testing.T) {
	cc := map[string]struct {
		DEMA   DEMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			DEMA:  DEMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			DEMA: DEMA{
				valid: true,
				ema: EMA{
					sma: SMA{
						length: 3,
						valid:  true,
					},
					valid: true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			DEMA: DEMA{
				valid: true,
				ema: EMA{
					sma: SMA{
						length: 3,
						valid:  true,
					},
					valid: true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(1),
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			Result: decimal.RequireFromString("6.75"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.DEMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_DEMA_Count(t *testing.T) {
	assert.Equal(t, 29, DEMA{
		valid: false,
		ema: EMA{
			sma: SMA{
				length: 15,
			},
		},
	}.Count())
}

func Test_NewEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successfully created new EMA": {
			Length: 1,
			Result: EMA{
				valid: true,
				sma: SMA{
					length: 1,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewEMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_EMA_Calc(t *testing.T) {
	cc := map[string]struct {
		EMA    EMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			EMA:   EMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			EMA: EMA{
				valid: true,
				sma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			EMA: EMA{
				valid: true,
				sma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(1),
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			Result: decimal.RequireFromString("4.75"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.EMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_EMA_CalcNext(t *testing.T) {
	cc := map[string]struct {
		EMA    EMA
		Last   decimal.Decimal
		Next   decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			EMA:   EMA{},
			Error: ErrInvalidIndicator,
		},
		"Successful calculation": {
			EMA: EMA{
				valid: true,
				sma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Last:   decimal.NewFromInt(5),
			Next:   decimal.NewFromInt(5),
			Result: decimal.NewFromInt(5),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.EMA.CalcNext(c.Last, c.Next)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_EMA_Count(t *testing.T) {
	assert.Equal(t, 29, EMA{
		sma: SMA{
			length: 15,
		},
	}.Count())
}

func Test_EMA_multiplier(t *testing.T) {
	assert.Equal(t, decimal.RequireFromString("0.5").String(), EMA{
		sma: SMA{
			length: 3,
		},
	}.multiplier().String())
}

func Test_NewHMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result HMA
		Error  error
	}{
		"NewWMA returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new HMA": {
			Length: 2,
			Result: HMA{
				valid: true,
				wma: WMA{
					length: 2,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewHMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_HMA_Calc(t *testing.T) {
	cc := map[string]struct {
		HMA    HMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			HMA:   HMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			HMA: HMA{
				valid: true,
				wma: WMA{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			HMA: HMA{
				valid: true,
				wma: WMA{
					length: 4,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(32),
				decimal.NewFromInt(29),
				decimal.NewFromInt(38),
				decimal.NewFromInt(34),
				decimal.NewFromInt(29),
			},
			Result: decimal.RequireFromString("33.8"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.HMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.Round(8).String(), res.Round(8).String())
		})
	}
}

func Test_HMA_Count(t *testing.T) {
	assert.Equal(t, 17, HMA{
		wma: WMA{
			length: 15,
		},
	}.Count())
}

func Test_NewROC(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ROC
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new ROC": {
			Length: 1,
			Result: ROC{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewROC(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_ROC_validate(t *testing.T) {
	cc := map[string]struct {
		ROC   ROC
		Error error
	}{
		"Invalid length": {
			ROC: ROC{
				length: -1,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			ROC: ROC{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.ROC.validate())
			if c.Error == nil {
				assert.True(t, c.ROC.valid)
			}
		})
	}
}

func Test_ROC_Calc(t *testing.T) {
	cc := map[string]struct {
		ROC    ROC
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			ROC:   ROC{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			ROC: ROC{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			ROC: ROC{
				valid:  true,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(7),
				decimal.NewFromInt(16),
				decimal.NewFromInt(24),
				decimal.NewFromInt(16),
				decimal.NewFromInt(10),
			},
			Result: decimal.RequireFromString("-30"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.ROC.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_ROC_Count(t *testing.T) {
	assert.Equal(t, 15, ROC{
		length: 15,
	}.Count())
}

func Test_NewRSI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result RSI
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new RSI": {
			Length: 1,
			Result: RSI{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewRSI(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_RSI_validate(t *testing.T) {
	cc := map[string]struct {
		RSI   RSI
		Error error
	}{
		"Invalid length": {
			RSI: RSI{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			RSI: RSI{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.RSI.validate())
			if c.Error == nil {
				assert.True(t, c.RSI.valid)
			}
		})
	}
}

func Test_RSI_Calc(t *testing.T) {
	cc := map[string]struct {
		RSI    RSI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			RSI:   RSI{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when average gain 0": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(16),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation when average loss 0": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(4),
				decimal.NewFromInt(8),
			},
			Result: _hundred,
		},
		"Successful calculation": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(8),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(50),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.RSI.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_RSI_Count(t *testing.T) {
	assert.Equal(t, 15, RSI{
		length: 15,
	}.Count())
}

func Test_NewSMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result SMA
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new SMA": {
			Length: 1,
			Result: SMA{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewSMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_SMA_validate(t *testing.T) {
	cc := map[string]struct {
		SMA   SMA
		Error error
	}{
		"Invalid length": {
			SMA: SMA{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			SMA: SMA{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.SMA.validate())
			if c.Error == nil {
				assert.True(t, c.SMA.valid)
			}
		})
	}
}

func Test_SMA_Calc(t *testing.T) {
	cc := map[string]struct {
		SMA    SMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			SMA:   SMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			SMA: SMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			SMA: SMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromInt(31),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.SMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_SMA_Count(t *testing.T) {
	assert.Equal(t, 15, SMA{
		length: 15,
	}.Count())
}

func Test_NewSRSI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result SRSI
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new SRSI": {
			Length: 1,
			Result: SRSI{
				valid: true,
				rsi: RSI{
					length: 1,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewSRSI(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_SRSI_Calc(t *testing.T) {
	cc := map[string]struct {
		SRSI   SRSI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			SRSI:  SRSI{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			SRSI: SRSI{
				valid: true,
				rsi: RSI{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successfully handled division by 0": {
			SRSI: SRSI{
				valid: true,
				rsi: RSI{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(8),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.Zero,
		},
		"Successful calculation": {
			SRSI: SRSI{
				valid: true,
				rsi: RSI{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(11),
				decimal.NewFromInt(12),
				decimal.NewFromInt(11),
				decimal.NewFromInt(11),
				decimal.NewFromInt(11),
			},
			Result: decimal.RequireFromString("0.5"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.SRSI.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_SRSI_Count(t *testing.T) {
	assert.Equal(t, 29, SRSI{
		valid: false,
		rsi: RSI{
			length: 15,
		},
	}.Count())
}

func Test_NewStoch(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Stoch
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new Stoch": {
			Length: 1,
			Result: Stoch{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewStoch(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Stoch_validate(t *testing.T) {
	cc := map[string]struct {
		Stoch Stoch
		Error error
	}{
		"Invalid length": {
			Stoch: Stoch{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			Stoch: Stoch{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Stoch.validate())
			if c.Error == nil {
				assert.True(t, c.Stoch.valid)
			}
		})
	}
}

func Test_Stoch_Calc(t *testing.T) {
	cc := map[string]struct {
		Stoch  Stoch
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			Stoch: Stoch{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when new lows are reached": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
		"Successfully handled division by 0": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(150),
				decimal.NewFromInt(150),
			},
			Result: decimal.Zero,
		},
		"Successful calculation when new highs are reached": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(120),
				decimal.NewFromInt(145),
				decimal.NewFromInt(135),
			},
			Result: decimal.NewFromInt(60),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Stoch.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_NewWMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result WMA
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new WMA": {
			Length: 1,
			Result: WMA{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewWMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_WMA_validate(t *testing.T) {
	cc := map[string]struct {
		WMA   WMA
		Error error
	}{
		"Invalid length": {
			WMA: WMA{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			WMA: WMA{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.WMA.validate())
			if c.Error == nil {
				assert.True(t, c.WMA.valid)
			}
		})
	}
}

func Test_WMA_Calc(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			WMA:   WMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			WMA: WMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			WMA: WMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(30),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromInt(31),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.WMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_WMA_Count(t *testing.T) {
	assert.Equal(t, 15, WMA{
		length: 15,
	}.Count())
}
