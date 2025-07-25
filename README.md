# onboardingexercise
As part of Run:AI's onboarding, we need to create a Simple `ExposeDeployment` Kubernetes resource that acts similarly to `Deployment`, but with 
a service `built-in` for network exposure.

Exercise description [here](https://runai.atlassian.net/wiki/spaces/EN/pages/2698805322/ExposeDeployment+k8s+Controller+exerecise)

## ✅ ExposeDeployment Controller - Development Checklist

### 🚧 Environment Setup

* [X] Install Go (≥ 1.21)
* [X] Install `operator-sdk`
* [X] Install `kubectl`
* [X] Set up a local Kubernetes cluster with Kind
* [X] Scaffold the operator project using `operator-sdk`
* [X] Enable necessary RBAC permissions

---

### 📦 Define Custom Resource Definition (CRD)

* [X] Define `ExposeDeployment` CRD with the following `spec` fields:

  * [X] `Image` (non-empty string, immutable)
  * [X] `Replicas` (max 10)
  * [X] `MinAvailableTimeSec` (double if odd)
  * [X] `Args` (command-line args)
  * [X] `Envs` (environment variables)

  * [X] `Ports` with:
    * [X] `TargetPort`
    * [X] `Port`

* [X] Define `status` with:
  * [X] `ReadyPods`
  * [X] `AvailablePods`
  * [X] `Conditions` array with fields:

    * [X] `Type` (LastReconcileSucceeded, Ready, Available)
    * [X] `Status` (bool)
    * [X] `Message` (string)
    * [X] `Reason` (camel case word)
    * [X] `LastTransitionTime` (timestamp)

---

### 🔁 Controller Logic

* [X] Create controller that watches `ExposeDeployment` CRs
* [X] Set up the manager with cache limited to owned pods
* [ ] Reconciliation logic:

  * [X] Create pods owned by the CR (label: `ExposeDeployment: <CR name>`)
  * [X] Create non-owned services for exposed ports
  * [X] Maintain desired replica count
  * [X] Validate and transform `MinAvailableTimeSec` if needed
  * [ ] Track pod readiness and availability
  * [ ] Update CR `status` and `conditions` accordingly
  * [X] Ensure reconciliation is idempotent and only affects `spec`-defined fields
  * [X] On deletion of CR, clean up associated pods and services

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
