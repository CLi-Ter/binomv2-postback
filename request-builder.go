package binomv2postback

import "strings"

// RequestBuilder allows you to construct Request interface
type RequestBuilder interface {
	Request(clickID string) Request
	WithPayout(payout float64) RequestBuilder
	WithEvents(events Events) RequestBuilder
	WithStatus(cnvStatus string, cnvStatus2 ...string) RequestBuilder
	WithPostbackMode(mode string) RequestBuilder
	DropStatus(keepPrimary bool) RequestBuilder
	DropConversion() RequestBuilder
	ClickID() string
	Mode() string
}

func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{
		req: &request{},
	}
}

func NewRequestBuilderWithClickID(clickID string) RequestBuilder {
	return &requestBuilder{
		req: &request{
			clickID: clickID,
		},
	}
}

type requestBuilder struct {
	req  *request
	mode string
}

// Return request clickID
func (r *requestBuilder) ClickID() string {
	return r.req.clickID
}

// Mode returns string name of postback mode
func (r *requestBuilder) Mode() string {
	return r.mode
}

// Request method create a copy of builder and apply clickID to it
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

// WithPostbackMode set information to builder about mode,
// that can be requested later with Mode() method when processing it
func (r *requestBuilder) WithPostbackMode(mode string) RequestBuilder {
	r.mode = mode
	return r
}

// DisablePostback turns off futured S2S postback for this postback
func (r *requestBuilder) DisablePostback() RequestBuilder {
	r.req.disablePostback = true
	return r
}

// WithCurrency setup conversion currency
func (r *requestBuilder) WithCurrency(currency string) RequestBuilder {
	// TODO: validate currency
	r.req.currency = &currency

	return r
}

// WithToOffer setup offer N from Path to set click on
func (r *requestBuilder) WithToOffer(toOffer uint64) RequestBuilder {
	r.req.toOffer = &toOffer

	return r
}

// DropStatus clear request data about conversion status.
// if keepPrimary is true it clear only secondary conversion status
func (r *requestBuilder) DropStatus(keepPrimary bool) RequestBuilder {
	if !keepPrimary {
		r.req.cnvStatus = nil
	}
	r.req.cnvStatus2 = nil

	return r
}

// DropConversion clear all request data about converion.
// Data contains conversion statuses, payout and isCnv flag.
// After that action postback will contains only clickid and events data
func (r *requestBuilder) DropConversion() RequestBuilder {
	r.DropStatus(false)
	r.req.payout = nil
	r.req.isCnv = false
	r.req.currency = nil
	r.req.toOffer = nil

	return r
}
