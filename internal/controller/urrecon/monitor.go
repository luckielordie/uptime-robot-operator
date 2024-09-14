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

type MonitorApiClient interface {
	uptimerobot.MonitorCreator
	uptimerobot.MonitorEditor
	uptimerobot.MonitorGetter
}

type Monitor struct {
	Id   string
	Name string
}

type MonitorApiReconciler struct {
	apiClient MonitorApiClient
}

func NewMonitorApiReconciler(apiClient MonitorApiClient) MonitorApiReconciler {
	return MonitorApiReconciler{
		apiClient: apiClient,
	}
}

func (reconciler *MonitorApiReconciler) CreateApiObject(ctx context.Context, monitor *Monitor) error {
	logger := log.FromContext(ctx)
	response, err := reconciler.apiClient.NewMonitor(ctx, uptimerobot.NewMonitorRequest{
		FriendlyName: monitor.Name,
	})
	if err != nil {
		logger.Info("failed api request", "response", response)
		return err
	}

	logger.Info("successful api request", "response", response)
	monitor.Id = strconv.Itoa(response.Monitor.Id)

	return nil
}

func (reconciler *MonitorApiReconciler) EditApiObject(ctx context.Context, monitor *Monitor) error {
	logger := log.FromContext(ctx)

	response, err := reconciler.apiClient.EditMonitor(ctx, uptimerobot.EditMonitorRequest{
		Id:           monitor.Id,
		FriendlyName: monitor.Name,
	})
	if err != nil {
		logger.Info("failed api request", "response", response)
		return err
	}
	logger.Info("successful api request", "response", response)

	return nil
}

func (reconciler *MonitorApiReconciler) ApiObjectExists(ctx context.Context, monitor *Monitor) (bool, error) {
	// if id is an empty string then the object can't exist
	if monitor.Id == "" {
		return false, nil
	}

	//check if object exists through the API
	_, err := reconciler.apiClient.GetMonitors(ctx, []string{monitor.Id})
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "not_found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (reconciler *MonitorApiReconciler) GetApiObject(ctx context.Context, monitor *Monitor) (*Monitor, error) {
	apiResponse, err := reconciler.apiClient.GetMonitors(ctx, []string{monitor.Id})
	if err != nil {
		return nil, fmt.Errorf("unexpected error with monitor api: %w", err)
	}

	if apiResponse.Monitors == nil {
		return nil, errors.New("api returned no monitors when one was expected")
	}

	if len(apiResponse.Monitors) != 1 {
		return nil, errors.New("api returned more than one monitor when only one was expected")
	}

	return &Monitor{
		Id:   apiResponse.Monitors[0].Id,
		Name: apiResponse.Monitors[0].FriendlyName,
	}, nil
}
