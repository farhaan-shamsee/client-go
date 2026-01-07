package main

import (
	// "context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	// "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// default to $HOME/.kube/config instead of a literal ~ which isn't expanded by Go
	home, _ := os.UserHomeDir()
	defaultKube := filepath.Join(home, ".kube", "config")
	kubeconfig := flag.String("kubeconfig", defaultKube, "absolute path to the kubeconfig file")
	// ctx := context.Background()
	
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
	
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s creating Kubernetes client\n", err.Error())
	}

	ch := make(chan struct{})
	informers := informers.NewSharedInformerFactory(clientSet, 10*time.Minute)
	c := newController(clientSet, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.run(ch)
	fmt.Println(informers)
}
