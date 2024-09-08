package alertcontact

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/luckielordie/uptime-robot-operator/internal/uptimerobot"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CreateOrUpdater interface {
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

func NewAlertContact(id string) *AlertContact {
	return &AlertContact{
		Id: id,
	}
}

func IsNotFoundErr(err error) bool {
	if err == nil {
		panic("checking nil error")
	}

	errStr := err.Error()
	return strings.Contains(errStr, "not_found")
}

func createApiResource(ctx context.Context, client CreateOrUpdater, alertContact *AlertContact, mutate func() error) (controllerutil.OperationResult, error) {
	logger := log.FromContext(ctx)
	err := mutate()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	response, err := client.NewAlertContact(ctx, strconv.Itoa(alertContact.Type), alertContact.Value, alertContact.Name)
	if err != nil {
		logger.Info("failed api request", "response", response)
		return controllerutil.OperationResultNone, err
	}

	logger.Info("successful api request", "response", response)
	alertContact.Id = strconv.Itoa(response.AlertContact.Id)
	alertContact.Status = 0

	return controllerutil.OperationResultCreated, nil
}

func updateApiResource(ctx context.Context, client CreateOrUpdater, local *AlertContact, remote AlertContact, mutate func() error) (controllerutil.OperationResult, error) {
	logger := log.FromContext(ctx)
	err := mutate()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	if reflect.DeepEqual(remote, *local) {
		return controllerutil.OperationResultNone, nil
	}

	response, err := client.EditAlertContact(ctx, local.Id, local.Value, local.Name)
	if err != nil {
		logger.Info("failed api request", "response", response)
		return controllerutil.OperationResultNone, err
	}
	logger.Info("successful api request", "response", response)

	*local = remote

	return controllerutil.OperationResultUpdated, nil
}

func CreateOrUpdate(ctx context.Context, client CreateOrUpdater, alertContact *AlertContact, mutate func() error) (controllerutil.OperationResult, error) {
	logger := log.FromContext(ctx)
	//going to assume 0 isn't a valid alert contact id. This will probably bite me in the arse. But they don't give me much choice.
	if alertContact.Id == "" {
		logger.Info("no kube object api resource exists, creating...")
		return createApiResource(ctx, client, alertContact, mutate)
	}

	apiResponse, err := client.GetAlertContacts(ctx, []string{alertContact.Id})
	if err != nil {
		if !IsNotFoundErr(err) {
			return controllerutil.OperationResultNone, err
		}

		logger.Info("kube object api resource unexpectedly deleted, recreating...")
		return createApiResource(ctx, client, alertContact, mutate)
	}

	if apiResponse.AlertContacts == nil {
		return controllerutil.OperationResultNone, errors.New("api returned no alert-contacts when one was expected")
	}

	if len(apiResponse.AlertContacts) != 1 {
		return controllerutil.OperationResultNone, errors.New("api returned more than one alert-contact when only one was expected")
	}

	logger.Info("kube object api resource exists, updating...")
	remote := AlertContact{
		Id:     apiResponse.AlertContacts[0].Id,
		Name:   apiResponse.AlertContacts[0].FriendlyName,
		Type:   apiResponse.AlertContacts[0].Type,
		Status: apiResponse.AlertContacts[0].Status,
		Value:  apiResponse.AlertContacts[0].Value,
	}
	return updateApiResource(ctx, client, alertContact, remote, mutate)
}
