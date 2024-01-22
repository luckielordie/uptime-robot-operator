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
	"errors"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	uptimerobotcomv1alpha1 "github.com/luckielordie/uptime-robot-operator/api/v1alpha1"
	"github.com/luckielordie/uptime-robot-operator/internal/controller/alertcontact"
)

// AlertContactReconciler reconciles a AlertContact object
type AlertContactReconciler struct {
	client.Client
	Scheme             *runtime.Scheme
	AlertContactClient alertcontact.Client
}

func getAlertContact(ctx context.Context, reader client.Reader, req ctrl.Request) (uptimerobotcomv1alpha1.AlertContact, error) {
	logger := log.FromContext(ctx)
	alertContact := uptimerobotcomv1alpha1.AlertContact{}
	err := reader.Get(ctx, req.NamespacedName, &alertContact)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "failed to retrieve alert contact resource")
		}

		logger.Info("requeue", "reason", "failed to get")
		return uptimerobotcomv1alpha1.AlertContact{}, err
	}

	return alertContact, nil
}

const utFinalizer = "uptimerobot.com/finalizer"

//+kubebuilder:rbac:groups=uptimerobot.com,resources=alertcontacts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimerobot.com,resources=alertcontacts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimerobot.com,resources=alertcontacts/finalizers,verbs=update

func (reconciler *AlertContactReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	alertContact, err := getAlertContact(ctx, reconciler, request)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !controllerutil.ContainsFinalizer(&alertContact, utFinalizer) {
		logger.Info("adding finalizer for alert contact")
		added := controllerutil.AddFinalizer(&alertContact, utFinalizer)
		if !added {
			err = errors.New("failed to add finalizer")
			logger.Error(err, "failed to add finalizer")
			return ctrl.Result{}, err
		}

		err = reconciler.Update(ctx, &alertContact)
		if err != nil {
			logger.Error(err, "failed to update resource")
			return ctrl.Result{}, err
		}
	}

	isToBeDeleted := alertContact.GetDeletionTimestamp() != nil
	if isToBeDeleted {
		if controllerutil.ContainsFinalizer(&alertContact, utFinalizer) {
			logger.Info("finalizing alert contact")

			_, err := reconciler.AlertContactClient.DeleteAlertContact(ctx, alertContact.Status.Id)
			if err != nil {
				logger.Error(err, "failed to delete alert contact", "id", alertContact.Status.Id)
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(&alertContact, utFinalizer)
			err = reconciler.Update(ctx, &alertContact)
			if err != nil {
				logger.Error(err, "failed to update resource after removing finalizer")
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil //nothing to be done, no requeue
	}

	//CreateOrUpdate AlertContact
	alertContactObj := alertcontact.AlertContact{
		Id: alertContact.Status.Id,
	}

	result, err := alertcontact.CreateOrUpdate(ctx, reconciler.AlertContactClient, &alertContactObj, func() error {
		alertContactObj.Name = alertContact.Spec.Name
		alertContactObj.Type = alertContact.Spec.Type
		alertContactObj.Value = alertContact.Spec.Value
		alertContactObj.Status = alertContact.Status.Status
		return nil
	})
	if err != nil {
		logger.Error(err, "failed updating alertcontact on api")
		return ctrl.Result{}, err
	}

	if result != controllerutil.OperationResultNone {
		alertContact.Status.Id = alertContactObj.Id
		alertContact.Status.Status = alertContactObj.Status

		statusClient := reconciler.Client.Status()
		err = statusClient.Update(ctx, &alertContact)
		if err != nil {
			logger.Error(err, "failed updating status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{
		RequeueAfter: time.Second * 30,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertContactReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimerobotcomv1alpha1.AlertContact{}).
		Complete(r)
}
