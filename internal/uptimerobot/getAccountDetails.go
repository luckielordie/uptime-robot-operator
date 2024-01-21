package uptimerobot

import (
	"context"
	"encoding/json"
)

type GetAccountDetailsResponse struct {
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

func GetAccountDetails(ctx context.Context, client UptimeRobotClient) (GetAccountDetailsResponse, error) {
	apiResponse, err := client.MakeApiRequest(ctx, "getAccountDetails", map[string]string{})
	if err != nil {
		return GetAccountDetailsResponse{}, err
	}

	var accountResponse GetAccountDetailsResponse
	err = json.Unmarshal(apiResponse, &accountResponse)
	if err != nil {
		return GetAccountDetailsResponse{}, err
	}

	return accountResponse, nil
}
