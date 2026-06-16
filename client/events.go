package client

import (
	"fmt"
	"strings"

	ver2 "github.com/CLi-Ter/binomv2-postback"
)

// События в трекере (всего их 30)
type Events [30]ver2.Event

// Params возвращает все события как параметры
func (e *Events) Params() []string {
	out := []string{}
	for _, v := range e {
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
func (e *Events) Set(ev ver2.Event, force bool) error {
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
