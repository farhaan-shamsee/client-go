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
