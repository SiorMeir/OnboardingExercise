# onboardingexercise
As part of Run:AI's onboarding, we need to create a Simple `ExposeDeployment` Kubernetes resource that acts similarly to `Deployment`, but with 
a service `built-in` for network exposure.

## ✅ ExposeDeployment Controller - Development Checklist

### 🚧 Environment Setup

* [ ] Install Go (≥ 1.21)
* [ ] Install `operator-sdk`
* [ ] Install `kubectl`
* [ ] Set up a local Kubernetes cluster (e.g., Minikube)
* [ ] Scaffold the operator project using `operator-sdk`
* [ ] Enable necessary RBAC permissions

---

### 📦 Define Custom Resource Definition (CRD)

* [ ] Define `ExposeDeployment` CRD with the following `spec` fields:

  * [ ] `Image` (non-empty string, immutable)
  * [ ] `Replicas` (max 10)
  * [ ] `MinAvailableTimeSec` (double if odd)
  * [ ] `Args` (command-line args)
  * [ ] `Envs` (environment variables)
  * [ ] `Ports` with:

    * [ ] `TargetPort`
    * [ ] `Port`
* [ ] Define `status` with:

  * [ ] `ReadyPods`
  * [ ] `AvailablePods`
  * [ ] `Conditions` array with fields:

    * [ ] `Type` (LastReconcileSucceeded, Ready, Available)
    * [ ] `Status` (bool)
    * [ ] `Message` (string)
    * [ ] `Reason` (camel case word)
    * [ ] `LastTransitionTime` (timestamp)

---

### 🔁 Controller Logic

* [ ] Create controller that watches `ExposeDeployment` CRs
* [ ] Set up the manager with cache limited to owned pods
* [ ] Reconciliation logic:

  * [ ] Create pods owned by the CR (label: `ExposeDeployment: <CR name>`)
  * [ ] Create non-owned services for exposed ports
  * [ ] Maintain desired replica count
  * [ ] Validate and transform `MinAvailableTimeSec` if needed
  * [ ] Track pod readiness and availability
  * [ ] Update CR `status` and `conditions` accordingly
  * [ ] Ensure reconciliation is idempotent and only affects `spec`-defined fields
  * [ ] On deletion of CR, clean up associated pods and services

---

### 🧪 Testing

* [ ] Test creation and update of `ExposeDeployment`
* [ ] Test pod readiness tracking
* [ ] Test service creation and accessibility
* [ ] Use `kubectl port-forward` and `curl` to verify service-pod connectivity
* [ ] Confirm replicas max constraint and odd `MinAvailableTimeSec` handling
* [ ] Validate cleanup on CR deletion
* [ ] Use `envtest` or controller-runtime test framework

---

### ⚠️ Error Handling & Logging

* [ ] Add structured logging throughout reconciliation
* [ ] Handle and log all errors gracefully
* [ ] Ensure requeues on failure where applicable

---

### 📚 Documentation

* [ ] Document CRD spec and usage
* [ ] Add deployment guide
* [ ] Describe testing steps and troubleshooting tips
* [ ] Mention required tools and prerequisites

---

### 📎 Optional: Advanced Features

* [ ] Admission Webhooks for validation/immutability of fields
* [ ] Readiness probes or lifecycle hooks

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/onboardingexercise:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/onboardingexercise:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/onboardingexercise:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/onboardingexercise/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
operator-sdk edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.
