package binomv2postback

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/CLi-Ter/binomv2-postback/binom"
)

type EventClient interface {
	// отправка события
	SendEvent(clickID string, event Event) error
	SendEvents(clickID string, events Events) error
	// работа с счетчиком события
	AddEvent(clickID string, index uint8) error
	SubEvent(clickID string, index uint8) error
	SetupEvent(clickID string, index uint8) error
	ResetEvent(clickID string, index uint8) error
}

type PostbackClient interface {
	SendPostbackRequest(postback Request, opts ...sendClickOpt) error
	SendPostback(clickID string, status *string, payout *float64, events Events) error
}

// Client это клиент для трекера Binom позволяющий работать с кликом.
type Client interface {
	EventClient
	PostbackClient
	DryRun()
	SetLogger(log Logger)
}

type client struct {
	dryRun       bool
	clickBaseURL string // Базовый URL для клика в трекере https://binom.tracker/click
	apiKey       string // API-ключ от Binom
	updKey       string // UPDKey из настроек Binom
	log          Logger

	httpClient *http.Client
}

func (cli *client) SetLogger(log Logger) {
	cli.log = log
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

type sendClickOpt func(cli *client, clkReq *clickReq) error

func OptWithClickBaseURL(clickBaseURL string) sendClickOpt {
	return func(cli *client, clkReq *clickReq) error {
		if clkReq != nil && clkReq.log != nil {
			clkReq.log.Debugf("Setup click request with clickBaseURL option: %s", clickBaseURL)
		}
		clkReq.clickBaseURL = clickBaseURL

		return nil
	}
}

func OptWithHost(host string) sendClickOpt {
	return func(cli *client, clkReq *clickReq) error {
		if clkReq != nil && clkReq.log != nil {
			clkReq.log.Debugf("setup click request with host option: %s", host)
		}
		url, err := url.Parse(clkReq.clickBaseURL)
		if err != nil {
			return err
		}
		url.Host = host

		clkReq.clickBaseURL = url.String()

		return nil
	}
}

func OptWithDryRun(dryRun bool) sendClickOpt {
	return func(cli *client, clkReq *clickReq) error {
		if clkReq != nil && clkReq.log != nil {
			clkReq.log.Debugf("setup click request with dryRun option: %b", dryRun)
		}
		clkReq.dryRun = dryRun

		return nil
	}
}

func OptDryRun() sendClickOpt {
	return OptWithDryRun(true)
}

func OptWithContext(ctx context.Context) sendClickOpt {
	return func(cli *client, clkReq *clickReq) error {
		if clkReq != nil && clkReq.log != nil {
			clkReq.log.Debugf("setup click request with context option: %v", ctx)
		}
		clkReq.ctx = ctx

		return nil
	}
}

type SendClickOptions []sendClickOpt

type clickReq struct {
	method       string
	clickBaseURL string
	dryRun       bool
	body         io.Reader
	ctx          context.Context
	log          Logger
}

// sendClick отправляет GET запрос в binom на обработчик клика.
// Это может быть базовый клик, lp клик, клик по кампании
// событие (если клик уже существует) или же конверсия.
func (cli *client) sendClick(query string, opt ...sendClickOpt) error {
	clkReq := &clickReq{
		method:       http.MethodGet,
		clickBaseURL: cli.clickBaseURL,
		dryRun:       cli.dryRun,
		body:         nil,
		ctx:          nil,
		log:          cli.log,
	}
	for _, f := range opt {
		if err := f(cli, clkReq); err != nil {
			return err
		}
	}
	// Создаем GET HTTP-запрос
	req, err := http.NewRequest(clkReq.method, clkReq.clickBaseURL, clkReq.body)
	if err != nil {
		return err
	}
	if clkReq.ctx != nil {
		req = req.WithContext(clkReq.ctx)
	}
	// добавляем параметры, в зависимости от них Binom понимает, что мы присылаем
	req.URL.RawQuery = query
	if clkReq.log != nil {
		clkReq.log.Infof("Send binom request: %v", req)
	}

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

	if clkReq.log != nil {
		clkReq.log.Infof("Binom request: %v Response: %v", req, response)
	}

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

func (cli *client) SendPostbackRequest(postback Request, opts ...sendClickOpt) error {
	return cli.sendClick(postback.URLParam())
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
	panic("not implemented. coming in v0.9")
}

// SetLPClick устанавливает клик по лендингу для клика clickID.
func (cli *client) SetLPClick(clickID string) error {
	panic("not implemented. coming in v0.9")
}

// SendClick производит клик по офферу.
func (cli *client) SendClick() error {
	panic("not implemented. coming in v0.9")
}
