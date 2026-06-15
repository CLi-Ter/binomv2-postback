package entity

import (
	"strconv"

	"github.com/CLi-Ter/binomv2-postback/click"
)

// Postback предназначен для работы со структурами постбеков
type Postback interface {
	ClickID() string
	Payout() string
	ConversionStatus() string
	ConversionStatus2() string
	Currency() string
	Events() Events
	ToOffer() string
}

func NewPostback(clk click.Click) Postback {
	return &postback{
		Click: clk,
	}
}

type postback struct {
	Click       click.Click
	EventValues Events
	Values      struct {
		Payout     *float64
		Currency   *string
		CnvStatus  *string
		CnvStatus2 *string
		ToOffer    *uint64
	}

	DisablePostback bool
}

func (p *postback) ClickID() string {
	return p.Click.String()
}

func (p *postback) Events() Events {
	return p.EventValues
}

func (p *postback) Payout() string {
	if p.Values.Payout == nil {
		return ""
	}

	return strconv.FormatFloat(*p.Values.Payout, 'g', -1, 64)
}

func (p *postback) ConversionStatus() string {
	if p.Values.CnvStatus == nil {
		return ""
	}

	return *p.Values.CnvStatus
}

func (p *postback) ConversionStatus2() string {
	if p.Values.CnvStatus2 == nil {
		return ""
	}

	return *p.Values.CnvStatus2
}

func (p *postback) Currency() string {
	if p.Values.Currency == nil {
		return ""
	}

	return *p.Values.Currency
}

func (p *postback) ToOffer() string {
	if p.Values.ToOffer == nil {
		return ""
	}

	return strconv.FormatUint(*p.Values.ToOffer, 10)
}
