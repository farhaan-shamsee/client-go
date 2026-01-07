package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type controller struct {
	clientSet      kubernetes.Interface
	depLister      appslisters.DeploymentLister
	depCacheSynced cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func newController(clientSet kubernetes.Interface, depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientSet:      clientSet,
		depLister:      depInformer.Lister(),
		depCacheSynced: depInformer.Informer().HasSynced,
		queue:          workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekspose"),
	}
	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handleAdd,
			DeleteFunc: handleDelete,
		})

	return c
}

// actual controller run
func (c *controller) run(ch <-chan struct{}) {
	fmt.Println("Starting controller")
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("waiting for cache to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, ch)
}

func (c *controller) worker() {
	fmt.Println("working...")
}

func handleAdd(obj interface{}) {
	fmt.Println("add was called")
}

func handleDelete(obj interface{}) {
	fmt.Println("delete was called")
}