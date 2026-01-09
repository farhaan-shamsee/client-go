package main

import (
	"fmt" // For formatted I/O operations like Printf
	"log" // For logging errors and fatal messages

	"k8s.io/apimachinery/pkg/runtime/schema" // Provides schema types like GroupVersionResource for Kubernetes API objects
	"k8s.io/cli-runtime/pkg/genericclioptions" // Provides common CLI flags and configuration options for kubectl-like tools
	// "k8s.io/client-go/restmapper" // Commented out - would provide RESTMapper utilities for resource mapping
	cmdutil "k8s.io/kubectl/pkg/cmd/util" // Provides utility functions from kubectl commands, including Factory for creating clients and mappers
)

func run() error {

	// Define the resource name we want to look up in the Kubernetes API
	res := "addons"

	// Create configuration flags for connecting to the Kubernetes cluster
	// Uses default kubeconfig path and enables deprecated password authentication
	configFlag := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	
	// Create match version flags to ensure API version compatibility
	matchVersionFlags := cmdutil.NewMatchVersionFlags(configFlag)

	// Create a RESTMapper from the factory to discover and map Kubernetes resources
	// The RESTMapper translates between resource names and their full GroupVersionResource
	m, err := cmdutil.NewFactory(matchVersionFlags).ToRESTMapper()
	if err != nil {
		return fmt.Errorf("Failed to get restMapper: %w", err)
	}

	// Look up the complete GroupVersionResource for the given resource name
	// This queries the API server to find the full group, version, and resource details
	gvr, err := m.ResourceFor(schema.GroupVersionResource{
		Resource: res,
	})
	if err != nil {
		return fmt.Errorf("Failed to get resource: %w", err)
	}

	// Print the complete GroupVersionResource information
	fmt.Printf("Complete GVR is, group %s, version %s, resource %s\n", gvr.Group, gvr.Version, gvr.Resource)

	return nil
}

func main() {
	// Execute the run function and log any errors that occur
	if err := run(); err != nil {
		log.Fatal("Application failed: ", err)
	}
}
