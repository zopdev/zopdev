// Package provider contains types and responses for interacting with cloud providers such as AWS, GCP, and Azure.
//
// It provides data structures representing clusters, node pools, namespaces, and cloud accounts, along with their details.
//
// Example usage:
//   - Retrieve cloud account details for AWS, GCP, or Azure.
//   - Fetch clusters and their associated node pools and namespaces in the cloud.
package provider

// ClusterResponse represents the response containing information about clusters.
// It includes a list of clusters and information about pagination.
type ClusterResponse struct {
	// Clusters is a list of clusters available for the provider.
	Clusters []Cluster `json:"options"`

	// Next contains pagination information for retrieving the next set of resources.
	Next Next `json:"next"`

	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Next provides pagination details for fetching additional data.
// It contains the name, path, and parameters required to get the next set of results.
type Next struct {
	// Name is the name of the next .
	Name string `json:"name"`

	// Path is the URL path to the next set of results.
	Path string `json:"path"`

	// Params holds the parameters required to fetch the next set.
	Params map[string]string `json:"params"`
}

// Cluster represents a cloud provider cluster, including details like its name,
// identifier, locations, region, node pools, and type.
type Cluster struct {
	// Name is the name of the cluster.
	Name string `json:"name"`

	// Identifier is a unique identifier for the cluster.
	Identifier string `json:"identifier"`

	// Locations lists the locations available for the cluster.
	Locations []string `json:"locations"`

	// to set key for sending response.
	Type string `json:"type"`

	// Region specifies the region where the cluster is located.
	Region string `json:"region"`

	// NodePools is a list of node pools associated with the cluster.
	NodePools []NodePool `json:"nodePools"`
}

// NodePool represents a node pool within a cluster, detailing machine type, availability zones,
// node version, current node count, and node name.
type NodePool struct {
	// MachineType specifies the machine type for the node pool.
	MachineType string `json:"machineType"`

	// NodeVersion indicates the version of the nodes in the pool.
	NodeVersion string `json:"nodeVersion,omitempty"`

	// NodeName is the name of the node pool.
	NodeName string `json:"nodeName"`

	// CurrentNode specifies the number of nodes currently in the node pool.
	CurrentNode int32 `json:"currentNode"`

	// AvailabilityZones lists the availability zones where nodes in the pool are located.
	AvailabilityZones []string `json:"availabilityZones"`
}

// CloudAccount represents a cloud account, including details such as its name,
// provider, provider-specific ID, provider details, and credentials.
type CloudAccount struct {
	// ID is a unique identifier for the cloud account.
	ID int64 `json:"id"`

	// Name is the name of the cloud account.
	Name string `json:"name"`

	// Provider is the name of the cloud service provider (e.g., AWS, GCP, Azure).
	Provider string `json:"provider"`

	// ProviderID is the unique identifier for the provider account.
	ProviderID string `json:"providerId"`

	// ProviderDetails contains additional details specific to the provider,
	// such as API keys or other configuration settings.
	ProviderDetails interface{} `json:"providerDetails"`

	// Credentials holds authentication information used to access the cloud provider.
	Credentials interface{} `json:"credentials,omitempty"`
}

// NamespaceResponse represents a response containing a list of namespaces available in a cloud provider.
// It includes an array of namespaces.
type NamespaceResponse struct {
	// Options is a list of available namespaces.
	Options []Namespace `json:"options"`

	Metadata Metadata `json:"metadata"`
}

// Namespace represents a namespace within a cloud provider. It contains the name and type of the namespace.
type Namespace struct {
	// Name is the name of the namespace.
	Name string `json:"name"`

	// to set key for sending response.
	Type string `json:"type"`
}

type ServiceResponse struct {
	Services []Service `json:"services"`
	Metadata Metadata  `json:"metadata"`
}

type ServiceList struct {
	Items []Service `json:"items"`
}

type Service struct {
	Metadata K8sMetadata `json:"metadata"`
	Spec     ServiceSpec `json:"spec"`
	Status   Status      `json:"status"`
}

type K8sMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

type ServiceSpec struct {
	Ports                    []Port            `json:"ports"`
	Selector                 map[string]string `json:"selector"`
	ClusterIP                string            `json:"clusterIP"`
	Type                     string            `json:"type"`
	SessionAffinity          string            `json:"sessionAffinity"`
	ExternalIPs              []string          `json:"externalIPs"`
	LoadBalancerIP           string            `json:"loadBalancerIP"`
	LoadBalancerSourceRanges []string          `json:"loadBalancerSourceRanges"`
}

type Port struct {
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort any    `json:"targetPort"`
	NodePort   any    `json:"nodePort"`
}

type Status struct {
	LoadBalancer LoadBalancer `json:"loadBalancer"`
}

type LoadBalancer struct {
	Ingress []Ingress `json:"ingress"`
}

type Ingress struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}
