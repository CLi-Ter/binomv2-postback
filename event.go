package binomv2postback

import (
	"fmt"
	"strings"
)

// Event представляет собой событие в биноме. https://docs.binom.org/events-v2.php
// всего событий в BinomV2 от 1 до 30 далее X. Их значение можно обновлять SetEvent или складывать AddEvent.
// в URL события имеют вид eventX=INT или add_eventX=INT
type Event interface {
	Type() string     // тип значения add_event или event
	Value() int64     // значение события
	Index() int8      // номер события в трекере
	Name() string     // имя URL-аргумента
	URLParam() string // форматирование значения в виде URL-аргумента
}

// События в трекере (всего их 30)
type Events [30]Event

// Params возвращает все события как параметры
func (e *Events) Params() []string {
	out := []string{}
	for _, v := range *e {
		if v == nil {
			continue
		}
		out = append(out, v.URLParam())
	}

	return out
}

// String преобразует массив Events в строку
func (e *Events) String() string {
	return strings.Join(e.Params(), ":")
}

// URLParams преобразует массив Events в строку URL-аргументов
func (e *Events) URLParams() string {
	return strings.Join(e.Params(), "&")
}

// Set проверяет наличие события в массиве и устанавливает конкретное событие index=X
// если force=true, либо выбрасывает ошибку (TODO: конкретная ошибка)
func (e *Events) Set(ev Event, force bool) error {
	index := ev.Index()
	if int(index) > cap(e) {
		return fmt.Errorf("event index out of range. Max: %d", cap(e))
	}
	if v := e[index]; v != nil && !force {
		return fmt.Errorf("event %d already set %v", index, v)
	}
	e[index] = ev

	return nil
}
