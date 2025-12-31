package binomv2postback

import "strings"

type RequestBuilder interface {
	Request() Request
	WithPayout(payout float64) RequestBuilder
	WithEvents(events Events) RequestBuilder
	WithStatus(cnvStatus string, cnvStatus2 ...string) RequestBuilder
}

type requestBuilder struct {
	req *request
}

// Request implements RequestBuilder.
func (r *requestBuilder) Request() Request {
	return r.req
}

// WithEvents implements RequestBuilder.
func (r *requestBuilder) WithEvents(events Events) RequestBuilder {
	r.req.events = events
	return r
}

// WithPayout implements RequestBuilder.
func (r *requestBuilder) WithPayout(payout float64) RequestBuilder {
	r.req.payout = &payout
	return r
}

// WithStatus implements RequestBuilder.
func (r *requestBuilder) WithStatus(cnvStatus string, cnvStatus2 ...string) RequestBuilder {
	r.req.cnvStatus = &cnvStatus
	if len(cnvStatus2) > 0 {
		cnvStatuses := strings.Join(cnvStatus2, "_")
		r.req.cnvStatus2 = &cnvStatuses
	}

	return r
}

func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{}
}
