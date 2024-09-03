package finch_pay

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

type GetLimitsRequest struct {
	Currency string `schema:"currency"`
}

type GetLimitsResponse struct {
	MinAmount string `json:"min_amount"`
	MaxAmount string `json:"max_amount"`
}

func (c Client) GetLimits(request GetLimitsRequest) (*GetLimitsResponse, error) {
	var response GetLimitsResponse
	httpErr, err := c.httpClient.SendGetRequest(http.SchemaRequestData{
		Body: request,
		Data: http.RequestData{
			Path:    "/v1/currencies/limits",
			Headers: []http.HeaderData{{Name: "x-api-key", Value: c.params.ApiKey}},
		},
	}, &response)

	if err != nil {
		return nil, errors.Errorf("can not get finch pay limits error: %s", err)
	}
	if httpErr != nil {
		return nil, errors.Errorf("can not get finch pay limits error: %s", httpErr)
	}

	return &response, nil
}
