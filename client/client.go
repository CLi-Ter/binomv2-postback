package client

import (
	bp "github.com/CLi-Ter/binomv2-postback"
	"github.com/CLi-Ter/binomv2-postback/entity"
)

// Client это клиент для трекера Binom позволяющий работать с кликом.
type Client interface {
	EventClient
	PostbackClient
	DryRun()
	SetLogger(log Logger)
}

type Request interface {
	entity.Postback

	Params() []string
	URLParam() string
	String() string
	IsConversion() bool
	IsDisabledPostback() bool

	SendClickOptions() SendClickOptions
}

type EventClient interface {
	// отправка события
	SendEvent(clickID string, event bp.Event, opts ...sendClickOpt) error
	SendEvents(clickID string, events entity.Events, opts ...sendClickOpt) error
	// работа с счетчиком события
	AddEvent(clickID string, index uint8, opts ...sendClickOpt) error
	SubEvent(clickID string, index uint8, opts ...sendClickOpt) error
	SetupEvent(clickID string, index uint8, opts ...sendClickOpt) error
	ResetEvent(clickID string, index uint8, opts ...sendClickOpt) error
}

type PostbackClient interface {
	SendPostbackRequest(postback Request, opts ...sendClickOpt) error
	SendPostback(clickID string, status *string, payout *float64, events entity.Events, opts ...sendClickOpt) error
}
