package binomv2postback

import (
	"strconv"
	"strings"
)

type Request interface {
	ClickID() string
	Payout() string
	ConversionStatus() string
	ConversionStatus2() string
	Currency() string
	Events() Events
	Params() []string
	URLParam() string
	String() string
	IsConversion() bool
	IsDisabledPostback() bool
	ToOffer() string
}

type request struct {
	clickID         string
	payout          *float64
	cnvStatus       *string
	cnvStatus2      *string
	currency        *string
	isCnv           bool
	events          Events
	disablePostback bool
	toOffer         *uint64
}

func (p *request) ClickID() string {
	return p.clickID
}

func (p *request) Payout() string {
	if p.payout == nil {
		return ""
	}

	return strconv.FormatFloat(*p.payout, 'g', -1, 64)
}

func (p *request) ConversionStatus() string {
	if p.cnvStatus == nil {
		return ""
	}

	return *p.cnvStatus
}

func (p *request) ConversionStatus2() string {
	if p.cnvStatus2 == nil {
		return ""
	}

	return *p.cnvStatus2
}

func (p *request) Currency() string {
	if p.currency == nil {
		return ""
	}

	return *p.currency
}

func (p *request) Events() Events {
	return p.events
}

func (p *request) IsConversion() bool {
	// Это конверсия - если это прописано явно в запросе, есть какой-то статус или payout.
	if p.isCnv || p.cnvStatus != nil || p.cnvStatus2 != nil || p.payout != nil {
		return true
	}

	return false
}

func (p *request) IsDisabledPostback() bool {
	return p.disablePostback
}

func (p *request) ToOffer() string {
	if p.toOffer == nil {
		return ""
	}

	return strconv.FormatUint(*p.toOffer, 10)
}

func (p *request) Params() []string {
	var output []string
	// чтобы была поддержка String как в бином, тут не добавляем cnv_id,
	// он добавляется в URLParam
	output = append(output, p.clickID)
	if p.payout != nil {
		output = append(output, "payout="+p.Payout())
	}
	if p.currency != nil {
		output = append(output, "cnv_currency="+p.Currency())
	}
	if p.cnvStatus != nil {
		output = append(output, "cnv_status="+p.ConversionStatus())
	}
	if p.cnvStatus2 != nil {
		output = append(output, "cnv_status2="+p.ConversionStatus2())
	}
	// добавляем события
	output = append(output, p.events.Params()...)

	// TODO: Следующие 2 параметра возможно будут перенесены в URLParam.
	// Мне пока не понятно поведение binom, если отправить в postbackManager строки,
	// которые вернет метод String()
	// ------------------------------>
	// записать клик на оффер N (N - порядковый номер оффера в пути).
	if p.toOffer != nil {
		output = append(output, "to_offer="+p.ToOffer())
	}
	// устанавливаем флаг не отсылать постбек, если включена опция
	if p.disablePostback {
		output = append(output, "disable_postback=1")
	}
	// <------------------------------

	return output
}

func (p *request) URLParam() string {
	params := p.Params()
	// добавляем cnv_id, т.к. в Params() он не устанавливается
	// если это конверсия, то ставим cnv_id, если
	params[0] = "cnv_id=" + params[0]

	return strings.Join(params, "&")
}

func (p *request) String() string {
	// для поддержки формата обновления Binom конверсий
	return strings.Join(p.Params(), ":")
}
