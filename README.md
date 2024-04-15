# AO-Autoscaler
A Kubernetes Operator designed to automate the scaling of a deployment on a cluster.
## Description
The operator is configured to scale based on the following metrics; Minimum Replicas, Maximum Replicas, CPU Utilisation and Manual Replicas Override. Use an instance of the 'Autoscaler' CR to automate your scaling requirements.
## Getting Started
You’ll need a Kubernetes or Openshift cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

## Autoscaler Custom Resource

`TargetDeploymentName:` Is the name of the Deployment that the Autoscaler will manage. This field is not optional

`MinReplicas:` (Optional) Is the minimum number of replicas that the Autoscaler can scale down to. This field is optional.

`MaxReplicas:` Is the maximum number of replicas that the Autoscaler can scale up to. This field is not optional.

`TargetCPUUtilizationPercentage:` (Optional) Is the target average CPU utilization (as a percentage) over all of the pods. If the average CPU utilization exceeds this threshold, the Autoscaler will scale up. This field is optional

`ManualReplicasOverride:` (Optional) Is used to manually set the number of desired pods. If set, this will supersede the other replica fields. This field is optional.

`Autoscaler Sample CR` can be viewed [Here](config/samples/scaling_v1alpha1_autoscaler.yaml)

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/ao-autoscaler:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/ao-autoscaler:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
Create a PR for review if you would like to contribute to my project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

