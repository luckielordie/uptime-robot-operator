package uptimerobot

import "context"

type NewMonitorResponse struct {
	Stat    string `json:"stat"`
	Monitor struct {
		Id     int `json:"id"`
		Status int `json:"status"`
	} `json:"monitor"`
}

func (c NewMonitorResponse) GetStat() string {
	return c.Stat
}

type NewMonitorRequest struct {
	FriendlyName                     string
	Url                              string
	MonitorType                      int
	SubType                          int
	Port                             int
	KeywordType                      int
	KeywordCaseType                  int
	KeywordValue                     string
	Interval                         int
	Timeout                          int
	HttpUsername                     string
	HttpPassword                     string
	HttpAuthType                     int
	PostType                         int
	PostValue                        string
	HttpMethod                       string
	PostContentType                  int
	AlertContacts                    string
	MaintenanceWindows               string
	CustomHttpHeaders                string
	CustomHttpStatuses               string
	IgnoreSSLErrors                  bool
	DisableDomainExpireNotifications bool
}

type MonitorCreator interface {
	NewMonitor(ctx context.Context, request NewMonitorRequest) (NewMonitorResponse, error)
}

type DeleteMonitorResponse struct {
	Stat    string `json:"stat"`
	Monitor struct {
		Id int `json:"id"`
	} `json:"monitor"`
}

func (c DeleteMonitorResponse) GetStat() string {
	return c.Stat
}

type MonitorDeleter interface {
	DeleteMonitor(ctx context.Context, id int) (DeleteMonitorResponse, error)
}
