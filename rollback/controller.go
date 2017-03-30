package rollback

import (
	"context"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// Controller is the rollback controller.
type Controller struct {
	// Client is
	Client *kubernetes.Clientset
	// If provided, the controller will limit.
	Namespace string
}

// Run starts the event loop of the controller.
func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
		case <-time.After(10 * time.Second):
			c.rollbackDeployments()
		}
	}
}

// rollbackDeployments identifies deployments that have failed to make
// progress and rolls them back to the last revision.
func (c *Controller) rollbackDeployments() {
	client := c.Client.ExtensionsV1beta1().Deployments(c.Namespace)

	list, err := client.List(v1.ListOptions{})
	if err != nil {
		glog.Infof("failed to list deployments: %v", err)
		return
	}
	glog.Infof("found %d deployments", len(list.Items))

	for _, d := range list.Items {
		if deploymentFailed(d) && d.Spec.RollbackTo == nil {
			d.Spec.RollbackTo = &v1beta1.RollbackConfig{Revision: 0}
			if _, err := client.Update(&d); err != nil {
				glog.Infof("failed to update deployment %s/%s", d.Name, d.Namespace)
			} else {
				glog.Infof("rolled back deployment %s/%s", d.Name, d.Namespace)
			}
		}
	}
}

// deploymentFailed determines if a deployment gone over its progress deadline.
func deploymentFailed(d v1beta1.Deployment) bool {
	for _, c := range d.Status.Conditions {
		// https://kubernetes.io/docs/user-guide/deployments/#failed-deployment
		if c.Type == v1beta1.DeploymentProgressing &&
			c.Status == "False" &&
			c.Reason == "ProgressDeadlineExceeded" {

			return true
		}
	}
	return false
}
