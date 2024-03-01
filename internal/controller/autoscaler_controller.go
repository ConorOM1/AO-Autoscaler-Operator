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

	log "github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	scalingv1alpha1 "github.com/ConorOM1/AO-Autoscaler-Operator/api/v1alpha1"
)

// AutoscalerReconciler reconciles a Autoscaler object
type AutoscalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=scaling.autoscaler.project.com,resources=autoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scaling.autoscaler.project.com,resources=autoscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scaling.autoscaler.project.com,resources=autoscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=metrics.k8s.io,resources=pods;nodes,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Autoscaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *AutoscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// TODO(user): your logic here
	// Fetch the Autoscaler instance
	var autoscaler scalingv1alpha1.Autoscaler

	if err := r.Get(ctx, req.NamespacedName, &autoscaler); err != nil {
		log.Error(err, "unable to fetch Autoscaler")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the target Deployment
	var deployment v1.Deployment
	deploymentName := types.NamespacedName{
		Namespace: autoscaler.Namespace,
		Name:      autoscaler.Spec.TargetDeploymentName,
	}
	if err := r.Get(ctx, deploymentName, &deployment); err != nil {
		log.Error(err, "unable to fetch target Deployment")
		return ctrl.Result{}, err
	}

	// Update Deployment replicas if necessary
	minReplicas := autoscaler.Spec.MinReplicas
	if *deployment.Spec.Replicas < minReplicas {
		log.Info("Updating Deployment replicas", "from ", *deployment.Spec.Replicas, "to ", minReplicas)
		deployment.Spec.Replicas = &minReplicas
		if err := r.Update(ctx, &deployment); err != nil {
			log.Error(err, "unable to update Deployment replicas")
			return ctrl.Result{}, err
		}
	}

	maxReplicas := autoscaler.Spec.MaxReplicas
	if *deployment.Spec.Replicas > maxReplicas {
		log.Info("Updating Deployment replicas", "from ", *deployment.Spec.Replicas, "to ", maxReplicas)
		deployment.Spec.Replicas = &maxReplicas
		if err := r.Update(ctx, &deployment); err != nil {
			log.Error(err, "unable to update Deployment replicas")
			return ctrl.Result{}, err
		}
	}

	// Update Autoscaler status
	autoscaler.Status.CurrentReplicas = deployment.Status.Replicas
	if err := r.Status().Update(ctx, &autoscaler); err != nil {
		log.Error(err, "failed to update Autoscaler status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1alpha1.Autoscaler{}).
		Complete(r)
}
