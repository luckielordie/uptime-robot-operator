package uptimerobot

import "context"

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

func (c GetAccountDetailsResponse) GetStat() string {
	return c.Stat
}

type AccountDetailsGetter interface {
	GetAccountDetails(ctx context.Context) (GetAccountDetailsResponse, error)
}
