package binomv2postback

const (
	PAYOUT_CURRENCY_GEL = "gel"
	PAYOUT_CURRENCY_EUR = "eur"
	PAYOUT_CURRENCY_USD = "usd"
	PAYOUT_CURRENCY_RUB = "rub"
)

// Event представляет собой событие в биноме. https://docs.binom.org/events-v2.php
// всего событий в BinomV2 от 1 до 30 далее X. Их значение можно обновлять SetEvent или складывать AddEvent.
// в URL события имеют вид eventX=INT или add_eventX=INT
type Event interface {
	Type() string     // тип значения add_event или event
	Value() int64     // значение события
	Index() int8      // номер события в трекере
	Name() string     // имя URL-аргумента
	URLParam() string // форматирование значения в виде URL-аргумента
}

// Conversion хранит данные по конверсии.
// Payout - выплата, содержит сумму Value и валюту Currency (по-умолчанию - USD)
type Conversion interface {
	Payout() Payout
	Status() string
	Status2() string
	HasStatus() bool
	HasStatus2() bool
	// ToOffer() string ??
}

type Payout interface {
	HasValue() bool
	HasCurrency() bool
	Value() float64
	Currency() string
}

// Postback предназначен для работы со структурами постбеков
type Postback interface {
	ClickID() string
	Payout() string
	ConversionStatus() string
	ConversionStatus2() string
	Currency() string
	Events() []Event
	ToOffer() string
}
