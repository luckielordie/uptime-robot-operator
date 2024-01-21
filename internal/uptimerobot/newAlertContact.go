package uptimerobot

import (
	"context"
	"encoding/json"
)

type NewAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id     int `json:"id"`
		Status int `json:"status"`
	} `json:"alertcontact"`
}

func NewAlertContact(ctx context.Context, client UptimeRobotClient, alertType string, value string, friendlyName string) (NewAlertContactResponse, error) {
	params := map[string]string{
		"type":  alertType,
		"value": value,
	}

	if friendlyName != "" {
		params["friendly_name"] = friendlyName
	}

	responseBytes, err := client.MakeApiRequest(ctx, "getAlertContacts", params)
	if err != nil {
		return NewAlertContactResponse{}, err
	}

	var newAlertContactResponse NewAlertContactResponse
	err = json.Unmarshal(responseBytes, &newAlertContactResponse)
	if err != nil {
		return NewAlertContactResponse{}, err
	}

	return newAlertContactResponse, nil
}
