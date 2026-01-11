package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// default to $HOME/.kube/config instead of a literal ~ which isn't expanded by Go
	home, _ := os.UserHomeDir()
	defaultKube := filepath.Join(home, ".kube", "config")
	kubeconfig := flag.String("kubeconfig", defaultKube, "absolute path to the kubeconfig file")
	ctx := context.Background()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error %s building config from flag\n", err.Error())
		// try in-cluster config as a fallback. This is useful when running inside a k8s cluster.
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error %s getting in-cluster config\n", err.Error())
			return
		}
	}

	// NewForCOnfig gives typed-client.

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s creating Kubernetes client\n", err.Error())
	}

	// Watch pods in the default namespace
	watcher, err := clientSet.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error %s watching pods\n", err.Error())
		return
	}
	defer watcher.Stop()

	fmt.Println("Watching pods in 'default' namespace...")
	fmt.Println("Press Ctrl+C to stop")

	// Process watch events
	for event := range watcher.ResultChan() {
		pod := event.Object.(*corev1.Pod)
		fmt.Printf("Event: %s - Pod: %s\n", event.Type, pod.Name)
	}
}
