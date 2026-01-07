package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("Can not get homedir %w", err)
	}
	defaultKubeconfigPath := filepath.Join(home, ".kube", "config")
	configFlag := flag.String("kubeconfig",defaultKubeconfigPath,"path of kubecionfig")
	config, err := clientcmd.BuildConfigFromFlags("",*configFlag)
	if err != nil {
		return fmt.Errorf("Can not get config %w", err)
	}
	
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Can not get config %w", err)
	}
	
	resources, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group: "helm.cattle.io",
		Version: "v1",
		Resource: "helmcharts",
	}).Namespace("kube-system").Get(context.Background(), "traefik", metav1.GetOptions{})
	// resources, err := dynamicClient.Resource(schema.GroupVersionResource{
	// 	Group: "helm.cattle.io",
	// 	Version: "v1",
	// 	Resource: "helmcharts",
	// }).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to list resources: %w", err)
	}
	fmt.Println(resources.GetName())
	// fmt.Println(resources)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal("Application failed: ", err)
	}
}