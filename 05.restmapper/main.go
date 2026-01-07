package main

import (
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	// "k8s.io/client-go/restmapper"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func run() error {

	res := "addons"

	configFlag := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionFlags := cmdutil.NewMatchVersionFlags(configFlag)

	m, err := cmdutil.NewFactory(matchVersionFlags).ToRESTMapper()
	if err != nil {
		return fmt.Errorf("Failed to get restMapper: %w", err)
	}

	gvr, err := m.ResourceFor(schema.GroupVersionResource{
		Resource: res,
	})
	if err != nil {
		return fmt.Errorf("Failed to get resource: %w", err)
	}

	fmt.Printf("Complete GVR is, group %s, version %s, resource %s\n", gvr.Group, gvr.Version ,gvr.Resource)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal("Application failed: ", err)
	}
}
