package uptimerobot

import (
	"context"
	"encoding/json"
	"fmt"
)

type GetAlertContactResponse struct {
	Stat          string `json:"stat"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
	Total         int    `json:"total"`
	AlertContacts []struct {
		Id           int    `json:"id"`
		FriendlyName string `json:"friendly_name"`
		Type         int    `json:"type"`
		Status       int    `json:"status"`
		Value        string `json:"value"`
	} `json:"alert_contacts"`
}

func GetAlertContacts(ctx context.Context, client UptimeRobotClient, alertContactIds []int) (GetAlertContactResponse, error) {
	params := map[string]string{}
	if alertContactIds != nil {
		for _, id := range alertContactIds {
			if params["alert_contacts"] == "" {
				params["alert_contacts"] = fmt.Sprintf("%d", id)
			} else {
				params["alert_contacts"] = fmt.Sprintf("%s-%d", params["alert_contacts"], id)
			}
		}
	}

	responseBytes, err := client.MakeApiRequest(ctx, "getAlertContacts", params)
	if err != nil {
		return GetAlertContactResponse{}, err
	}

	var getAlertContactResponse GetAlertContactResponse
	err = json.Unmarshal(responseBytes, &getAlertContactResponse)
	if err != nil {
		return GetAlertContactResponse{}, err
	}

	return getAlertContactResponse, nil
}
