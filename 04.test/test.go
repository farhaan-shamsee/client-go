package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic("cant find home directory", err)
	}
	configFilePath := filepath.Join(home, ".kube", "config")
	kubeconfig := flag.String("kubeconfig", configFilePath, "path to kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	deployments,err := clientSet.AppsV1().Deployments("kube-system").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, deps := range deployments.Items{
		fmt.Println(deps.Name)
	}
}
