package client

import (
	"context"
	"net/url"

	"github.com/CLi-Ter/binomv2-postback/entity"
)

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

func OptWithPostbackLevel(lvl entity.PostbackLevel) sendClickOpt {
	return func(cli *client, clkReq *clickReq) error {
		clkReq.pbLvl = lvl

		return nil
	}
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

func OptWithLogger(logger Logger) sendClickOpt {
	return func(cli *client, clkReq *clickReq) error {
		clkReq.log = logger

		return nil
	}
}

type SendClickOptions []sendClickOpt
