# Auto Rollback Controller

A Kubernetes controller example

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

[rollback-config]: https://github.com/kubernetes/kubernetes/blob/v1.5.0/pkg/apis/extensions/v1beta1/types.go#L292-L303
