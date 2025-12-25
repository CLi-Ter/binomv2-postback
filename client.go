package binomv2postback

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/CLi-Ter/binomv2-postback/binom"
)

type Client interface {
	SendPostback(clickID string, status *string, payout *float64, events Events) error
	UpdatePayout(clickID string, payout float64) error
	SendEvents(clickID string, events Events) error
	SendEvent(clickID string, event Event) error
	AddEvent(clickID string, index uint8) error
	SubEvent(clickID string, index uint8) error
	SetupEvent(clickID string, index uint8) error
	ResetEvent(clickID string, index uint8) error
	DryRun()
}

type client struct {
	dryRun       bool
	clickBaseURL string // Базовый URL для клика в трекере https://binom.tracker/click
	apiKey       string // API-ключ от Binom
	updKey       string // UPDKey из настроек Binom

	httpClient *http.Client
}

// AddEvent добавляет к событию index единицу
func (cli *client) AddEvent(clickID string, index uint8) error {
	return cli.SendEvent(clickID, binom.AddEvent(int8(index), 1))
}

// SubEvent вычитает у события index единицу
func (cli *client) SubEvent(clickID string, index uint8) error {
	return cli.SendEvent(clickID, binom.AddEvent(int8(index), -1))
}

// SetupEvent устанавливает событие index в единицу
func (cli *client) SetupEvent(clickID string, index uint8) error {
	return cli.SendEvent(clickID, binom.Event(int8(index), 1))
}

// ResetEvent устанавливает событие index в ноль
func (cli *client) ResetEvent(clickID string, index uint8) error {
	return cli.SendEvent(clickID, binom.Event(int8(index), 0))
}

// NewClient создает новый клиент для Binom-трекера, у которого клик адрес расположен по clickBaseURL.
// apiKey - нужен для создания базового клика, т.к. он создается в Binom через API.
// updKey - нужен для обновления данных по клику (отправка событий), если он установлен в настройках Binom.
func NewClient(clickBaseURL string, apiKey string, updKey string) Client {
	return &client{
		clickBaseURL: clickBaseURL,
		apiKey:       apiKey,
		updKey:       updKey,

		httpClient: &http.Client{},
	}
}

func (cli *client) DryRun() {
	cli.dryRun = true
}

// sendClick отправляет GET запрос в binom на обработчик клика.
// Это может быть базовый клик, lp клик, клик по кампании
// событие (если клик уже существует) или же конверсия.
func (cli *client) sendClick(query string) error {
	// Создаем GET HTTP-запрос
	req, err := http.NewRequest(http.MethodGet, cli.clickBaseURL, nil)
	if err != nil {
		return err
	}
	// добавляем параметры, в зависимости от них Binom понимает, что мы присылаем
	req.URL.RawQuery = query
	if cli.dryRun {
		fmt.Println("dryRun req URL:", req.URL.String())
		return nil
	}
	// Отправляем запрос, ожидаем 200-ый ответ
	response, err := cli.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Получив ошибку, пытаемся прочесть содержимое ответа и вернуть его как ошибку
	if response.StatusCode != http.StatusOK {
		var body []byte
		_, err = response.Body.Read(body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return fmt.Errorf("failed to send request, status code: %d, response %s", response.StatusCode, string(body))
	}

	return nil
}

// SendEvents обновляет клик событиями (конверсия не генерируется)
func (cli *client) SendEvents(clickID string, events Events) error {
	q := make(url.Values)
	q.Add("upd_clickid", clickID)
	q.Add("upd_key", cli.updKey)

	return cli.sendClick(q.Encode() + "&" + events.URLParams())
}

// SendEvent отправляет (postback.AddEvent) или обновляет (postback.SetEvent)
// событие с номером 1 <= index <= 30.
func (cli *client) SendEvent(clickID string, event Event) error {
	events := Events{}
	if err := events.Set(event, false); err != nil {
		return err
	}

	return cli.SendEvents(clickID, events)
}

// SendPostback отправляет/обновляет конверсию с cnv_id=clickID.
// не обновляет статус конверсии, если status=nil
// не обнволяет выплату, если payout=nil
// во время конверсии можно добавить-заменить события через events
func (cli *client) SendPostback(clickID string, status *string, payout *float64, events Events) error {
	q := make(url.Values)
	q.Add("cnv_id", clickID)
	if status != nil {
		q.Add("cnv_status", *status)
	}
	if payout != nil {
		q.Add("payout", fmt.Sprintf("%f", *payout))
	}

	var output []string
	output = append(output, q.Encode())
	eventsParams := events.URLParams()
	if eventsParams != "" {
		output = append(output, eventsParams)
	}

	return cli.sendClick(strings.Join(output, "&"))
}

// UpdatePayout implements Client.
func (cli *client) UpdatePayout(clickID string, payout float64) error {
	return cli.SendPostback(clickID, nil, &payout, Events{})
}

// SendBaseClick отправляет базовый клик на компанию с ключем campaignKey.
// если установлен lpbcid=true, то так же устанавливает LPClick.
func (cli *client) SendBaseClick(campaignKey string, lpbcid bool) error {
	panic("not implemented. coming in v0.4")
}

// SetLPClick устанавливает клик по лендингу для клика clickID.
func (cli *client) SetLPClick(clickID string) error {
	panic("not implemented. coming in v0.4")
}

// SendClick производит клик по офферу.
func (cli *client) SendClick() error {
	panic("not implemented. coming in v0.4")
}
