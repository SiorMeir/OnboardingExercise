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

package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	exposedeployv1alpha1 "github.com/example/ExposeDeployment/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var exposedeploymentlog = logf.Log.WithName("exposedeployment-resource")

// SetupExposeDeploymentWebhookWithManager registers the webhook for ExposeDeployment in the manager.
func SetupExposeDeploymentWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&exposedeployv1alpha1.ExposeDeployment{}).
		WithValidator(&ExposeDeploymentCustomValidator{}).
		WithDefaulter(&ExposeDeploymentCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-exposedeploy-example-com-v1alpha1-exposedeployment,mutating=true,failurePolicy=fail,sideEffects=None,groups=exposedeploy.example.com,resources=exposedeployments,verbs=create;update,versions=v1alpha1,name=mexposedeployment-v1alpha1.kb.io,admissionReviewVersions=v1

// ExposeDeploymentCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind ExposeDeployment when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ExposeDeploymentCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ExposeDeploymentCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind ExposeDeployment.
func (d *ExposeDeploymentCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	exposedeployment, ok := obj.(*exposedeployv1alpha1.ExposeDeployment)

	if !ok {
		return fmt.Errorf("expected an ExposeDeployment object but got %T", obj)
	}
	exposedeploymentlog.Info("Defaulting for ExposeDeployment", "name", exposedeployment.GetName())
	if exposedeployment.Spec.MinAvailableTimeSec%2 != 0 {
		exposedeploymentlog.Info("MinAvailableTimeSec is odd, doubling it", "name", exposedeployment.GetName(), "minAvailableTimeSec", exposedeployment.Spec.MinAvailableTimeSec)
		exposedeployment.Spec.MinAvailableTimeSec *= 2
	}

	return nil
}

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-exposedeploy-example-com-v1alpha1-exposedeployment,mutating=false,failurePolicy=fail,sideEffects=None,groups=exposedeploy.example.com,resources=exposedeployments,verbs=create;update,versions=v1alpha1,name=vexposedeployment-v1alpha1.kb.io,admissionReviewVersions=v1

// ExposeDeploymentCustomValidator struct is responsible for validating the ExposeDeployment resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ExposeDeploymentCustomValidator struct {
}

var _ webhook.CustomValidator = &ExposeDeploymentCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ExposeDeployment.
func (v *ExposeDeploymentCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	exposedeployment, ok := obj.(*exposedeployv1alpha1.ExposeDeployment)
	if !ok {
		return nil, fmt.Errorf("expected a ExposeDeployment object but got %T", obj)
	}
	exposedeploymentlog.Info("Validation for ExposeDeployment upon creation", "name", exposedeployment.GetName())
	// should add code here
	// RESPONSE: all validations have been done using kubebuilder markers. Do we still need to add code here?
	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ExposeDeployment.
func (v *ExposeDeploymentCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ExposeDeployment.
func (v *ExposeDeploymentCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
