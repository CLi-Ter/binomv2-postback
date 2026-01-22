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
	Events() Events
	Params() []string
	URLParam() string
	String() string
	IsConversion() bool
}

type request struct {
	clickID    string
	payout     *float64
	cnvStatus  *string
	cnvStatus2 *string
	isCnv      bool
	events     Events
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

func (p *request) Params() []string {
	var output []string
	// чтобы была поддержка String как в бином, тут не добавляем cnv_id,
	// он добавляется в URLParam
	output = append(output, p.clickID)
	if p.payout != nil {
		output = append(output, "payout="+p.Payout())
	}
	if p.cnvStatus != nil {
		output = append(output, "cnv_status="+(*p.cnvStatus))
	}
	if p.cnvStatus2 != nil {
		output = append(output, "cnv_status2="+(*p.cnvStatus2))
	}
	output = append(output, p.events.Params()...)

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
