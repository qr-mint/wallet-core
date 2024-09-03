package simple_swap

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/golib4/http-client/http"
)

type Params struct {
	ApiKey string
}

type Client struct {
	httpClient *http.Client
	params     Params
}

func NewClient(httpClient *http.Client, params Params) *Client {
	return &Client{httpClient: httpClient, params: params}
}

type CreateExchangeRequest struct {
	Fixed        bool    `json:"fixed"`
	CurrencyFrom string  `json:"currency_from"`
	CurrencyTo   string  `json:"currency_to"`
	Amount       float64 `json:"amount"`
	AddressTo    string  `json:"address_to"`
}

type CreateExchangeResponse struct {
	RedirectUrl string `json:"redirect_url"`
}

func (c Client) CreateExchange(request CreateExchangeRequest) (*CreateExchangeResponse, error) {
	var response CreateExchangeResponse
	httpErr, err := c.httpClient.SendJsonPostRequest(http.JsonRequestData{
		Body: request,
		Data: http.RequestData{
			Path: fmt.Sprintf("/create_exchange?apiKey=%s", c.params.ApiKey),
		},
	}, &response)

	if err != nil {
		return nil, errors.Errorf("can not create simple swap exchange error: %s", err)
	}
	if httpErr != nil {
		return nil, errors.Errorf("can not create simple swap exchange error: %s", httpErr)
	}

	return &response, nil
}

type GetRangesRequest struct {
	Fixed        bool   `schema:"fixed"`
	CurrencyFrom string `schema:"currency_from"`
	CurrencyTo   string `schema:"currency_to"`
	ApiKey       string `schema:"api_key"`
}

type GetRangesResponse struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

func (c Client) GetRanges(request GetRangesRequest) (*GetRangesResponse, error) {
	request.ApiKey = c.params.ApiKey
	var response GetRangesResponse
	httpErr, err := c.httpClient.SendGetRequest(http.SchemaRequestData{
		Body: request,
		Data: http.RequestData{
			Path: "/get_ranges",
		},
	}, &response)

	if err != nil {
		return nil, errors.Errorf("can not get simple swap range error: %s", err)
	}
	if httpErr != nil {
		return nil, errors.Errorf("can not get simple swap range error: %s", httpErr)
	}

	return &response, nil
}
