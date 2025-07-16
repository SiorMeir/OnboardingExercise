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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	exposedeployv1alpha1 "github.com/example/ExposeDeployment/api/v1alpha1"
)

var _ = Describe("ExposeDeployment Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "my-namespace",
		}
		exposedeployment := &exposedeployv1alpha1.ExposeDeployment{}

		BeforeEach(func() {
			By("creating the namespace if not exists")
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-namespace",
				},
			}
			Expect(k8sClient.Create(ctx, namespace)).To(Succeed())

			By("creating the custom resource for the Kind ExposeDeployment")
			err := k8sClient.Get(ctx, typeNamespacedName, exposedeployment)
			if err != nil && errors.IsNotFound(err) {
				resource := &exposedeployv1alpha1.ExposeDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "my-namespace",
					},
					Spec: exposedeployv1alpha1.ExposeDeploymentSpec{
						Replicas:            1,
						Image:               "busybox:latest",
						MinAvailableTimeSec: 5,
						CustomEnv:           []string{"echo hello"},
						PortDefinition:      exposedeployv1alpha1.PortDefinition{Port: 80, TargetPort: 80},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &exposedeployv1alpha1.ExposeDeployment{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ExposeDeployment")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			By("deleting the namespace")
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-namespace",
				},
			}
			Expect(k8sClient.Delete(ctx, namespace)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ExposeDeploymentReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
			By("Listing all pods in namespace and asserting that len is 1")
			pods := &corev1.PodList{}
			err = k8sClient.List(ctx, pods, client.InNamespace("my-namespace"))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(pods.Items)).To(Equal(1))
			Expect(pods.Items[0].Name).To(ContainSubstring(resourceName))
			By("Listing all services in namespace and asserting that len is 1")
			services := &corev1.ServiceList{}
			err = k8sClient.List(ctx, services, client.InNamespace("my-namespace"))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(services.Items)).To(Equal(1))
			Expect(services.Items[0].Name).To(ContainSubstring(resourceName))
		})
	})
})
