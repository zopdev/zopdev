package provider

import (
	"gofr.dev/pkg/gofr"
)

// Provider defines the interface for interacting with a cloud provider's resources.
// It includes methods for listing all clusters and retrieving namespaces for a given cluster.
//
// This interface can be implemented for various cloud providers such as AWS, GCP, or Azure.
// It allows users to interact with cloud infrastructure, retrieve clusters, and list namespaces.
type Provider interface {
	// ListAllClusters lists all clusters available for a given cloud account.
	//
	// ctx: The context for the request.
	// cloudAccount: The cloud account associated with the provider (e.g., AWS, GCP, Azure).
	// credentials: The authentication credentials used to access the provider's resources.
	//
	// Returns a ClusterResponse containing details of the available clusters, or an error if the request fails.
	ListAllClusters(ctx *gofr.Context, cloudAccount *CloudAccount, credentials interface{}) (*ClusterResponse, error)

	// ListNamespace retrieves namespaces for a given cluster within a cloud account.
	//
	// ctx: The context for the request.
	// cluster: The cluster for which to list namespaces.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	//
	// Returns the namespaces for the specified cluster, or an error if the request fails.
	ListNamespace(ctx *gofr.Context, cluster *Cluster, cloudAccount *CloudAccount, credentials interface{}) (interface{}, error)

	// ListServices lists all services within a specified namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster for which to list services.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace within the cluster to list services.
	//
	// Returns the services for the specified namespace, or an error if the request fails.
	ListServices(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace string) (interface{}, error)

	// ListDeployments lists all deployments within a specified namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster for which to list deployments.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace within the cluster to list deployments.
	//
	// Returns the deployments for the specified namespace, or an error if the request fails.
	ListDeployments(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace string) (interface{}, error)

	// ListPods lists all pods within a specified namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster for which to list pods.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace within the cluster to list pods.
	//
	// Returns the pods for the specified namespace, or an error if the request fails.
	ListPods(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace string) (interface{}, error)

	// ListCronJobs lists all cron jobs within a specified namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster for which to list cron jobs.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace within the cluster to list cron jobs.
	//
	// Returns the cron jobs for the specified namespace, or an error if the request fails.
	ListCronJobs(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace string) (interface{}, error)

	// GetService retrieves details of a specific service within a namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster containing the service.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace containing the service.
	// name: The name of the service to retrieve.
	//
	// Returns the details of the specified service, or an error if the request fails.
	GetService(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace, name string) (interface{}, error)

	// GetDeployment retrieves details of a specific deployment within a namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster containing the deployment.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace containing the deployment.
	// name: The name of the deployment to retrieve.
	//
	// Returns the details of the specified deployment, or an error if the request fails.
	GetDeployment(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace, name string) (interface{}, error)

	// GetPod retrieves details of a specific pod within a namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster containing the pod.
	// cloudAccount: The cloud account associated with the provider.
	// credentials: The authentication credentials used to access the provider's resources.
	// namespace: The namespace containing the pod.
	// name: The name of the pod to retrieve.
	//
	// Returns the details of the specified pod, or an error if the request fails.
	GetPod(ctx *gofr.Context, cluster *Cluster,
		cloudAccount *CloudAccount, credentials interface{}, namespace, name string) (interface{}, error)

	// GetCronJob retrieves details of a specific cron job within a namespace in a cluster.
	//
	// ctx: The context for the request.
	// cluster: The cluster containing the cron job.
	// cloudAcc: The cloud account associated with the provider.
	// creds: The authentication credentials used to access the provider's resources.
	// namespace: The namespace containing the cron job.
	// name: The name of the cron job to retrieve.
	//
	// Returns the details of the specified cron job, or an error if the request fails.
	GetCronJob(ctx *gofr.Context, cluster *Cluster,
		cloudAcc *CloudAccount, creds any, namespace, name string) (any, error)
}
