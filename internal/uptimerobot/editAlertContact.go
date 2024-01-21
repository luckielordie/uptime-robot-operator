package uptimerobot

import (
	"context"
	"encoding/json"
	"strconv"
)

type EditAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id int `json:"id"`
	} `json:"alert_contact"`
}

func EditAlertContact(ctx context.Context, client UptimeRobotClient, id int, value string, friendlyName string) (EditAlertContactResponse, error) {
	params := map[string]string{
		"id":    strconv.Itoa(id),
		"value": value,
	}

	if friendlyName != "" {
		params["friendly_name"] = friendlyName
	}

	responseBytes, err := client.MakeApiRequest(ctx, "editAlertContact", params)
	if err != nil {
		return EditAlertContactResponse{}, err
	}

	var editAlertContactResponse EditAlertContactResponse
	err = json.Unmarshal(responseBytes, &editAlertContactResponse)
	if err != nil {
		return EditAlertContactResponse{}, err
	}

	return editAlertContactResponse, nil
}
