package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/CLi-Ter/binomv2-postback/binom"
	"github.com/CLi-Ter/binomv2-postback/entity"
	ver2 "github.com/CLi-Ter/binomv2-postback/ver2"
)

type ClientOptions struct {
	SendEmptyUpdates bool
	DryRun           bool // Режим без отправки HTTP-запроса
}

type client struct {
	clickBaseURL string  // Базовый URL для клика в трекере https://binom.tracker/click
	apiKey       string  // API-ключ от Binom
	updKey       *string // UPDKey из настроек Binom

	opt ClientOptions

	log        Logger
	httpClient *http.Client
}

func (cli *client) SetLogger(log Logger) {
	cli.log = log
}

// AddEvent добавляет к событию index единицу
func (cli *client) AddEvent(clickID string, index uint8, opts ...sendClickOpt) error {
	return cli.SendEvent(clickID, binom.AddEvent(int8(index), 1), opts...)
}

// SubEvent вычитает у события index единицу
func (cli *client) SubEvent(clickID string, index uint8, opts ...sendClickOpt) error {
	return cli.SendEvent(clickID, binom.AddEvent(int8(index), -1), opts...)
}

// SetupEvent устанавливает событие index в единицу
func (cli *client) SetupEvent(clickID string, index uint8, opts ...sendClickOpt) error {
	return cli.SendEvent(clickID, binom.Event(int8(index), 1), opts...)
}

// ResetEvent устанавливает событие index в ноль
func (cli *client) ResetEvent(clickID string, index uint8, opts ...sendClickOpt) error {
	return cli.SendEvent(clickID, binom.Event(int8(index), 0), opts...)
}

// NewClient создает новый клиент для Binom-трекера, у которого клик адрес расположен по clickBaseURL.
// apiKey - нужен для создания базового клика, т.к. он создается в Binom через API.
// updKey - нужен для обновления данных по клику (отправка событий), если он установлен в настройках Binom.
func NewClient(clickBaseURL string, apiKey string, updKey string) Client {
	var uk *string
	if updKey != "" {
		uk = &updKey
	}
	return &client{
		clickBaseURL: clickBaseURL,
		apiKey:       apiKey,
		updKey:       uk,

		httpClient: &http.Client{},
	}
}

func (cli *client) DryRun() {
	cli.opt.DryRun = true
}

type clickReq struct {
	ctx  context.Context
	log  Logger
	body io.Reader

	method       string
	clickBaseURL string
	dryRun       bool

	pbLvl entity.PostbackLevel
}

// sendClick отправляет GET запрос в binom на обработчик клика.
// Это может быть базовый клик, lp клик, клик по кампании
// событие (если клик уже существует) или же конверсия.
func (cli *client) sendClick(query string, opt ...sendClickOpt) error {
	clkReq := &clickReq{
		method:       http.MethodGet,
		clickBaseURL: cli.clickBaseURL,
		dryRun:       cli.opt.DryRun,
		body:         nil,
		ctx:          nil,
		log:          cli.log,
	}
	// Применяем опции клика на запрос
	for _, f := range opt {
		if err := f(cli, clkReq); err != nil {
			return err
		}
	}

	// Проверка уровня отправки постбека для этого клика.
	// Если клик нельзя отправлять, то возвращаем nil
	if !clkReq.pbLvl.CanPostback() {
		return nil
	}
	// Если уровень отправки клика только в трекер, то добавляем опцию
	// disable_postback к запросу
	// TODO: this is all?
	if clkReq.pbLvl == entity.PB_LVL_NO_TS {
		query = query + "&disable_postback=1"
	}

	// Создаем GET HTTP-запрос
	req, err := http.NewRequest(clkReq.method, clkReq.clickBaseURL, clkReq.body)
	if err != nil {
		return err
	}
	// Если в опциях есть контекст, то добавляем его к запросу
	if clkReq.ctx != nil {
		req = req.WithContext(clkReq.ctx)
	}

	// добавляем параметры, в зависимости от них Binom понимает, что мы присылаем
	req.URL.RawQuery = query
	if clkReq.dryRun {
		fmt.Println("dryRun req URL:", req.URL.String())
		return nil
	}
	if clkReq.log != nil {
		clkReq.log.Infof("Send binom request: %v", req.URL.String())
	}

	// Отправляем запрос, ожидаем 200-ый ответ
	response, err := cli.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if clkReq.log != nil {
		clkReq.log.Debugf("Binom request: %v Response: %v", req, response)
	}

	return parsePostbackResponse(response)
}

// SendEvents обновляет клик событиями (конверсия не генерируется)
func (cli *client) SendEvents(clickID string, events entity.Events, opts ...sendClickOpt) error {
	eventParams := events.URLParams()
	// не посылать пустые события !!!
	// если опция SendEmptyUpdates не включена,
	// то производим проверку значений событий на пустоту перед отправкой.
	if !cli.opt.SendEmptyUpdates {
		if eventParams == "" {
			if cli.log != nil {
				cli.log.Debugf("SendEvents>cli.opt.SendEmptyUpdates: empty update")
			}
			return nil
		}
	}

	q := make(url.Values)
	q.Add("upd_clickid", clickID)
	if cli.updKey != nil {
		q.Add("upd_key", *cli.updKey)
	}

	return cli.sendClick(q.Encode()+"&"+eventParams, opts...)
}

// SendEvent отправляет (postback.AddEvent) или обновляет (postback.SetEvent)
// событие с номером 1 <= index <= 30.
func (cli *client) SendEvent(clickID string, event ver2.Event, opts ...sendClickOpt) error {
	events := entity.Events{}
	if err := events.Set(event, false); err != nil {
		return err
	}

	return cli.SendEvents(clickID, events, opts...)
}

func (cli *client) SendPostbackRequest(postback Request, opts ...sendClickOpt) error {
	// если это не конверсия, то отправляем через SendEvents, чтобы не триггерить postback в биноме
	if !postback.IsConversion() {
		return cli.SendEvents(postback.ClickID(), postback.Events(), opts...)
	}

	return cli.sendClick(postback.URLParam(), opts...)
}

// SendPostback отправляет/обновляет конверсию с cnv_id=clickID.
// не обновляет статус конверсии, если status=nil
// не обнволяет выплату, если payout=nil
// во время конверсии можно добавить-заменить события через events
func (cli *client) SendPostback(clickID string, status *string, payout *float64, events entity.Events, opts ...sendClickOpt) error {
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

	return cli.sendClick(strings.Join(output, "&"), opts...)
}

// UpdatePayout implements Client.
func (cli *client) UpdatePayout(clickID string, payout float64) error {
	return cli.SendPostback(clickID, nil, &payout, entity.Events{})
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

func parsePostbackResponse(resp *http.Response) error {
	var body []byte
	_, err := resp.Body.Read(body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Получив ошибку, пытаемся прочесть содержимое ответа и вернуть его как ошибку
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send request, status code: %d, response %s", resp.StatusCode, string(body))
	}

	// Binom сейчас возвращает 200 даже при ошибках
	// При ошибках внутри тела ответа status=fail
	return parsePostbackResponseBody(body)
}

func parsePostbackResponseBody(body []byte) error {
	var resp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if len(body) > 0 { // Проверяем, что тело не пустое
		if err := json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("unmarshal response body failed: %v", err)
		} else {
			if resp.Status == "fail" {
				return fmt.Errorf("postback request failed with error: %s", resp.Message)
			}
		}
	}

	return nil
}
