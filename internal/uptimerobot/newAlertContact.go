package uptimerobot

import (
	"context"
	"encoding/json"
	"errors"
)

type NewAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id     int `json:"id"`
		Status int `json:"status"`
	} `json:"alertcontact"`
}

func NewAlertContact(ctx context.Context, client Client, alertType string, value string, friendlyName string) (NewAlertContactResponse, error) {
	if alertType == "" || value == "" {
		return NewAlertContactResponse{}, errors.New("required parameters are empty")
	}

	params := map[string]string{
		"type":  alertType,
		"value": value,
	}

	if friendlyName != "" {
		params["friendly_name"] = friendlyName
	}

	responseBytes, err := makeApiRequest(ctx, client, "getAlertContacts", params)
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
