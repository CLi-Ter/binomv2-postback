package postback

import (
	ver2 "github.com/CLi-Ter/binomv2-postback"
	bpClient "github.com/CLi-Ter/binomv2-postback/client"
)

// PostbackService предназначен для обслуживания постбеков и отправки их в правильный трекер.
type Service interface {
	// SendEvent - обновляет клик данными о событиях.
	SendEvent(clickID string, ev ver2.Event) error
	// SendPostback - отправляет постбек по клику.
	// Создает конверсию или обновляет клик в зависимости от содержания запроса.
	SendPostback(bpClient.Request) error
	// BuildPostbackRequest - строит и возвращает запрос.
	BuildPostbackRequest(postback bpClient.RequestBuilder) bpClient.Request
}
