package controller

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Finalize(ctx context.Context, reconciler client.Client, object client.Object, finalizerString string, finaliser func(context.Context) error) (controllerutil.OperationResult, error) {
	isMarkedToBeDeleted := object.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(object, finalizerString) {
			err := finaliser(ctx)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}

			controllerutil.RemoveFinalizer(object, finalizerString)
			err = reconciler.Update(ctx, object)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}

			return controllerutil.OperationResultUpdated, nil
		}
	}

	if !controllerutil.ContainsFinalizer(object, finalizerString) {
		controllerutil.AddFinalizer(object, finalizerString)
		err := reconciler.Update(ctx, object)
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		return controllerutil.OperationResultUpdated, nil
	}

	return controllerutil.OperationResultNone, nil
}
