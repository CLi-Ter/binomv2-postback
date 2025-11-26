package postback

import (
	"fmt"
	"strings"
)

type Postback struct {
	ID     string  `json:"cnv_id"`
	Payout float64 `json:"payout"`
	Events Events  `json:"events"`
	Status *string `json:"status,omitempty"`
}

type AddEvent int

func (ev AddEvent) Param(index uint8) string {
	return fmt.Sprintf("add_event%d=%d", index, ev)
}

type SetEvent int

func (ev SetEvent) Param(index uint8) string {
	return fmt.Sprintf("event%d=%d", index, ev)
}

type Event interface {
	Param(index uint8) string
}

type Events [30]Event

func (e *Events) URLParams() string {
	out := []string{}
	for i, v := range *e {
		if v == nil {
			continue
		}
		out = append(out, v.Param(uint8(i)))
	}

	return strings.Join(out, "&")
}

func (e *Events) Set(index uint8, ev Event, force bool) error {
	if int(index) > cap(e) {
		return fmt.Errorf("event index out of range. Max: %d", cap(e))
	}
	if v := e[index]; v != nil && !force {
		return fmt.Errorf("event %d already set %v", index, v)
	}
	e[index] = ev

	return nil
}

// func NewEvents() Events {
// 	return &events{}
// }

// type Events interface {
// 	URLParams() string
// }

type Service interface {
	Postback(cnv_id string, payout float64, events Events) error
}

type service struct {
	// binomClient BinomClient
	// logger      logger.Logger
}

func (s *service) Postback(cnv_id string, payout float64, events Events) error {

	return nil
}
