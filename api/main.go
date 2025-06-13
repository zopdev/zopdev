package main

import (
	"gofr.dev/pkg/gofr"

	appHandler "github.com/zopdev/zopdev/api/applications/handler"
	appService "github.com/zopdev/zopdev/api/applications/service"
	appStore "github.com/zopdev/zopdev/api/applications/store"

	auditHandler "github.com/zopdev/zopdev/api/audit/handler"
	auditService "github.com/zopdev/zopdev/api/audit/service"
	auditStore "github.com/zopdev/zopdev/api/audit/store"

	caHandler "github.com/zopdev/zopdev/api/cloudaccounts/handler"
	caService "github.com/zopdev/zopdev/api/cloudaccounts/service"
	caStore "github.com/zopdev/zopdev/api/cloudaccounts/store"

	clService "github.com/zopdev/zopdev/api/deploymentspace/cluster/service"
	clStore "github.com/zopdev/zopdev/api/deploymentspace/cluster/store"

	deployHandler "github.com/zopdev/zopdev/api/deploymentspace/handler"
	deployService "github.com/zopdev/zopdev/api/deploymentspace/service"
	deployStore "github.com/zopdev/zopdev/api/deploymentspace/store"

	envHandler "github.com/zopdev/zopdev/api/environments/handler"
	envService "github.com/zopdev/zopdev/api/environments/service"
	envStore "github.com/zopdev/zopdev/api/environments/store"
	"github.com/zopdev/zopdev/api/migrations"
	"github.com/zopdev/zopdev/api/provider/gcp"

	resourceClient "github.com/zopdev/zopdev/api/resources/client"
	resourceHandler "github.com/zopdev/zopdev/api/resources/handler/resource"
	resourceService "github.com/zopdev/zopdev/api/resources/service/resource"
	resourceStore "github.com/zopdev/zopdev/api/resources/store/resource"

	"github.com/zopdev/zopdev/api/resources/providers"

	resGroupHandler "github.com/zopdev/zopdev/api/resources/handler/resourcegroup"
	resGroupService "github.com/zopdev/zopdev/api/resources/service/resourcegroup"
	resGroupStore "github.com/zopdev/zopdev/api/resources/store/resourcegroup"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())
	app.Metrics().NewCounter("db_error_count", "Count of DB errors")

	gkeSvc := gcp.New()

	trustedPrincipalARN := app.Config.Get("AWS_TRUSTED_PRINCIPAL_ARN")

	cloudAccountStore := caStore.New()
	cloudAccountService := caService.New(cloudAccountStore, gkeSvc, trustedPrincipalARN)
	cloudAccountHandler := caHandler.New(cloudAccountService)

	deploymentStore := deployStore.New()
	clusterStore := clStore.New()
	clusterService := clService.New(clusterStore)
	deploymentService := deployService.New(deploymentStore, clusterService, cloudAccountService, gkeSvc)

	environmentStore := envStore.New()
	deploymentHandler := deployHandler.New(deploymentService)

	environmentService := envService.New(environmentStore, deploymentService)
	environmentHandler := envHandler.New(environmentService)

	applicationStore := appStore.New()
	applicationService := appService.New(applicationStore, environmentService)
	applicationHandler := appHandler.New(applicationService)

	app.AddHTTPService("cloud-account", "http://localhost:8000")

	app.POST("/cloud-accounts", cloudAccountHandler.AddCloudAccount)
	app.GET("/cloud-accounts", cloudAccountHandler.ListCloudAccounts)
	app.GET("/cloud-accounts/{id}/deployment-space/clusters", cloudAccountHandler.ListDeploymentSpace)
	app.GET("/cloud-accounts/{id}/deployment-space/namespaces", cloudAccountHandler.ListNamespaces)
	app.GET("/cloud-accounts/{id}/deployment-space/options", cloudAccountHandler.ListDeploymentSpaceOptions)
	app.GET("/cloud-accounts/{id}/credentials", cloudAccountHandler.GetCredentials)

	// Endpoints for one-click setup
	app.GET("/cloud-accounts/{provider}/connection/info", cloudAccountHandler.GetCloudAccountConnectionInfo)
	app.POST("/cloud-accounts/{provider}/connection", cloudAccountHandler.CreateCloudAccountConnection)
	app.GET("/cloud-accounts/{provider}/stack-status/{integrationId}", cloudAccountHandler.GetStackStatus)

	app.POST("/applications", applicationHandler.AddApplication)
	app.GET("/applications", applicationHandler.ListApplications)
	app.GET("/applications/{id}", applicationHandler.GetApplication)

	app.POST("/applications/{id}/environments", environmentHandler.Add)
	app.GET("/applications/{id}/environments", environmentHandler.List)
	app.PATCH("/applications/{id}/environments", environmentHandler.Update)

	app.POST("/environments/{id}/deploymentspace", deploymentHandler.Add)
	app.GET("/environments/{id}/deploymentspace/service/{name}", deploymentHandler.GetService)
	app.GET("/environments/{id}/deploymentspace/service", deploymentHandler.ListServices)
	app.GET("/environments/{id}/deploymentspace/deployment/{name}", deploymentHandler.GetDeployment)
	app.GET("/environments/{id}/deploymentspace/deployment", deploymentHandler.ListDeployments)
	app.GET("/environments/{id}/deploymentspace/pod/{name}", deploymentHandler.GetPod)
	app.GET("/environments/{id}/deploymentspace/pod", deploymentHandler.ListPods)
	app.GET("/environments/{id}/deploymentspace/cronjob/{name}", deploymentHandler.GetCronJob)
	app.GET("/environments/{id}/deploymentspace/cronjob", deploymentHandler.ListCronJobs)

	registerAuditAPIRoutes(app)
	registerCloudResourceRoutes(app)

	app.Run()
}

func registerAuditAPIRoutes(app *gofr.App) {
	adStore := auditStore.New()
	adSvc := auditService.New(adStore)
	adHandler := auditHandler.New(adSvc)

	app.POST("/audit/cloud-accounts/{id}/all", adHandler.RunAll)
	app.POST("/audit/cloud-accounts/{id}/category/{category}", adHandler.RunByCategory)
	app.POST("/audit/cloud-accounts/{id}/rule/{ruleId}", adHandler.RunByID)
	app.GET("/audit/cloud-accounts/{id}/results", adHandler.GetAllResults)
	app.GET("/audit/cloud-accounts/{id}/results/{ruleId}", adHandler.GetResultByID)
}

func registerCloudResourceRoutes(app *gofr.App) {
	client := resourceClient.New()

	gcpClient, err := providers.NewCloudResourceProvider("GCP")
	if err != nil {
		app.Logger().Fatal(err)
	}

	awsClient, err := providers.NewCloudResourceProvider("AWS")
	if err != nil {
		app.Logger().Fatal(err)
	}

	resStore := resourceStore.New()
	resSvc := resourceService.New(gcpClient, awsClient, client, resStore)
	resHld := resourceHandler.New(resSvc)

	// TODO: Figure out a way to sync resources on startup.

	app.AddCronJob("0 * * * *", "resource-sync", resSvc.SyncCron)

	app.GET("/cloud-account/{id}/resources", resHld.GetResources)
	app.POST("/cloud-account/{id}/resources/state", resHld.ChangeState)
	app.POST("/cloud-account/{id}/resources/sync", resHld.SyncResources)

	rgStr := resGroupStore.New()
	rgSvc := resGroupService.New(rgStr, resSvc)
	rgHld := resGroupHandler.New(rgSvc)

	app.GET("/cloud-account/{id}/resource-groups", rgHld.GetAllResourceGroups)
	app.GET("/cloud-account/{id}/resource-groups/{rgID}", rgHld.GetResourceGroup)
	app.POST("/cloud-account/{id}/resource-groups", rgHld.CreateResourceGroup)
	app.PUT("/cloud-account/{id}/resource-groups/{rgID}", rgHld.UpdateResourceGroup)
	app.DELETE("/cloud-account/{id}/resource-groups/{rgID}", rgHld.DeleteResourceGroup)
}
