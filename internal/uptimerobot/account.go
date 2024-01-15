package uptimerobot

import (
	"context"
	"encoding/json"
)

type Account struct {
	Email           string
	MonitorLimit    int
	MonitorInterval int
	UpMonitors      int
	DownMonitors    int
	PausedMonitors  int
}

func NewAccountFromResponse(response accountResponse) Account {
	return Account{
		Email:           response.Account.Email,
		MonitorLimit:    response.Account.MonitorLimit,
		MonitorInterval: response.Account.MonitorInterval,
		UpMonitors:      response.Account.UpMonitors,
		DownMonitors:    response.Account.DownMonitors,
		PausedMonitors:  response.Account.PausedMonitors,
	}
}

type accountResponse struct {
	Stat    string `json:"stat"`
	Account struct {
		Email           string `json:"email"`
		MonitorLimit    int    `json:"monitor_limit"`
		MonitorInterval int    `json:"monitor_interval"`
		UpMonitors      int    `json:"up_monitors"`
		DownMonitors    int    `json:"down_monitors"`
		PausedMonitors  int    `json:"paused_monitors"`
	} `json:"account"`
}

func GetAccountDetails(ctx context.Context, client Client) (Account, error) {
	apiResponse, err := makeApiRequest(ctx, client, "getAccountDetails", map[string]string{})
	if err != nil {
		return Account{}, err
	}

	var accountResponse accountResponse
	err = json.Unmarshal(apiResponse, &accountResponse)
	if err != nil {
		return Account{}, err
	}

	return NewAccountFromResponse(accountResponse), nil
}
