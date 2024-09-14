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
	FriendlyName                     string   `json:"friendly_name"`
	Url                              string   `json:"url"`
	MonitorType                      int      `json:"type"`
	SubType                          int      `json:"sub_type"`
	Port                             int      `json:"port"`
	KeywordType                      int      `json:"keyword_type"`
	KeywordCaseType                  int      `json:"keyword_case_type"`
	KeywordValue                     string   `json:"keyword_value"`
	Interval                         int      `json:"interval"`
	Timeout                          int      `json:"timeout"`
	HttpUsername                     string   `json:"http_username"`
	HttpPassword                     string   `json:"http_password"`
	HttpAuthType                     int      `json:"http_auth_type"`
	PostType                         int      `json:"post_type"`
	PostValue                        string   `json:"post_value"`
	HttpMethod                       string   `json:"http_method"`
	PostContentType                  int      `json:"post_content_type"`
	AlertContacts                    []string `json:"alert_contacts"`
	MaintenanceWindows               string   `json:"mwindows"`
	CustomHttpHeaders                string   `json:"custom_http_headers"`
	CustomHttpStatuses               string   `json:"custom_http_statuses"`
	IgnoreSSLErrors                  bool     `json:"ignore_ssl_errors"`
	DisableDomainExpireNotifications bool     `json:"disable_domain_expire_notifications"`
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

type EditMonitorRequest struct {
	Id                               string   `json:"id"`
	FriendlyName                     string   `json:"friendly_name"`
	Url                              string   `json:"url"`
	SubType                          int      `json:"sub_type"`
	Port                             int      `json:"port"`
	KeywordType                      int      `json:"keyword_type"`
	KeywordCaseType                  int      `json:"keyword_case_type"`
	KeywordValue                     string   `json:"keyword_value"`
	Interval                         int      `json:"interval"`
	Timeout                          int      `json:"timeout"`
	Status                           int      `json:"status"`
	HttpUsername                     string   `json:"http_username"`
	HttpPassword                     string   `json:"http_password"`
	HttpAuthType                     int      `json:"http_auth_type"`
	HttpMethod                       string   `json:"http_method"`
	PostType                         int      `json:"post_type"`
	PostValue                        string   `json:"post_value"`
	PostContentType                  int      `json:"post_content_type"`
	AlertContacts                    []string `json:"alert_contacts"`
	MaintenanceWindows               string   `json:"mwindows"`
	CustomHttpHeaders                string   `json:"custom_http_headers"`
	CustomHttpStatuses               string   `json:"custom_http_statuses"`
	IgnoreSSLErrors                  bool     `json:"ignore_ssl_errors"`
	DisableDomainExpireNotifications bool     `json:"disable_domain_expire_notifications"`
}

type EditMonitorResponse struct {
	Stat    string `json:"stat"`
	Monitor struct {
		Id int `json:"id"`
	} `json:"monitor"`
}

func (c EditMonitorResponse) GetStat() string {
	return c.Stat
}

type MonitorEditor interface {
	EditMonitor(ctx context.Context, request EditMonitorRequest) (EditMonitorResponse, error)
}

type GetMonitorResponse struct {
	Stat       string `json:"stat"`
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Monitors []struct {
		Id              string `json:"id"`
		FriendlyName    string `json:"friendly_name"`
		Url             string `json:"url"`
		MonitorType     int    `json:"type"`
		SubType         int    `json:"sub_type"`
		KeywordType     int    `json:"keyword_type"`
		KeywordCaseType int    `json:"keyword_case_type"`
		KeywordValue    string `json:"keyword_value"`
		HttpUsername    string `json:"http_username"`
		HttpPassword    string `json:"http_password"`
		Port            int    `json:"port"`
		Interval        int    `json:"interval"`
		Status          int    `json:"status"`
		CreateDatetime  int    `json:"create_datetime"`
		MonitorGroup    int    `json:"monitor_group"`
		IsGroupMain     int    `json:"is_group_main"`
		Logs            struct {
			Type     int `json:"type"`
			Datetime int `json:"datetime"`
			Duration int `json:"duration"`
		} `json:"logs"`
	}
}

func (c GetMonitorResponse) GetStat() string {
	return c.Stat
}

type MonitorGetter interface {
	GetMonitors(ctx context.Context, monitorIds []string) (GetMonitorResponse, error)
}
