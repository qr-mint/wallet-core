package change_hero

import (
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

type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type GetMinAmountParams struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type GetMinAmountResponse struct {
	Result string `json:"result"`
}

func (c Client) GetMinAmount(params GetMinAmountParams) (*GetMinAmountResponse, error) {
	var response GetMinAmountResponse
	httpErr, err := c.httpClient.SendJsonPostRequest(http.JsonRequestData{
		Body: Request{
			Method: "getMinAmount",
			Params: params,
		},
		Data: http.RequestData{
			Path:    "/v2/",
			Headers: []http.HeaderData{{Name: "api-key", Value: c.params.ApiKey}},
		},
	}, &response)

	if err != nil {
		return nil, errors.Errorf("can not get change hero min amount error: %s", err)
	}
	if httpErr != nil {
		return nil, errors.Errorf("can not get change hero min amount error: %s", httpErr)
	}

	return &response, nil
}

type CreateTransactionParams struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	AddressTo string  `json:"address"`
	Amount    float64 `json:"amount"`
}

type CreateTransactionResponse struct {
	Result struct {
		Id           string `json:"id"`
		PayInAddress string `json:"payinAddress"`
	} `json:"result"`
}

func (c Client) CreateTransaction(params CreateTransactionParams) (*CreateTransactionResponse, error) {
	var response CreateTransactionResponse
	httpErr, err := c.httpClient.SendJsonPostRequest(http.JsonRequestData{
		Body: Request{
			Method: "createTransaction",
			Params: params,
		},
		Data: http.RequestData{
			Path:    "/v2/",
			Headers: []http.HeaderData{{Name: "api-key", Value: c.params.ApiKey}},
		},
	}, &response)

	if err != nil {
		return nil, errors.Errorf("can not create change hero transaction error: %s", err)
	}
	if httpErr != nil {
		return nil, errors.Errorf("can not create change hero transaction error: %s", httpErr)
	}

	return &response, nil
}

type GetExchangeAmountParams struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type GetExchangeAmountResponse struct {
	Result string `json:"result"`
}

func (c Client) GetExchangeAmount(params GetExchangeAmountParams) (*GetExchangeAmountResponse, error) {
	var response GetExchangeAmountResponse
	httpErr, err := c.httpClient.SendJsonPostRequest(http.JsonRequestData{
		Body: Request{
			Method: "getExchangeAmount",
			Params: params,
		},
		Data: http.RequestData{
			Path:    "/v2/",
			Headers: []http.HeaderData{{Name: "api-key", Value: c.params.ApiKey}},
		},
	}, &response)

	if err != nil {
		return nil, errors.Errorf("can not create change hero transaction error: %s", err)
	}
	if httpErr != nil {
		return nil, errors.Errorf("can not create change hero transaction error: %s", httpErr)
	}

	return &response, nil
}
