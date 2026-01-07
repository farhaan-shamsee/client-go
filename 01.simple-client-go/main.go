package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

	pods, err := clientSet.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error %s listing pods\n", err.Error())
	}
	fmt.Println("Pods:================")
	for _, pod := range pods.Items {
		fmt.Printf("%v\n", pod.Name)
	}

	deployments, err := clientSet.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error %s listing deployments\n", err.Error())
	}
	fmt.Println("Deployments:===============")
	for _, deploy := range deployments.Items {
		fmt.Printf("%v\n", deploy.Name)
	}

	ingress, err := clientSet.NetworkingV1().Ingresses("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Errorf("cant get ingresses: %w", err)
	}
	for _, ingressList := range ingress.Items {
		fmt.Println(ingressList.Name)
	}
}
