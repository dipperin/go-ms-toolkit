package msproxy

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/json"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type RestyRequester struct {
	url    string
	client *resty.Client
}

func (r *RestyRequester) Post(api string, req interface{}, respData interface{}) error {
	resp, err := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(json.StringifyJsonToBytes(req)).
		Post(r.url + api)

	if err != nil {
		return err
	}

	if err = json.ParseJsonFromBytes(resp.Body(), respData); err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusOK, http.StatusBadRequest, http.StatusForbidden:
		return nil
	default:
		return fmt.Errorf("request api: %s get resp status code: %d, raw resp: %s", api, resp.StatusCode(), string(resp.Body()))
	}
}

func NewRestyRequester(url string) *RestyRequester {
	r := &RestyRequester{url: url}
	r.client = resty.New()
	return r
}
