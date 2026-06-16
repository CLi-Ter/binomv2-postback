package entity

const (
	PAYOUT_CURRENCY_GEL = "gel"
	PAYOUT_CURRENCY_EUR = "eur"
	PAYOUT_CURRENCY_USD = "usd"
	PAYOUT_CURRENCY_RUB = "rub"
)

type Payout interface {
	HasValue() bool
	HasCurrency() bool
	Value() float64
	Currency() string
}

type payout struct {
	val      *float64
	currency *string
}

// NewEmptyPayout создает пустой Payout
func NewEmptyPayout() Payout {
	return &payout{}
}

// NewCurrencyPayout создает Payout для валюты cur.
// Значение val может быть равным nil,
// тогда такой Payout должен обновлять лишь валюту выплаты.
func NewCurrencyPayout(cur string, val *float64) Payout {
	return &payout{
		val:      val,
		currency: &cur,
	}
}

func NewGelPayout(val float64) Payout {
	cur := PAYOUT_CURRENCY_GEL
	return &payout{
		val:      &val,
		currency: &cur,
	}
}

func NewEurPayout(val float64) Payout {
	cur := PAYOUT_CURRENCY_EUR
	return &payout{
		val:      &val,
		currency: &cur,
	}
}

func NewUsdPayout(val float64) Payout {
	return &payout{
		val: &val,
	}
}

func NewRubPayout(val float64) Payout {
	cur := PAYOUT_CURRENCY_RUB
	return &payout{
		val:      &val,
		currency: &cur,
	}
}

func (p *payout) HasValue() bool {
	return p.val != nil
}

func (p *payout) HasCurrency() bool {
	return p.currency != nil
}

func (p *payout) Value() float64 {
	if p.HasValue() {
		return *p.val
	}

	return 0
}

func (p *payout) Currency() string {
	if p.HasCurrency() {
		return *p.currency
	}

	return ""
}
