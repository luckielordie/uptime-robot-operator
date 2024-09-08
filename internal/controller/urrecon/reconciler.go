package urrecon

import (
	"context"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ApiObjectCreator[ApiObject any] interface {
	CreateApiObject(ctx context.Context, object *ApiObject) error
}

func createApiResource[ApiObject any](ctx context.Context, creator ApiObjectCreator[ApiObject], object *ApiObject, mutate func() error) (controllerutil.OperationResult, error) {
	logger := log.FromContext(ctx)
	err := mutate()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	logger.Info("no api resource exists, creating...", "object", *object)
	err = creator.CreateApiObject(ctx, object)
	if err != nil {
		logger.Info("error creating api object")
		return controllerutil.OperationResultNone, err
	}

	return controllerutil.OperationResultCreated, nil
}

type ApiObjectEditor[ApiObject any] interface {
	EditApiObject(ctx context.Context, object *ApiObject) error
}

func updateApiResource[ApiObject any](ctx context.Context, editor ApiObjectEditor[ApiObject], local *ApiObject, remote *ApiObject, mutate func() error) (controllerutil.OperationResult, error) {
	logger := log.FromContext(ctx)
	err := mutate()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	if reflect.DeepEqual(remote, local) {
		logger.Info("kube object api resource is in sync")
		return controllerutil.OperationResultNone, nil
	}

	logger.Info("kube object and api resource out of sync, updating api resource", "from", *remote, "to", *local)
	err = editor.EditApiObject(ctx, local)
	if err != nil {
		logger.Info("error editing api object")
		return controllerutil.OperationResultNone, err
	}

	*local = *remote

	return controllerutil.OperationResultUpdated, nil
}

type ApiObjectReconciler[ApiObject any] interface {
	ApiObjectCreator[ApiObject]
	ApiObjectEditor[ApiObject]
	ApiObjectExists(ctx context.Context, object *ApiObject) (bool, error)
	GetApiObject(ctx context.Context, object *ApiObject) (*ApiObject, error)
}

func ReconcileApiObject[ApiObject any](ctx context.Context, reconciler ApiObjectReconciler[ApiObject], object *ApiObject, mutate func() error) (controllerutil.OperationResult, error) {
	exists, err := reconciler.ApiObjectExists(ctx, object)
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	if !exists {
		return createApiResource[ApiObject](ctx, reconciler, object, mutate)
	}

	remote, err := reconciler.GetApiObject(ctx, object)
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	return updateApiResource[ApiObject](ctx, reconciler, object, remote, mutate)
}
