package postback

import (
	binomPostback "github.com/CLi-Ter/binomv2-postback"
)

// PostbackService предназначен для обслуживания постбеков и отправки их в правильный трекер.
type Service interface {
	// SendEvent - обновляет клик данными о событиях.
	SendEvent(clickID string, ev binomPostback.Event) error
	// SendPostback - отправляет постбек по клику.
	// Создает конверсию или обновляет клик в зависимости от содержания запроса.
	SendPostback(binomPostback.Request) error
	// BuildPostbackRequest - строит и возвращает запрос.
	BuildPostbackRequest(postback binomPostback.RequestBuilder) binomPostback.Request
}
