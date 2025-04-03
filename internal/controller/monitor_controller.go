/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strconv"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	uptimerobotcomv1alpha1 "github.com/luckielordie/uptime-robot-operator/api/v1alpha1"
	"github.com/luckielordie/uptime-robot-operator/internal/controller/urrecon"
	"github.com/luckielordie/uptime-robot-operator/internal/uptimerobot"
)

type MonitorAPIClient interface {
	uptimerobot.MonitorCreator
	uptimerobot.MonitorEditor
	uptimerobot.MonitorGetter
	uptimerobot.MonitorDeleter
}

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	urrecon.MonitorApiReconciler
	Scheme        *runtime.Scheme
	MonitorClient MonitorAPIClient
}

func getMonitor(ctx context.Context, reader client.Reader, req ctrl.Request) (uptimerobotcomv1alpha1.Monitor, error) {
	logger := log.FromContext(ctx)
	monitor := uptimerobotcomv1alpha1.Monitor{}
	err := reader.Get(ctx, req.NamespacedName, &monitor)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "failed to retrieve alert contact resource")
		}

		logger.Info("requeue", "reason", "failed to get")
		return uptimerobotcomv1alpha1.Monitor{}, err
	}

	return monitor, nil
}

func getListOfAlertContactIds(ctx context.Context, reader client.Reader, labels map[string]string) ([]string, error) {
	logger := log.FromContext(ctx)
	alertContacts := uptimerobotcomv1alpha1.AlertContactList{}
	var matchingLabels client.MatchingLabels = labels
	err := reader.List(ctx, &alertContacts, matchingLabels)
	if err != nil {
		logger.Info("failed to retrieve alert contacts for monitor with labels", "labels", labels)
		return nil, err
	}

	var ids []string
	for _, ac := range alertContacts.Items {
		if ac.Status.Id != "" {
			ids = append(ids, ac.Status.Id)
		}
	}

	return ids, nil
}

//+kubebuilder:rbac:groups=uptimerobot.com,resources=monitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimerobot.com,resources=monitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimerobot.com,resources=monitors/finalizers,verbs=update
//+kubebuilder:rbac:groups=uptimerobot.com,resources=alertcontacts,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Monitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (reconciler *MonitorReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	monitor, err := getMonitor(ctx, reconciler, request)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	result, err := Finalize(ctx, reconciler.Client, &monitor, FINALIZER_TOKEN, func(context.Context) error {
		idInt, err := strconv.Atoi(monitor.Status.Id)
		if err != nil {
			return err
		}

		_, err = reconciler.MonitorClient.DeleteMonitor(ctx, idInt)
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "not_found") {
				return nil
			}
			logger.Error(err, "failed to delete monitor", "id", monitor.Status.Id)
			return err
		}

		return nil
	})
	if err != nil || result != controllerutil.OperationResultNone {
		if result != controllerutil.OperationResultNone {
			logger.Error(err, "failed finalizing monitor")
		}

		return ctrl.Result{}, err
	}

	alertContactIds, err := getListOfAlertContactIds(ctx, reconciler, monitor.Spec.AlertContacts.MatchLabels)
	if err != nil {
		return ctrl.Result{}, err
	}

	monitorObj := urrecon.Monitor{
		Id:            monitor.Status.Id,
		AlertContacts: alertContactIds,
	}

	result, err = urrecon.ReconcileApiObject[urrecon.Monitor](ctx, reconciler, &monitorObj, func() error {
		monitorObj.Name = monitor.Spec.Name
		monitorObj.Url = monitor.Spec.Url
		monitorObj.AlertContacts = alertContactIds
		return nil
	})
	if err != nil {
		logger.Error(err, "failed updating monitor on api")
		return ctrl.Result{}, err
	}

	if result != controllerutil.OperationResultNone {
		monitor.Status.Id = monitorObj.Id
		monitor.Status.Name = monitorObj.Name
		monitor.Status.Url = monitorObj.Url

		statusClient := reconciler.Client.Status()
		err = statusClient.Update(ctx, &monitor)
		if err != nil {
			logger.Error(err, "failed updating status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{
		RequeueAfter: time.Second * 15,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimerobotcomv1alpha1.Monitor{}).
		Complete(r)
}
