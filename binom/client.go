package binom

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/CLi-Ter/binomv2-postback/binom/postback"
)

type Client interface {
	SendPostback(clickID string, status *string, payout *float64, events postback.Events) error
	SendEvents(clickID string, events postback.Events) error
	SendEvent(clickID string, index uint8, event postback.Event) error
}

type client struct {
	clickBaseURL string
	apiKey       string
	updKey       string

	httpClient *http.Client
}

func NewClient(clickBaseURL string, apiKey string, updKey string) Client {
	return &client{
		clickBaseURL: clickBaseURL,
		apiKey:       apiKey,
		updKey:       updKey,

		httpClient: &http.Client{},
	}
}

func (cli *client) sendClick(query string) error {
	req, err := http.NewRequest(http.MethodGet, cli.clickBaseURL, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = query

	response, err := cli.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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

func (cli *client) SendEvents(clickID string, events postback.Events) error {
	var q url.Values
	q.Add("upd_clickid", clickID)
	q.Add("upd_key", cli.updKey)

	return cli.sendClick(q.Encode() + "&" + events.URLParams())
}

func (cli *client) SendEvent(clickID string, index uint8, event postback.Event) error {
	events := postback.Events{}
	if err := events.Set(index, event, false); err != nil {
		return err
	}

	return cli.SendEvents(clickID, events)
}

func (cli *client) SendPostback(clickID string, status *string, payout *float64, events postback.Events) error {
	var q url.Values
	q.Add("cnv_id", clickID)
	if status != nil {
		q.Add("cnv_status", *status)
	}
	if payout != nil {
		q.Add("payout", fmt.Sprintf("%f", *payout))
	}

	return cli.sendClick(q.Encode() + "&" + events.URLParams())
}

func (cli *client) SendBaseClick(campaignKey string, lpbcid bool) error {
	return nil
}

func (cli *client) SetLPClick(clickId string) error {
	return nil
}

func (cli *client) SendClick() error {
	return nil
}
