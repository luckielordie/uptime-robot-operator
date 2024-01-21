package uptimerobot

import (
	"context"
	"encoding/json"
	"strconv"
)

type DeleteAlertContactResponse struct {
	Stat         string `json:"stat"`
	AlertContact struct {
		Id int `json:"id"`
	} `json:"alert_contact"`
}

func DeleteAlertContact(ctx context.Context, client UptimeRobotClient, id int) (DeleteAlertContactResponse, error) {
	params := map[string]string{
		"id": strconv.Itoa(id),
	}

	responseBytes, err := client.MakeApiRequest(ctx, "deleteAlertContact", params)
	if err != nil {
		return DeleteAlertContactResponse{}, err
	}

	var deleteAlertContactResponse DeleteAlertContactResponse
	err = json.Unmarshal(responseBytes, &deleteAlertContactResponse)
	if err != nil {
		return DeleteAlertContactResponse{}, err
	}

	return deleteAlertContactResponse, nil
}
