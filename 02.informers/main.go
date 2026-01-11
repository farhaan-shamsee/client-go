package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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
	
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s creating Kubernetes client\n", err.Error())
	}

	// 30*time.Second is the resync period. This means that every 30 seconds, the informer will re-list all resources and update its cache from API server.
	informerFactory := informers.NewSharedInformerFactory(clientSet, 30*time.Second)

	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add was called")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("update was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete was called")
		},
	})

	// Start the informer factory
	informerFactory.Start(ctx.Done())
	// Wait for the caches in the in-memory cache to be synced before using the informer
	informerFactory.WaitForCacheSync(ctx.Done())

	// Lister is provided by the informer to list resources from the local cache.
	pod, err := podInformer.Lister().Pods("kube-system").Get("coredns-ccb96694c-ks6vz")
	if err != nil {
		fmt.Printf("Error %s getting pod\n", err.Error())
		return
	}
	fmt.Println(pod)
}
