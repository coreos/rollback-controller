package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/coreos/rollback-controller/rollback"
)

func main() {
	var (
		kubeconfigPath string
		namespace      string
	)

	// client-go writes the the default flag's flagset, adding a bunch of flags
	// that aren't related to this controller. Create a new flagset so only are
	// flags are displayed and settable.
	fs := flag.NewFlagSet("rollback-controller", flag.ExitOnError)
	fs.StringVar(&kubeconfigPath, "kubeconfig", "",
		"Path to kubeconfig. If not provided, the controller assumes it's running in the cluster.")
	fs.StringVar(&namespace, "namepsace", "",
		"Namespace to monitor deployments in. If not provided, defaults to cluster wide.")
	fs.Parse(os.Args[1:])

	var (
		client *kubernetes.Clientset
		err    error
	)

	if kubeconfigPath == "" {
		client, err = newInClusterClient()
	} else {
		client, err = newClient(kubeconfigPath)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	controller := rollback.Controller{Client: client, Namespace: namespace}
	controller.Run(context.Background())
}

// newInClusterClient constructs a client using the well known environment variables
// and credentials configured in a Kubernetes pod.
func newInClusterClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("loading rest config: %v", err)
	}
	return kubernetes.NewForConfig(config)
}

func newClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("loading kubeconfig: %v", err)
	}
	return kubernetes.NewForConfig(config)
}
