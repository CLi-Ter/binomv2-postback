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
}

type request struct {
	clickID    string
	payout     *float64
	cnvStatus  *string
	cnvStatus2 *string
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

func (p *request) Params() []string {
	var output []string
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
	for _, ev := range p.events {
		if ev == nil {
			continue
		}

		if ev.Index() > 0 {
			output = append(output, ev.URLParam())
		}
	}
	return output
}

func (p *request) URLParam() string {
	return strings.Join(p.Params(), "&")
}

func (p *request) String() string {
	return strings.Join(p.Params(), ":")
}
