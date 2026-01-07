# Simple Client Go

- If any go struct implements the `runtime.Object` interface from `k8s.io/apimachinery`, then we can say that the struct is a Kubernetes object.
- It can be stored in etcd, listed/watched, created/updated/deleted via Kubernetes API.
- Pod implements the typeMeta and objectMeta fields from runtime.Object.

```go
type Pod struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   PodSpec   `json:"spec,omitempty"`
    Status PodStatus `json:"status,omitempty"`
}
```

- Embeds `metav1.TypeMeta` — provides `APIVersion` and `Kind` fields.
- Embeds `metav1.ObjectMeta` — provides metadata fields like `name`, `namespace`, `labels`, `annotations`, etc.

- Pod implements the deepCopyObject method from runtime.Object.
  
```go
func (in *Pod) DeepCopyObject() runtime.Object {
    if c := in.DeepCopy(); c != nil {
        return c
    }
    return nil
}
```

- `NewForConfig()` gives us `Clientset`, meaning it gives us clients for all Kubernetes resources, exception being CRDs.
- NewForConfig() -> Clientset -> AppsV1() -> DeploymentInterface.
- The DeploymentInterface has methods like List(), Get(), Create(), Update(), Delete() to interact with Deployment resources.
- The Watch() method on DeploymentInterface returns a `watch.Interface` that has `ResultChan()` that allows us to monitor changes to Deployment resources in real-time.

Reference: [YouTube Video](https://youtu.be/2s_dOZB7ebo?si=OQPIQipBiAKSKIZx)

## Typed vs Dynamic Client

- Typed clients are generated from the OpenAPI specification and provide a strongly typed interface for interacting with Kubernetes resources.
- Dynamic clients use the unstructured API and provide a more flexible, but less type-safe, way to interact with Kubernetes resources.
- `clientSet, err := kubernetes.NewForConfig(config)` this gives us typed-client.

```go
pods, err := clientSet.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
```

- `dynamicClient, err := dynamic.NewForConfig(config)` this gives us dynamic-client.

    ```go
    resources, err := dynamicClient.Resource(schema.GroupVersionResource{
        Group: "helm.cattle.io",
        Version: "v1",
        Resource: "helmcharts",
    }).List(context.Background(), metav1.ListOptions{})
    ```

## Informers

- <img width="1534" height="542" alt="image" src="https://github.com/user-attachments/assets/66bf11d9-65b8-4df8-b1b4-4a8fd64b3b71" />

- Informers provide a high-level API to watch and cache Kubernetes resources.
- `Watch()` method gives us a watch. Interface, which provides a channel to receive events. But using it to much will increase load on the API server.
- Informers use List-Watch mechanism to efficiently monitor resources.
- Informers maintain a local cache of resources, reducing the load on the API server.
- Informers provide event handlers (AddFunc, UpdateFunc, DeleteFunc) to react to changes in resources.
- Informers automatically handle reconnections and resyncs, making them more robust than using `Watch()` directly.
- Shared Informers allow multiple components to share the same cache and watch connection, further reducing load on the API server.
- Resource versions are used by informers to ensure they have the latest state of resources and to handle updates correctly.
- We should never update the objects received in event handlers directly, as they are from the informer's cache. Instead, we should use the `clientset` `DeepCopy` to update resources in the API server.
- `NewFilteredSharedInformerFactory()` allows us to create informers with custom list options, like filtering by label selectors.

## Queues

- <img width="1534" height="542" alt="image" src="https://github.com/user-attachments/assets/1a28877e-ef6d-486c-93a5-6c43004426fa" />
- Informers use work queues to manage the processing of events.
- When an event occurs (add, update, delete), the informer adds a key (usually namespace/name) to the work queue.
- A separate worker goroutine processes items from the work queue.
- The worker retrieves the resource from the informer's cache using the key and performs the necessary processing.
- Using a work queue allows for rate limiting, retries, and batching of events, improving efficiency and reliability.
- 

## Controller Development

- <img width="1265" height="511" alt="image" src="https://github.com/user-attachments/assets/b343542f-2635-44ce-a41f-bfecc0bb305a" />

- We create a channel in controller to signal when to stop the controller. And we make it run until we receive a signal on that channel. If we don't pass any signal in the channel, the controller will run indefinitely.


## RestMapper

- restmapper helps us map GroupVersionKinds (GVKs) to the appropriate REST endpoints in the Kubernetes API.
- It is used by dynamic clients and other components that need to interact with Kubernetes resources without knowing their exact API paths.
- Package: `k8s.io/apimachinery/pkg/api/meta/restmapper`

## API Machinery

- `GVK`: Group, Version, Kind
- Deployment is the Kind.
- `GVR`: Group, Version, Resource
- Resource: plural, lowercase form of Kind.
- deployments is the Resource for Deployment Kind.
- RESTMapper maps GVKs to GVRs and vice versa.

## Scheme

- The scheme defines the structure and validation rules for API objects.
- It is used to ensure that objects conform to the expected format and contain all required fields.
- Package: `k8s.io/apimachinery/pkg/runtime/schema`
- We can use ObjectKinds() method from Scheme to get GVKs for a given object.

    ```go
    gvks, _, err := scheme.ObjectKinds(obj)
    if err != nil {
        return nil, err
    }
    if len(gvks) == 0 {
        return nil, fmt.Errorf("no GVK found for object")
    }
    return &gvks[0], nil
    ```

- This will only work if the object is a registered type in the scheme. We can use AddKnownTypes() method to register custom types.

SO the path is: Go struct -> Use AddKnownTypes() to register the struct in Scheme -> Now we can use ObjectKinds() to get GVK for the struct -> Use RESTMapper to map GVK to GVR -> Use GVR to interact with the resource via dynamic client.