package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	home, err := os.UserHomeDir()
	ctx := context.Background()
	if err != nil {
		fmt.Println("Error getting the homedir:", err)
		os.Exit(1)
	}
	kubeconfigPath := filepath.Join(home,".kube","config")
	config, err := clientcmd.BuildConfigFromFlags("",kubeconfigPath)
	if err != nil {
		fmt.Println("Error getting the config:", err)
		os.Exit(1)
	}
	
	m, err := cmdutil.NewFactory(config).ToRESTMapper()
	
}