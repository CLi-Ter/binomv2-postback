package binomv2postback

import "strings"

// RequestBuilder allows you to construct Request interface
type RequestBuilder interface {
	Request(clickID string) Request
	WithPayout(payout float64) RequestBuilder
	WithEvents(events Events) RequestBuilder
	WithStatus(cnvStatus string, cnvStatus2 ...string) RequestBuilder
	ClickID() string
}

func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{
		req: &request{},
	}
}

func NewRequestBuilderWithClickID(clickID string) RequestBuilder {
	return &requestBuilder{
		req: &request{},
	}
}

type requestBuilder struct {
	req *request
}

func (r *requestBuilder) ClickID() string {
	return r.req.clickID
}

// Request method create a copy of builder and apply clickID to it.
func (r *requestBuilder) Request(clickID string) Request {
	req := *r.req
	req.clickID = clickID

	return &req
}

// WithEvents add click events to builded Request
func (r *requestBuilder) WithEvents(events Events) RequestBuilder {
	r.req.events = events
	return r
}

// WithPayout add conversion payout to builded Request
func (r *requestBuilder) WithPayout(payout float64) RequestBuilder {
	r.req.payout = &payout
	return r
}

// WithStatus add conversion status to builded Request
func (r *requestBuilder) WithStatus(cnvStatus string, cnvStatus2 ...string) RequestBuilder {
	r.req.cnvStatus = &cnvStatus
	if len(cnvStatus2) > 0 {
		cnvStatuses := strings.Join(cnvStatus2, "_")
		r.req.cnvStatus2 = &cnvStatuses
	}

	return r
}
