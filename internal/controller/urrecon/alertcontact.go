package urrecon

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/luckielordie/uptime-robot-operator/internal/uptimerobot"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type AlertContactApiClient interface {
	uptimerobot.AlertContactCreator
	uptimerobot.AlertContactEditor
	uptimerobot.AlertContactGetter
}

type AlertContact struct {
	Id     string
	Name   string
	Type   int
	Status int
	Value  string
}

type AlertContactApiReconciler struct {
	apiClient AlertContactApiClient
}

func NewAlertContactApiReconciler(apiClient AlertContactApiClient) AlertContactApiReconciler {
	return AlertContactApiReconciler{
		apiClient: apiClient,
	}
}

func (reconciler *AlertContactApiReconciler) CreateApiObject(ctx context.Context, alertContact *AlertContact) error {
	logger := log.FromContext(ctx)
	response, err := reconciler.apiClient.NewAlertContact(ctx, strconv.Itoa(alertContact.Type), alertContact.Value, alertContact.Name)
	if err != nil {
		logger.Info("failed api request", "response", response)
		return err
	}

	logger.Info("successful api request", "response", response)
	alertContact.Id = strconv.Itoa(response.AlertContact.Id)
	alertContact.Status = 0

	return nil
}

func (reconciler *AlertContactApiReconciler) EditApiObject(ctx context.Context, alertContact *AlertContact) error {
	logger := log.FromContext(ctx)

	response, err := reconciler.apiClient.EditAlertContact(ctx, alertContact.Id, alertContact.Value, alertContact.Name)
	if err != nil {
		logger.Info("failed api request", "response", response)
		return err
	}
	logger.Info("successful api request", "response", response)

	return nil
}

func (reconciler *AlertContactApiReconciler) ApiObjectExists(ctx context.Context, alertContact *AlertContact) (bool, error) {
	// if id is an empty string then the object can't exist
	if alertContact.Id == "" {
		return false, nil
	}

	//check if object exists through the API
	_, err := reconciler.apiClient.GetAlertContacts(ctx, []string{alertContact.Id})
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "not_found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (reconciler *AlertContactApiReconciler) GetApiObject(ctx context.Context, alertContact *AlertContact) (*AlertContact, error) {
	apiResponse, err := reconciler.apiClient.GetAlertContacts(ctx, []string{alertContact.Id})
	if err != nil {
		return nil, fmt.Errorf("unexpected error with alert-contact api: %w", err)
	}

	if apiResponse.AlertContacts == nil {
		return nil, errors.New("api returned no alert-contacts when one was expected")
	}

	if len(apiResponse.AlertContacts) != 1 {
		return nil, errors.New("api returned more than one alert-contact when only one was expected")
	}

	return &AlertContact{
		Id:     apiResponse.AlertContacts[0].Id,
		Name:   apiResponse.AlertContacts[0].FriendlyName,
		Type:   apiResponse.AlertContacts[0].Type,
		Status: apiResponse.AlertContacts[0].Status,
		Value:  apiResponse.AlertContacts[0].Value,
	}, nil
}
