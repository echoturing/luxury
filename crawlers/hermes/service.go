package hermes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var defaultHTTPClient = http.Client{
	Jar:     &cookiejar.Jar{},
	Timeout: time.Second * 20,
}

func CrawlGoods(ctx context.Context, query string) (*ProductResponse, error) {
	req, _ := http.NewRequest("GET", searchBaseURL+url.QueryEscape(query), nil)
	req = req.WithContext(ctx)
	response, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w|query:%s", err, query)
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w|query:%s", err, query)
	}
	productResponse := ProductResponse{}
	err = json.Unmarshal(responseBytes, &productResponse)
	if err != nil {
		return nil, fmt.Errorf("%w|resp:%s", err, string(responseBytes))
	}
	return &productResponse, nil
}
