package entity

import (
	binv2post "github.com/CLi-Ter/binomv2-postback"
)

type conversion struct {
	payout  binv2post.Payout
	status  *string
	status2 *string
	toOffer *string
}

func (conv *conversion) Payout() binv2post.Payout {
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
