package entity

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

type conversion struct {
	payout  Payout
	status  *string
	status2 *string
	toOffer *string
}

func (conv *conversion) Payout() Payout {
	return conv.payout
}

func (conv *conversion) Status() string {
	if conv.status == nil {
		return ""
	}

	return *conv.status
}

func (conv *conversion) Status2() string {
	if conv.status2 == nil {
		return ""
	}

	return *conv.status2
}

func (conv *conversion) HasStatus() bool {
	return conv.status != nil
}

func (conv *conversion) HasStatus2() bool {
	return conv.status2 != nil
}
