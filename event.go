package binomv2postback

import (
	"fmt"
	"strconv"
	"strings"
)

// Event представляет собой событие в биноме. https://docs.binom.org/events-v2.php
// всего событий в BinomV2 от 1 до 30 далее X. Их значение можно обновлять SetEvent или складывать AddEvent.
// в URL события имеют вид eventX=INT или add_eventX=INT
type Event interface {
	Type() string                // тип значения add_event или event
	Value() int                  // значение события
	Name(index uint8) string     // имя URL-аргумента
	URLParam(index uint8) string // форматирование значения в виде URL-аргумента
}

// AddEvent - значение счетчика
type AddEvent int

func (ev AddEvent) Type() string {
	return "add_event"
}
func (ev AddEvent) Value() int {
	return int(ev)
}
func (ev AddEvent) Name(index uint8) string {
	return ev.Type() + strconv.Itoa(int(index))
}
func (ev AddEvent) URLParam(index uint8) string {
	return ev.Name(index) + "=" + strconv.Itoa(ev.Value())
}

// SetEvent - обновление значения
type SetEvent int

func (ev SetEvent) Type() string {
	return "event"
}
func (ev SetEvent) Value() int {
	return int(ev)
}
func (ev SetEvent) Name(index uint8) string {
	return ev.Type() + strconv.Itoa(int(index))
}
func (ev SetEvent) URLParam(index uint8) string {
	return ev.Name(index) + "=" + strconv.Itoa(ev.Value())
}

// События в трекере (всего их 30)
type Events [30]Event

// URLParams преобразует массив Events в строку URL-аргументов
func (e *Events) URLParams() string {
	out := []string{}
	for i, v := range *e {
		if v == nil {
			continue
		}
		out = append(out, v.URLParam(uint8(i)))
	}

	return strings.Join(out, "&")
}

// Set проверяет наличие события в массиве и устанавливает конкретное событие index=X
// если force=true, либо выбрасывает ошибку (TODO: конкретная ошибка)
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
