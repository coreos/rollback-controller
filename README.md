# Auto Rollback Controller

A Kubernetes controller example for learning purposes.

## Usage

In v1.5 Kubernetes added a new field to the `Deployments` spec, [`progressDeadlineSeconds`][rollback-config]. As of `v1.5` a `Deployment` will be marked as failed if it fails to make progress within the allotted time. However, it won't automatically roll back.

This means the following deployment:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: hello
spec:
  progressDeadlineSeconds: 5
  replicas: 3
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
      - name: hello
        image: alpine:3.5
        command:
        - /bin/sh
        - -c
        - "while :; do echo 'Goodbye'; exit 1; sleep 1; done"
```

doesn't roll back in v1.5.

The `rollback-controller` loops, looking for failed deployments, then automatically does this rollback.

## Example

In one terminal, start the rollback controller:

```
$ go get github.com/coreos/rollback-controller/cmd/rollback-controller
$ rollback-controller --kubeconfig=$PATH_TO_KUBECONFIG --namespace=default
```

In another create a deployment, then roll to a bad version of the deployment:

```
$ kubectl create -f examples/good.yaml
$ # Wait a bit for the deployment to succeed...
$ kubectl replace -f examples/bad.yaml
```

The second deployment should roll back to the first.

## Exercises

Once you get the rollback controller working, fork the project and try adding the following features:

### Opt-in annotation

Rather then rolling back all deployments in the cluster, modify the controller to only rollback deployments with a specific annotation. This is a common pattern for new controller features, letting users slowly opt into new functionality before making it the default.

An example annotation:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: hello
  annotations:
    coreos.com/auto-rollback: true
```

The rollback controller should ignore all other deployments.

### Run the controller in the cluster

Containerize the rollback controller using a [Dockerfile][dockerfile], push it to an image registry (e.g. [Quay][quay]), and write a [deployment manifest][deployments] for this app. Deploy the controller as a pod in the cluster instead of running it from your laptop.

### Optimizing performance

Instead of looping and listing, use client-go's [imformers framework][informers] to optimize the controllers performance using the watch API.

[rollback-config]: https://github.com/kubernetes/kubernetes/blob/v1.5.0/pkg/apis/extensions/v1beta1/types.go#L292-L303
[dockerfile]: https://docs.docker.com/engine/reference/builder/
[quay]: https://quay.io/
[deployments]: https://kubernetes.io/docs/user-guide/deployments/
[informers]: https://godoc.org/github.com/kubernetes/client-go/informers
