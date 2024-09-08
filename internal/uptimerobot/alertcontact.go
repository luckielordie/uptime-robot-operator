package uptimerobot

import "context"

type NewAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id int `json:"id"`
	} `json:"alertcontact"`
}

func (c NewAlertContactResponse) GetStat() string {
	return c.Stat
}

type AlertContactCreator interface {
	NewAlertContact(ctx context.Context, alertType string, value string, friendlyName string) (NewAlertContactResponse, error)
}

type GetAlertContactResponse struct {
	Stat          string `json:"stat"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
	Total         int    `json:"total"`
	AlertContacts []struct {
		Id           string `json:"id"`
		FriendlyName string `json:"friendly_name"`
		Type         int    `json:"type"`
		Status       int    `json:"status"`
		Value        string `json:"value"`
	} `json:"alert_contacts"`
}

func (c GetAlertContactResponse) GetStat() string {
	return c.Stat
}

type AlertContactGetter interface {
	GetAlertContacts(ctx context.Context, alertContactIds []string) (GetAlertContactResponse, error)
}

type EditAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id int `json:"id"`
	} `json:"alert_contact"`
}

func (c EditAlertContactResponse) GetStat() string {
	return c.Stat
}

type AlertContactEditor interface {
	EditAlertContact(ctx context.Context, id string, value string, friendlyName string) (EditAlertContactResponse, error)
}

type DeleteAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id string `json:"id"`
	} `json:"alert_contact"`
}

func (c DeleteAlertContactResponse) GetStat() string {
	return c.Stat
}

type AlertContactDeleter interface {
	DeleteAlertContact(ctx context.Context, id string) (DeleteAlertContactResponse, error)
}
