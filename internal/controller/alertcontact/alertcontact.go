package alertcontact

import (
	"context"
	"reflect"
	"strconv"

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
	return false
}

func CreateOrUpdate(ctx context.Context, client CreateOrUpdater, alertContact *AlertContact, mutate func() error) (controllerutil.OperationResult, error) {
	logger := log.FromContext(ctx)
	//going to assume 0 isn't a valid alert contact id. This will probably bite me in the arse. But they don't give me much choice.
	if alertContact.Id == "" {
		logger.Info("no api resource exists, creating...")
		err := mutate()
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		response, err := client.NewAlertContact(ctx, strconv.Itoa(alertContact.Type), alertContact.Value, alertContact.Name)
		if err != nil {
			logger.Info("failed api request", "response", response, "request", alertContact)
			return controllerutil.OperationResultNone, err
		}

		logger.Info("successful api request", "response", response, "request", alertContact)
		alertContact.Id = response.AlertContact.Id
		alertContact.Status = response.AlertContact.Status

		return controllerutil.OperationResultCreated, nil
	}

	logger.Info("api resource exists, updating...")

	apiResponse, err := client.GetAlertContacts(ctx, []string{alertContact.Id})
	if !IsNotFoundErr(err) {
		return controllerutil.OperationResultNone, err
	}

	if apiResponse.AlertContacts == nil {
		return controllerutil.OperationResultNone, err
	}

	if len(apiResponse.AlertContacts) != 1 {
		return controllerutil.OperationResultNone, err
	}

	existing := AlertContact{
		Id:     apiResponse.AlertContacts[0].Id,
		Name:   apiResponse.AlertContacts[0].FriendlyName,
		Type:   apiResponse.AlertContacts[0].Type,
		Status: apiResponse.AlertContacts[0].Status,
		Value:  apiResponse.AlertContacts[0].Value,
	}

	err = mutate()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	if reflect.DeepEqual(existing, alertContact) {
		return controllerutil.OperationResultNone, nil
	}

	_, err = client.EditAlertContact(ctx, alertContact.Id, alertContact.Value, alertContact.Name)
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	alertContact = &existing

	return controllerutil.OperationResultUpdated, nil
}
