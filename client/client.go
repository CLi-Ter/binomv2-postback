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

// Event представляет собой событие в биноме. https://docs.binom.org/events-v2.php
// всего событий в BinomV2 от 1 до 30 далее X. Их значение можно обновлять SetEvent или складывать AddEvent.
// в URL события имеют вид eventX=INT или add_eventX=INT
type Event interface {
	bp.Event
}

type EventClient interface {
	// отправка события
	SendEvent(clickID string, event bp.Event, opts ...SendClickOpt) error
	SendEvents(clickID string, events bp.Events, opts ...SendClickOpt) error
	// работа с счетчиком события
	AddEvent(clickID string, index uint8, opts ...SendClickOpt) error
	SubEvent(clickID string, index uint8, opts ...SendClickOpt) error
	SetupEvent(clickID string, index uint8, opts ...SendClickOpt) error
	ResetEvent(clickID string, index uint8, opts ...SendClickOpt) error
}

type PostbackClient interface {
	SendPostbackRequest(postback Request, opts ...SendClickOpt) error
	SendPostback(clickID string, status *string, payout *float64, events bp.Events, opts ...SendClickOpt) error
}
