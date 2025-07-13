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
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	exposedeployv1alpha1 "github.com/example/ExposeDeployment/api/v1alpha1"
)

// ExposeDeploymentReconciler reconciles a ExposeDeployment object
type ExposeDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const finalizerName = "exposedeploy.finalizers.example.com"

// +kubebuilder:rbac:groups=exposedeploy.example.com,resources=exposedeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=exposedeploy.example.com,resources=exposedeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=exposedeploy.example.com,resources=exposedeployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
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

	// examine DeletionTimestamp to determine if object is under deletion
	if exposedeploy.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then let's add the finalizer and update the object. This is equivalent
		// to registering our finalizer.
		if !controllerutil.ContainsFinalizer(exposedeploy, finalizerName) {
			controllerutil.AddFinalizer(exposedeploy, finalizerName)
			if err := r.Update(ctx, exposedeploy); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(exposedeploy, finalizerName) {
			// our finalizer is present, so let's handle any external dependency
			if _, err := r.gracefullyDeleteResourceAndService(ctx, podList, exposedeploy); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried.
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(exposedeploy, finalizerName)
			if err := r.Update(ctx, exposedeploy); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

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
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "exposedeploy.example.com/v1alpha1",
							Kind:       "ExposeDeployment",
							Name:       exposedeploy.Name,
							UID:        exposedeploy.UID,
						},
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:    "expose-deployment",
						Image:   exposedeploy.Spec.Image,
						Command: []string{"sleep", "infinity"},
						Args:    exposedeploy.Spec.CustomEnv,
					}},
				},
			}
			if err := r.Create(ctx, &pod); err != nil {
				return ctrl.Result{}, err
			}

		}
	}

	// 4. Scale down if needed ,exclude terminating pods by deletionTimestamp
	numOfTerminatingPods := calculateNumOfTerminatingPods(podList)
	if currentPodCount-numOfTerminatingPods > desiredReplicas {
		excess := currentPodCount - desiredReplicas - numOfTerminatingPods
		for i := range excess {
			pod := podList.Items[i]
			if err := r.Delete(ctx, &pod); err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	// update status of exposedeploy
	r.updateExposeDeploymentStatus(exposedeploy, podList)

	// create service of exposedeploy based on the portdefinition
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      exposedeploy.Name + "-service",
			Namespace: exposedeploy.Namespace,
			Labels:    labelsForApp(exposedeploy.Name),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: exposedeploy.Spec.PortDefinition.Port,
					TargetPort: intstr.IntOrString{
						IntVal: exposedeploy.Spec.PortDefinition.TargetPort,
					},
				},
			},
		},
	}
	if err = r.Create(ctx, &service); err != nil {
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

// Helper functions //

// labelsForApp creates a simple set of labels for ExposeDeployment.
func labelsForApp(name string) map[string]string {
	return map[string]string{"cr_name": name}
}

func calculateNumOfTerminatingPods(podList *corev1.PodList) int {
	numOfTerminatingPods := 0
	for _, pod := range podList.Items {
		if pod.DeletionTimestamp != nil {
			numOfTerminatingPods++
		}
	}
	return numOfTerminatingPods
}

// calculateAvailablePods creates a list of available pods with respect to the minavailabletimesec
func (r *ExposeDeploymentReconciler) calculateAvailablePods(exposedeploy *exposedeployv1alpha1.ExposeDeployment, podList *corev1.PodList) []corev1.Pod {
	availablePods := []corev1.Pod{}
	for _, pod := range podList.Items {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady {
				timeSinceCreation := time.Since(condition.LastTransitionTime.Time)
				if pod.DeletionTimestamp == nil && timeSinceCreation > time.Duration(exposedeploy.Spec.MinAvailableTimeSec)*time.Second {
					availablePods = append(availablePods, pod)
				}
			}
		}
	}
	return availablePods
}

func (r *ExposeDeploymentReconciler) updateExposeDeploymentStatus(exposedeploy *exposedeployv1alpha1.ExposeDeployment, podList *corev1.PodList) {
	availablePods := r.calculateAvailablePods(exposedeploy, podList)
	exposedeploy.Status.AvailablePods = int32(len(availablePods))
	exposedeploy.Status.ReadyPods = int32(len(availablePods))

	// loop through conditions, if the condition is LastReconcileSucceeded, update the status to true
	for _, condition := range exposedeploy.Status.Conditions {
		switch condition.Type {
		case exposedeployv1alpha1.LastReconcileSucceeded:
			condition.Status = true
			condition.LastTransitionTime = metav1.Now()
			condition.Reason = "updateResource"
			condition.Message = "ExposeDeployment reconciled successfully"
		case exposedeployv1alpha1.Available:
			if condition.Type == exposedeployv1alpha1.Available {
				if exposedeploy.Status.AvailablePods == exposedeploy.Spec.Replicas {
					condition.Status = true
					condition.LastTransitionTime = metav1.Now()
					condition.Reason = "allPodsAvailable"
					condition.Message = "ExposeDeployment is available and service is as expected"
				} else {
					condition.Status = false
					condition.LastTransitionTime = metav1.Now()
					condition.Reason = "notAllPodsAvailable"
					condition.Message = "ExposeDeployment is not available and service is not as expected"
				}
			}
		}
	}
}

func (r *ExposeDeploymentReconciler) gracefullyDeleteResourceAndService(ctx context.Context, podList *corev1.PodList, exposedeploy *exposedeployv1alpha1.ExposeDeployment) (ctrl.Result, error) {
	// delete all pods of exposedeploy and wait for them to be deleted
	for _, pod := range podList.Items {
		if err := r.Delete(ctx, &pod); err != nil {
			return ctrl.Result{}, err
		}
	}
	// wait for all pods to be deleted
	for _, pod := range podList.Items {
		if err := r.Get(ctx, client.ObjectKey{Name: pod.Name, Namespace: pod.Namespace}, &pod); err != nil {
			return ctrl.Result{}, err
		}
	}
	svc := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{
		Name:      exposedeploy.Name + "-service",
		Namespace: exposedeploy.Namespace,
	}, svc)
	if err == nil {
		if err := r.Delete(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
	} else if !errors.IsNotFound(err) {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
