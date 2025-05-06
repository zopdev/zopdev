package main

import (
	appHandler "github.com/zopdev/zopdev/api/applications/handler"
	appService "github.com/zopdev/zopdev/api/applications/service"
	appStore "github.com/zopdev/zopdev/api/applications/store"
	auditHandler "github.com/zopdev/zopdev/api/audit/handler"
	auditService "github.com/zopdev/zopdev/api/audit/service"
	auditStore "github.com/zopdev/zopdev/api/audit/store"

	caHandler "github.com/zopdev/zopdev/api/cloudaccounts/handler"
	caService "github.com/zopdev/zopdev/api/cloudaccounts/service"
	caStore "github.com/zopdev/zopdev/api/cloudaccounts/store"
	clStore "github.com/zopdev/zopdev/api/deploymentspace/cluster/store"
	"github.com/zopdev/zopdev/api/provider/gcp"

	envHandler "github.com/zopdev/zopdev/api/environments/handler"
	envService "github.com/zopdev/zopdev/api/environments/service"
	envStore "github.com/zopdev/zopdev/api/environments/store"

	deployHandler "github.com/zopdev/zopdev/api/deploymentspace/handler"
	deployService "github.com/zopdev/zopdev/api/deploymentspace/service"
	deployStore "github.com/zopdev/zopdev/api/deploymentspace/store"

	clService "github.com/zopdev/zopdev/api/deploymentspace/cluster/service"

	"github.com/zopdev/zopdev/api/migrations"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	gkeSvc := gcp.New()

	cloudAccountStore := caStore.New()
	cloudAccountService := caService.New(cloudAccountStore, gkeSvc)
	cloudAccountHandler := caHandler.New(cloudAccountService)

	deploymentStore := deployStore.New()
	clusterStore := clStore.New()
	clusterService := clService.New(clusterStore)
	deploymentService := deployService.New(deploymentStore, clusterService, cloudAccountService, gkeSvc)

	environmentStore := envStore.New()
	deploymentHandler := deployHandler.New(deploymentService)

	environmentService := envService.New(environmentStore, deploymentService)
	envrionmentHandler := envHandler.New(environmentService)

	applicationStore := appStore.New()
	applicationService := appService.New(applicationStore, environmentService)
	applicationHandler := appHandler.New(applicationService)

	auditStore := auditStore.New()
	auditService := auditService.New(auditStore)
	auditHandler := auditHandler.New(auditService)

	app.AddHTTPService("cloud-account", "http://localhost:8000")

	app.POST("/cloud-accounts", cloudAccountHandler.AddCloudAccount)
	app.GET("/cloud-accounts", cloudAccountHandler.ListCloudAccounts)
	app.GET("/cloud-accounts/{id}/deployment-space/clusters", cloudAccountHandler.ListDeploymentSpace)
	app.GET("/cloud-accounts/{id}/deployment-space/namespaces", cloudAccountHandler.ListNamespaces)
	app.GET("/cloud-accounts/{id}/deployment-space/options", cloudAccountHandler.ListDeploymentSpaceOptions)
	app.GET("/cloud-accounts/{id}/credentials", cloudAccountHandler.GetCredentials)

	app.POST("/audit/cloud-accounts/{id}/all", auditHandler.RunAll)
	app.POST("/audit/cloud-accounts/{id}/category/{category}", auditHandler.RunByCategory)
	app.POST("/audit/cloud-accounts/{id}/rule/{ruleId}", auditHandler.RunById)

	app.POST("/applications", applicationHandler.AddApplication)
	app.GET("/applications", applicationHandler.ListApplications)
	app.GET("/applications/{id}", applicationHandler.GetApplication)

	app.POST("/applications/{id}/environments", envrionmentHandler.Add)
	app.GET("/applications/{id}/environments", envrionmentHandler.List)
	app.PATCH("/applications/{id}/environments", envrionmentHandler.Update)

	app.POST("/environments/{id}/deploymentspace", deploymentHandler.Add)
	app.GET("/environments/{id}/deploymentspace/service/{name}", deploymentHandler.GetService)
	app.GET("/environments/{id}/deploymentspace/service", deploymentHandler.ListServices)
	app.GET("/environments/{id}/deploymentspace/deployment/{name}", deploymentHandler.GetDeployment)
	app.GET("/environments/{id}/deploymentspace/deployment", deploymentHandler.ListDeployments)
	app.GET("/environments/{id}/deploymentspace/pod/{name}", deploymentHandler.GetPod)
	app.GET("/environments/{id}/deploymentspace/pod", deploymentHandler.ListPods)
	app.GET("/environments/{id}/deploymentspace/cronjob/{name}", deploymentHandler.GetCronJob)
	app.GET("/environments/{id}/deploymentspace/cronjob", deploymentHandler.ListCronJobs)

	app.Run()
}
