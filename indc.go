// Package indc provides types and functions to calculate values of various
// market indicators.
package indc

import (
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

// Aroon holds all the necessary information needed to calculate Aroon.
// The zero value is not usable.
type Aroon struct {
	// valid specifies whether Aroon paremeters were validated.
	valid bool

	// trend specifies which Aroon trend to use during the
	// calculation process.
	trend Trend

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewAroon validates provided configuration options and
// creates new Aroon indicator instance.
func NewAroon(trend Trend, length int) (Aroon, error) {
	aroon := Aroon{
		trend:  trend,
		length: length,
	}

	if err := aroon.validate(); err != nil {
		return Aroon{}, err
	}

	return aroon, nil
}

// validate checks whether the indicator has valid configuration properties.
func (aroon *Aroon) validate() error {
	if err := aroon.trend.Validate(); err != nil {
		return err
	}

	if aroon.length < 1 {
		return ErrInvalidLength
	}

	aroon.valid = true

	return nil
}

// Calc calculates Aroon from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/a/aroon.asp.
// All credits are due to Tushar Chande who developed Aroon indicator.
func (aroon Aroon) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !aroon.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != aroon.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := dd[0]
	prd := decimal.Zero

	refresh := func(val decimal.Decimal) bool {
		fn := res.LessThanOrEqual
		if aroon.trend == TrendDown {
			fn = res.GreaterThanOrEqual
		}

		return fn(val)
	}

	for i := 0; i < len(dd); i++ {
		if refresh(dd[i]) {
			res = dd[i]
			prd = decimal.NewFromInt(int64(aroon.length - i - 1))
		}
	}

	return decimal.NewFromInt(int64(aroon.length)).Sub(prd).
		Mul(_hundred).Div(decimal.NewFromInt(int64(aroon.length))), nil
}

// Count determines the total amount of data points needed for Aroon
// calculation.
func (aroon Aroon) Count() int {
	return aroon.length
}

// BB holds all the necessary information needed to calculate Bollinger Bands.
// The zero value is not usable.
type BB struct {
	// valid specifies whether BB paremeters were validated.
	valid bool

	// percent specifies whether returned number should be in units (if false)
	// or percent (true).
	percent bool

	// band specifies which bollinger band to calculate.
	band Band

	// stdDev specifies how to adjust standard deviation.
	stdDev decimal.Decimal

	// sma specifies SMA indicator configuration.
	sma SMA
}

// NewBB validates provided configuration options and creates
// new BB indicator.
func NewBB(percent bool, band Band, stdDev decimal.Decimal, length int) (BB, error) {
	sma, err := NewSMA(length)
	if err != nil {
		return BB{}, err
	}

	bb := BB{
		percent: percent,
		band:    band,
		stdDev:  stdDev,
		sma:     sma,
	}

	if err := bb.validate(); err != nil {
		return BB{}, err
	}

	return bb, nil
}

// validate checks whether the indicator has valid configuration properties.
func (bb *BB) validate() error {
	if err := bb.band.Validate(); err != nil {
		return err
	}

	if bb.percent && bb.band == BandWidth {
		return errors.New("invalid bb configuration")
	}

	bb.valid = true

	return nil
}

// Calc calculates BB from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/b/bollingerbands.asp.
// All credits are due to John Bollinger who developed BB indicator.
func (bb BB) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !bb.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != bb.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res, err := bb.sma.Calc(dd)
	if err != nil {
		// unlikely to happen
		return decimal.Zero, err
	}

	sdev := sdev(dd).Mul(bb.stdDev)

	switch bb.band {
	case BandUpper:
		if bb.percent {
			return res.Add(sdev).Div(res).Sub(_one).Mul(_hundred), nil
		}

		return res.Add(sdev), nil
	case BandLower:
		if bb.percent {
			return res.Sub(sdev).Div(res).Sub(_one).Mul(_hundred), nil
		}

		return res.Sub(sdev), nil
	default: // BB is validated, only BandWidth is left.
		return res.Add(sdev).Sub(res.Sub(sdev)).Div(res).Mul(_hundred), nil
	}
}

// Count determines the total amount of data points needed for BB
// calculation.
func (bb BB) Count() int {
	return bb.sma.Count()
}

// CCI holds all the necessary information needed to calculate commodity
// channel index.
// The zero value is not usable.
type CCI struct {
	// valid specifies whether CCI paremeters were validated.
	valid bool

	// ma specifies moving average indicator configuration.
	ma Indicator

	// factor is used to scale CCI to provide more readable numbers.
	// default is 0.015f.
	factor decimal.Decimal
}

// NewCCI validates provided configuration options and creates
// new CCI indicator.
// If provided factor is zero, default value is going to be used (0.015f).
func NewCCI(mat MAType, length int, factor decimal.Decimal) (CCI, error) {
	if factor.Equal(decimal.Zero) {
		factor = decimal.RequireFromString("0.015")
	}

	ma, err := mat.Initialize(length)
	if err != nil {
		return CCI{}, err
	}

	cci := CCI{
		ma:     ma,
		factor: factor,
	}

	if err := cci.validate(); err != nil {
		return CCI{}, err
	}

	return cci, nil
}

// validate checks whether the indicator has valid configuration properties.
func (cci *CCI) validate() error {
	if cci.factor.LessThanOrEqual(decimal.Zero) {
		return errors.New("invalid factor")
	}

	cci.valid = true

	return nil
}

// Calc calculates CCI from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/c/commoditychannelindex.asp.
// All credits are due to Donald Lambert who developed CCI indicator.
func (cci CCI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !cci.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != cci.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res, err := cci.ma.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	dnm := cci.factor.Mul(mdev(dd))

	if dnm.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return dd[len(dd)-1].Sub(res).Div(dnm), nil
}

// Count determines the total amount of data points needed for CCI
// calculation.
func (cci CCI) Count() int {
	return cci.ma.Count()
}

// DEMA holds all the necessary information needed to calculate
// double exponential moving average.
// The zero value is not usable.
type DEMA struct {
	// valid specifies whether DEMA paremeters were validated.
	valid bool

	// ema specifies what ema should be used for dema calculations.
	ema EMA
}

// NewDEMA validates provided configuration options and creates
// new DEMA indicator.
func NewDEMA(length int) (DEMA, error) {
	ema, err := NewEMA(length)
	if err != nil {
		return DEMA{}, err
	}

	return DEMA{
		valid: true,
		ema:   ema,
	}, nil
}

// Calc calculates DEMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/d/double-exponential-moving-average.asp.
// All credits are due to Patrick Mulloy who developed DEMA indicator.
func (dema DEMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !dema.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != dema.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	pres := make([]decimal.Decimal, dema.ema.sma.length)

	var err error

	pres[0], err = dema.ema.sma.Calc(dd[:dema.ema.sma.length])
	if err != nil {
		// unlikely to happen
		return decimal.Zero, err
	}

	for i := dema.ema.sma.length; i < len(dd); i++ {
		pres[i-dema.ema.sma.length+1], err = dema.ema.CalcNext(pres[i-dema.ema.sma.length], dd[i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	res := pres[0]

	for i := 0; i < len(pres); i++ {
		res, err = dema.ema.CalcNext(res, pres[i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	return res, nil
}

// Count determines the total amount of data points needed for DEMA
// calculation.
func (dema DEMA) Count() int {
	return dema.ema.Count()
}

// EMA holds all the necessary information needed to calculate exponential
// moving average.
// The zero value is not usable.
type EMA struct {
	// valid specifies whether DEMA paremeters were validated.
	valid bool

	// sma specifies what sma should be used for ema calculations.
	sma SMA
}

// NewEMA validates provided configuration options and
// creates new EMA indicator.
func NewEMA(length int) (EMA, error) {
	sma, err := NewSMA(length)
	if err != nil {
		return EMA{}, err
	}

	return EMA{
		valid: true,
		sma:   sma,
	}, nil
}

// Calc calculates EMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/e/ema.asp.
func (ema EMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !ema.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != ema.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res, err := ema.sma.Calc(dd[:ema.sma.length])
	if err != nil {
		// unlikely to happen
		return decimal.Zero, err
	}

	for i := ema.sma.length; i < len(dd); i++ {
		res, err = ema.CalcNext(res, dd[i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	return res, nil
}

// CalcNext calculates sequential EMA by using previous EMA.
func (ema EMA) CalcNext(lres, dec decimal.Decimal) (decimal.Decimal, error) {
	if !ema.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	mtp := ema.multiplier()

	return dec.Mul(mtp).Add(lres.Mul(decimal.NewFromInt(1).Sub(mtp))), nil
}

// multiplier calculates EMA multiplier.
func (ema EMA) multiplier() decimal.Decimal {
	return decimal.NewFromInt(2).Div(decimal.NewFromInt(int64(ema.sma.length) + 1))
}

// Count determines the total amount of data points needed for EMA
// calculation.
func (ema EMA) Count() int {
	return ema.sma.length*2 - 1
}

// HMA holds all the necessary information needed to calculate
// hull moving average.
// The zero value is not usable.
type HMA struct {
	// valid specifies whether HMA paremeters were validated.
	valid bool

	// wma specifies the base moving average.
	wma WMA
}

// NewHMA validates provided configuration options and
// creates new HMA indicator.
func NewHMA(length int) (HMA, error) {
	wma, err := NewWMA(length)
	if err != nil {
		return HMA{}, err
	}

	return HMA{
		valid: true,
		wma:   wma,
	}, nil
}

// Calc calculates HMA from the provided data points slice.
// Calculation is based on formula provided by fidelity.
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/hull-moving-average.
// All credits are due to Alan Hull who developed HMA indicator.
func (h HMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !h.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != h.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	wma1 := WMA{length: h.wma.length / 2, valid: true}
	wma2 := WMA{length: int(math.Sqrt(float64(h.wma.length))), valid: true}

	res := make([]decimal.Decimal, wma2.length)

	for i := 0; i < wma2.length; i++ {
		res1, err := wma1.Calc(dd[i : wma1.length+i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}

		res2, err := h.wma.Calc(dd[i : h.wma.length+i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}

		res[i] = res1.Mul(decimal.NewFromInt(2)).Sub(res2)
	}

	return wma2.Calc(res)
}

// Count determines the total amount of data points needed for HMA
// calculation.
func (h HMA) Count() int {
	return int(math.Sqrt(float64(h.wma.length))) + h.wma.length - 1
}

// ROC holds all the necessary information needed to calculate rate
// of change.
// The zero value is not usable.
type ROC struct {
	// valid specifies whether ROC paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewROC validates provided configuration options and
// creates new ROC indicator.
func NewROC(length int) (ROC, error) {
	roc := ROC{length: length}

	if err := roc.validate(); err != nil {
		return ROC{}, err
	}

	return roc, nil
}

// validate checks whether the indicator has valid configuration properties.
func (roc *ROC) validate() error {
	if roc.length < 1 {
		return ErrInvalidLength
	}

	roc.valid = true

	return nil
}

// Calc calculates ROC from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/p/pricerateofchange.asp.
func (roc ROC) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !roc.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != roc.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	curr := dd[0]
	last := dd[len(dd)-1]

	return curr.Div(last).Sub(_one).Mul(_hundred), nil
}

// Count determines the total amount of data points needed for ROC
// calculation.
func (roc ROC) Count() int {
	return roc.length
}

// RSI holds all the necessary information needed to calculate relative
// strength index.
// The zero value is not usable.
type RSI struct {
	// valid specifies whether RSI paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewRSI validates provided configuration options and
// creates new RSI indicator.
func NewRSI(length int) (RSI, error) {
	rsi := RSI{
		length: length,
	}

	if err := rsi.validate(); err != nil {
		return RSI{}, err
	}

	return rsi, nil
}

// validate checks whether the indicator has valid configuration properties.
func (rsi *RSI) validate() error {
	if rsi.length < 1 {
		return ErrInvalidLength
	}

	rsi.valid = true

	return nil
}

// Calc calculates RSI from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/r/rsi.asp.
// All credits are due to J. Welles Wilder Jr. who developed RSI indicator.
func (rsi RSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !rsi.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != rsi.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	ag := decimal.Zero
	al := decimal.Zero
	length := decimal.NewFromInt(int64(rsi.length))

	for i := 1; i < len(dd); i++ {
		if dd[i].Sub(dd[i-1]).LessThan(decimal.Zero) {
			al = al.Add(dd[i].Sub(dd[i-1]).Abs())
		} else {
			ag = ag.Add(dd[i].Sub(dd[i-1]))
		}
	}

	if ag == decimal.Zero {
		return decimal.NewFromInt(0), nil
	}

	if al == decimal.Zero {
		return _hundred, nil
	}

	ag = ag.Div(length)

	al = al.Div(length)

	return _hundred.Sub(_hundred.Div(decimal.NewFromInt(1).Add(ag.Div(al)))), nil
}

// Count determines the total amount of data points needed for RSI
// calculation.
func (rsi RSI) Count() int {
	return rsi.length
}

// SMA holds all the necessary information needed to calculate simple
// moving average.
// The zero value is not usable.
type SMA struct {
	// valid specifies whether SMA paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewSMA validates provided configuration options and
// creates new SMA indicator.
func NewSMA(length int) (SMA, error) {
	sma := SMA{
		length: length,
	}

	if err := sma.validate(); err != nil {
		return SMA{}, err
	}

	return sma, nil
}

// validate checks whether the indicator has valid configuration properties.
func (sma *SMA) validate() error {
	if sma.length < 1 {
		return ErrInvalidLength
	}

	sma.valid = true

	return nil
}

// Calc calculates SMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/s/sma.asp.
func (sma SMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !sma.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != sma.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := decimal.Zero
	for i := 0; i < len(dd); i++ {
		res = res.Add(dd[i])
	}

	return res.Div(decimal.NewFromInt(int64(sma.length))), nil
}

// Count determines the total amount of data points needed for SMA
// calculation.
func (sma SMA) Count() int {
	return sma.length
}

// SRSI holds all the necessary information needed to calculate stoch
// relative strength index.
// The zero value is not usable.
type SRSI struct {
	// valid specifies whether SRSI paremeters were validated.
	valid bool

	// rsi specifies the base relative strength index.
	rsi RSI
}

// NewSRSI validates provided configuration options and
// creates new SRSI indicator.
func NewSRSI(length int) (SRSI, error) {
	rsi, err := NewRSI(length)
	if err != nil {
		return SRSI{}, err
	}

	return SRSI{
		valid: true,
		rsi:   rsi,
	}, nil
}

// Calc calculates SRSI from the provided data slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/s/stochrsi.asp.
func (srsi SRSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !srsi.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != srsi.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := make([]decimal.Decimal, srsi.rsi.length)

	var err error
	for i := 0; i < srsi.rsi.length; i++ {
		res[i], err = srsi.rsi.Calc(dd[i : srsi.rsi.length+i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	curr := res[0]
	max := res[0]
	min := res[0]

	for i := 1; i < len(res); i++ {
		if max.LessThan(res[i]) {
			max = res[i]
		}

		if min.GreaterThan(res[i]) {
			min = res[i]
		}
	}

	if max.Equal(min) {
		return decimal.Zero, nil
	}

	return curr.Sub(min).Div(max.Sub(min)), nil
}

// Count determines the total amount of data needed for SRSI
// calculation.
func (srsi SRSI) Count() int {
	return srsi.rsi.length*2 - 1
}

// Stoch holds all the necessary information needed to calculate stochastic
// oscillator.
// The zero value is not usable.
type Stoch struct {
	// valid specifies whether Stoch paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewStoch validates provided configuration options and
// creates new Stoch indicator.
func NewStoch(length int) (Stoch, error) {
	stoch := Stoch{
		length: length,
	}

	if err := stoch.validate(); err != nil {
		return Stoch{}, err
	}

	return stoch, nil
}

// validate checks whether the indicator has valid configuration properties.
func (stoch *Stoch) validate() error {
	if stoch.length < 1 {
		return ErrInvalidLength
	}

	stoch.valid = true

	return nil
}

// Calc calculates Stoch from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/s/stochasticoscillator.asp.
func (stoch Stoch) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !stoch.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != stoch.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	low := dd[0]
	high := dd[0]

	for i := 0; i < len(dd); i++ {
		if dd[i].LessThan(low) {
			low = dd[i]
		}

		if dd[i].GreaterThan(high) {
			high = dd[i]
		}
	}

	dnm := high.Sub(low)
	if dnm.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return dd[len(dd)-1].Sub(low).Div(dnm).Mul(_hundred), nil
}

// Count determines the total amount of data points needed for Stoch
// calculation.
func (stoch Stoch) Count() int {
	return stoch.length
}

// WMA holds all the necessary information needed to calculate weighted
// moving average.
// The zero value is not usable.
type WMA struct {
	// valid specifies whether WMA paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewWMA validates provided configuration options and
// creates new WMA indicator.
func NewWMA(length int) (WMA, error) {
	wma := WMA{
		length: length,
	}

	if err := wma.validate(); err != nil {
		return WMA{}, err
	}

	return wma, nil
}

// validate checks whether the indicator has valid configuration properties.
func (wma *WMA) validate() error {
	if wma.length < 1 {
		return ErrInvalidLength
	}

	wma.valid = true

	return nil
}

// Calc calculates WMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/articles/technical/060401.asp.
func (wma WMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !wma.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != wma.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := decimal.Zero

	weight := decimal.NewFromInt(int64(wma.length * (wma.length + 1))).Div(decimal.NewFromInt(2))

	for i := 0; i < len(dd); i++ {
		res = res.Add(dd[i].Mul(decimal.NewFromInt(int64(i + 1)).Div(weight)))
	}

	return res, nil
}

// Count determines the total amount of data points needed for WMA
// calculation.
func (wma WMA) Count() int {
	return wma.length
}
