/*
Copyright 2025.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	exposedeployv1alpha1 "github.com/example/ExposeDeployment/api/v1alpha1"
)

// ExposeDeploymentReconciler reconciles a ExposeDeployment object
type ExposeDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=exposedeploy.example.com,resources=exposedeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=exposedeploy.example.com,resources=exposedeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=exposedeploy.example.com,resources=exposedeployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ExposeDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *ExposeDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	exposedeploy := &exposedeployv1alpha1.ExposeDeployment{}
	err := r.Get(ctx, req.NamespacedName, exposedeploy)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// list all pods of exposedeploy
	podList := &corev1.PodList{}
	if err = r.List(ctx, podList, client.InNamespace(exposedeploy.Namespace), client.MatchingLabels(labelsForApp(exposedeploy.Name))); err != nil {
		return ctrl.Result{}, err
	}

	// currentPodCount := len(podList.Items)
	desiredReplicas := int(exposedeploy.Spec.Replicas)
	currentPodCount := len(podList.Items)

	// 3. Scale up if needed
	if currentPodCount < desiredReplicas {
		for i := len(podList.Items); i < desiredReplicas; i++ {
			pod := corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: exposedeploy.Name + "-pod-",
					Namespace:    req.Namespace,
					Labels:       labelsForApp(exposedeploy.Name),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:    "expose-deployment",
						Image:   exposedeploy.Spec.Image,
						Command: []string{"sleep", "infinity"},
					}},
				},
			}
			if err := r.Create(ctx, &pod); err != nil {
				return ctrl.Result{}, err
			}

		}
	}

	// 4. Scale down if needed ,exclude terminating pods by deletionTimestamp
	if len(podList.Items) > desiredReplicas {
		excess := len(podList.Items) - desiredReplicas
		for i := range excess {
			pod := podList.Items[i]
			if err := r.Delete(ctx, &pod); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// create service of exposedeploy based on the portdefinition
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "expose-deploy",
			Namespace: exposedeploy.Namespace,
		},
	}
	if err = r.Create(ctx, service); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExposeDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&exposedeployv1alpha1.ExposeDeployment{}).
		Named("exposedeployment").
		Complete(r)
}

// labelsForApp creates a simple set of labels for ExposeDeployment.
func labelsForApp(name string) map[string]string {
	return map[string]string{"cr_name": name}
}
