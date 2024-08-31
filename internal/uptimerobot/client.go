package uptimerobot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	apiKey string
}

type apiRequester interface {
	makeApiRequest(ctx context.Context, methodName string, params map[string]string) ([]byte, error)
}

func NewClient(apiKey string) Client {
	return Client{
		apiKey: apiKey,
	}
}

func (client Client) makeApiRequest(ctx context.Context, methodName string, params map[string]string) ([]byte, error) {
	url := fmt.Sprintf("https://api.uptimerobot.com/v2/%s", methodName)
	payloadString := fmt.Sprintf("api_key=%s&format=json", client.apiKey)
	if len(params) > 0 {
		for key, value := range params {
			payloadString = fmt.Sprintf("%s&%s=%s", payloadString, key, value)
		}
	}

	payload := strings.NewReader(payloadString)

	request, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		return []byte{}, err
	}

	request.Header.Add("cache-control", "no-cache")
	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return bodyBytes, nil
}

type APIResponse interface {
	GetStat() string
}

func request[Response APIResponse](ctx context.Context, apiCall string, requester apiRequester, paramBuilder func() (map[string]string, error)) (Response, error) {
	params, err := paramBuilder()
	if err != nil {
		return *new(Response), err
	}

	responseBytes, err := requester.makeApiRequest(ctx, apiCall, params)
	if err != nil {
		return *new(Response), err
	}

	var response Response
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return *new(Response), err
	}

	stat := response.GetStat()
	if stat == "fail" {
		return *new(Response), errors.New(string(responseBytes))
	}

	return response, nil
}

func (client Client) GetAccountDetails(ctx context.Context) (GetAccountDetailsResponse, error) {
	response, err := request[GetAccountDetailsResponse](ctx, "getAccountDetails", client, func() (map[string]string, error) {
		return map[string]string{}, nil
	})

	return response, err
}

func (client Client) DeleteAlertContact(ctx context.Context, id string) (DeleteAlertContactResponse, error) {
	response, err := request[DeleteAlertContactResponse](ctx, "deleteAlertContact", client, func() (map[string]string, error) {
		return map[string]string{
			"id": id,
		}, nil
	})

	return response, err
}

func (client Client) EditAlertContact(ctx context.Context, id string, value string, friendlyName string) (EditAlertContactResponse, error) {
	response, err := request[EditAlertContactResponse](ctx, "editAlertContact", client, func() (map[string]string, error) {
		params := map[string]string{
			"id":    id,
			"value": value,
		}

		if friendlyName != "" {
			params["friendly_name"] = friendlyName
		}

		return params, nil
	})

	return response, err
}

func (client Client) GetAlertContacts(ctx context.Context, alertContactIds []string) (GetAlertContactResponse, error) {
	response, err := request[GetAlertContactResponse](ctx, "getAlertContacts", client, func() (map[string]string, error) {
		params := map[string]string{}
		for _, id := range alertContactIds {
			if params["alert_contacts"] == "" {
				params["alert_contacts"] = fmt.Sprint(id)
			} else {
				params["alert_contacts"] = fmt.Sprintf("%s-%s", params["alert_contacts"], id)
			}
		}

		return params, nil
	})

	return response, err
}

func (client Client) NewAlertContact(ctx context.Context, alertType string, value string, friendlyName string) (NewAlertContactResponse, error) {
	response, err := request[NewAlertContactResponse](ctx, "newAlertContact", client, func() (map[string]string, error) {
		params := map[string]string{
			"type":  alertType,
			"value": value,
		}

		if friendlyName != "" {
			params["friendly_name"] = friendlyName
		}

		return params, nil
	})

	return response, err
}
