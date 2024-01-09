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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	uptimerobotcomv1alpha1 "github.com/luckielordie/uptime-robot-operator/api/v1alpha1"
)

// AccountReconciler reconciles a Account object
type AccountReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func getAccount(ctx context.Context, reader client.Reader, req ctrl.Request) (uptimerobotcomv1alpha1.Account, error) {
	account := uptimerobotcomv1alpha1.Account{}
	err := reader.Get(ctx, req.NamespacedName, &account)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "failed to retrieve bucket resource")
		}

		logger.Info("requeue", "reason", "failed to get")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
}

//+kubebuilder:rbac:groups=uptimerobot.com.uptimerobot.com,resources=accounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimerobot.com.uptimerobot.com,resources=accounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimerobot.com.uptimerobot.com,resources=accounts/finalizers,verbs=update

func (reconciler *AccountReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	account, err := getAccount(ctx, reconciler, request)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Error(err, "failed to retrieve account resource")
		}

		logger.Info("requeue", "reason", "failed to get")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AccountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimerobotcomv1alpha1.Account{}).
		Complete(r)
}
