package binomv2postback

import (
	"strconv"
	"strings"
)

type Postback interface {
	ClickID() string
	Payout() string
	ConversionStatus() string
	ConversionStatus2() string
	Events() Events
	Params() []string
	String() string
}

type postback struct {
	clickID    string
	payout     *float64
	cnvStatus  *string
	cnvStatus2 *string
	events     Events
}

func (p *postback) ClickID() string {
	return p.clickID
}

func (p *postback) Payout() string {
	if p.payout == nil {
		return ""
	}

	return strconv.FormatFloat(*p.payout, 'g', -1, 64)
}

func (p *postback) ConversionStatus() string {
	if p.cnvStatus == nil {
		return ""
	}

	return *p.cnvStatus
}

func (p *postback) ConversionStatus2() string {
	if p.cnvStatus2 == nil {
		return ""
	}

	return *p.cnvStatus2
}

func (p *postback) Events() Events {
	return p.events
}

func (p *postback) URLParams() string {
	return ""
}

func (p *postback) Params() []string {
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

func (p *postback) String() string {
	return strings.Join(p.Params(), ":")
}
