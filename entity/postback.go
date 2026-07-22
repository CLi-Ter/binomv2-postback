package entity

import (
	"strconv"

	binv2post "github.com/CLi-Ter/binomv2-postback"
	"github.com/CLi-Ter/binomv2-postback/click"
)

// Postback предназначен для работы со структурами постбеков
type Postback interface {
	ClickID() string
	Payout() string
	ConversionStatus() string
	ConversionStatus2() string
	Currency() string
	Events() binv2post.Events
	ToOffer() string
}

func NewPostback(clk click.Click, conv binv2post.Conversion) Postback {
	return &postback{
		Click:      clk,
		Conversion: conv,
	}
}

type postback struct {
	Click       click.Click
	EventValues binv2post.Events
	Conversion  binv2post.Conversion
	toOffer     *uint64

	DisablePostback bool
}

func (p *postback) ClickID() string {
	return p.Click.String()
}

func (p *postback) Events() binv2post.Events {
	return p.EventValues
}

func (p *postback) Payout() string {
	if p.Conversion == nil || p.Conversion.Payout() == nil || !p.Conversion.Payout().HasCurrency() {
		return ""
	}

	return strconv.FormatFloat(p.Conversion.Payout().Value(), 'g', -1, 64)
}

func (p *postback) ConversionStatus() string {
	return p.Conversion.Status()
}

func (p *postback) ConversionStatus2() string {
	return p.Conversion.Status2()
}

func (p *postback) Currency() string {
	if p.Conversion == nil || p.Conversion.Payout() == nil || !p.Conversion.Payout().HasCurrency() {
		return ""
	}

	return p.Conversion.Payout().Currency()
}

func (p *postback) ToOffer() string {
	if p.toOffer == nil {
		return ""
	}

	return strconv.FormatUint(*p.toOffer, 10)
}
