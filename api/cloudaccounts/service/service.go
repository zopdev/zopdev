package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/cloudaccounts/store"
	"github.com/zopdev/zopdev/api/provider"
)

type Service struct {
	store           store.CloudAccountStore
	deploymentSpace provider.Provider
}

// New creates a new CloudAccountService with the provided CloudAccountStore.
func New(clStore store.CloudAccountStore, deploySpace provider.Provider) CloudAccountService {
	return &Service{store: clStore, deploymentSpace: deploySpace}
}

// AddCloudAccount adds a new cloud account to the store if it doesn't already exist.
func (s *Service) AddCloudAccount(ctx *gofr.Context, cloudAccount *store.CloudAccount) (*store.CloudAccount, error) {
	// TODO : validation is only checking if the values are present - we also need to check if the values are valid
	// and able to connect to a cloud account, we would need to keep that code in a separate package where all gcp, aws code is present.
	switch strings.ToUpper(cloudAccount.Provider) {
	case gcp:
		err := fetchGCPProviderDetails(ctx, cloudAccount)
		if err != nil {
			return nil, err
		}
	case aws:
		err := validateAWSProviderDetails(ctx, cloudAccount)
		if err != nil {
			return nil, err
		}
	case oci:
		err := fetchOCIProviderDetails(ctx, cloudAccount)
		if err != nil {
			return nil, err
		}
	default:
		return nil, http.ErrorInvalidParam{Params: []string{"provider"}}
	}

	tempCloudAccount, err := s.store.GetCloudAccountByProvider(ctx, cloudAccount.Provider, cloudAccount.ProviderID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if tempCloudAccount != nil {
		return nil, http.ErrorEntityAlreadyExist{}
	}

	cloudAccount.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	return s.store.InsertCloudAccount(ctx, cloudAccount)
}

// FetchAllCloudAccounts retrieves all cloud accounts from the store.
func (s *Service) FetchAllCloudAccounts(ctx *gofr.Context) ([]store.CloudAccount, error) {
	return s.store.GetALLCloudAccounts(ctx)
}

// fetchGCPProviderDetails retrieves and assigns GCP details for a cloud account.
func fetchGCPProviderDetails(ctx *gofr.Context, cloudAccount *store.CloudAccount) error {
	var gcpCred gcpCredentials

	body, err := json.Marshal(cloudAccount.Credentials)
	if err != nil {
		ctx.Error(err.Error())
		return http.ErrorInvalidParam{}
	}

	err = json.Unmarshal(body, &gcpCred)
	if err != nil {
		return err
	}

	if gcpCred.ProjectID == "" {
		return http.ErrorInvalidParam{Params: []string{"credentials"}}
	}

	cloudAccount.ProviderID = gcpCred.ProjectID

	return nil
}

func validateAWSProviderDetails(_ *gofr.Context, account *store.CloudAccount) error {
	var awsCred awsCredentials

	awsCredBody, _ := json.Marshal(account.Credentials)

	err := json.Unmarshal(awsCredBody, &awsCred)
	if err != nil {
		return err
	}

	if awsCred.AccessKey == "" && awsCred.AccessSecret == "" {
		return http.ErrorMissingParam{Params: []string{"AWSAccessKeyID", "AWSecretAccessKey"}}
	}

	if awsCred.AccessKey == "" {
		return http.ErrorMissingParam{Params: []string{"AWSAccessKeyID"}}
	}

	if awsCred.AccessSecret == "" {
		return http.ErrorMissingParam{Params: []string{"AWSecretAccessKey"}}
	}

	account.ProviderID, err = getAWSAccountID(awsCred)
	if err != nil {
		return err
	}

	return nil
}

// TODO this logic will be strictly moved after the resource pause PR is merged.
func getAWSAccountID(awsCred awsCredentials) (string, error) {
	sess, err := session.NewSession(&awsSDK.Config{
		Credentials: credentials.NewStaticCredentials(awsCred.AccessKey, awsCred.AccessSecret, "")})
	if err != nil {
		return "", err
	}

	// Create an STS client
	svc := sts.New(sess)

	// Get the caller identity
	resp, err := svc.GetCallerIdentity(nil)
	if err != nil {
		return "", err
	}

	return *resp.Account, nil
}

// fetchOCIProviderDetails retrieves and assigns OCI details for a cloud account.
func fetchOCIProviderDetails(ctx *gofr.Context, cloudAccount *store.CloudAccount) error {
	var ociCred ociCredentials

	body, err := json.Marshal(cloudAccount.Credentials)
	if err != nil {
		ctx.Error(err.Error())
		return http.ErrorInvalidParam{}
	}

	err = json.Unmarshal(body, &ociCred)
	if err != nil {
		return err
	}

	if ociCred.TenancyOCID == "" || ociCred.UserOCID == "" {
		return http.ErrorInvalidParam{Params: []string{"credentials"}}
	}

	// For OCI, we'll use the tenancy OCID as the provider ID since it's unique
	cloudAccount.ProviderID = ociCred.TenancyOCID

	return nil
}

func (s *Service) FetchDeploymentSpace(ctx *gofr.Context, cloudAccountID int64) (interface{}, error) {
	cloudAccount, err := s.store.GetCloudAccountByID(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	creds, err := s.store.GetCredentials(ctx, cloudAccount.ID)
	if err != nil {
		return nil, err
	}

	deploymentSpaceAccount := provider.CloudAccount{
		ID:              cloudAccount.ID,
		Name:            cloudAccount.Name,
		Provider:        cloudAccount.Provider,
		ProviderID:      cloudAccount.ProviderID,
		ProviderDetails: cloudAccount.ProviderDetails,
	}

	clusters, err := s.deploymentSpace.ListAllClusters(ctx, &deploymentSpaceAccount, creds)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (s *Service) ListNamespaces(ctx *gofr.Context, id int64, clusterName, clusterRegion string) (interface{}, error) {
	cloudAccount, err := s.store.GetCloudAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	creds, err := s.store.GetCredentials(ctx, cloudAccount.ID)
	if err != nil {
		return nil, err
	}

	deploymentSpaceAccount := provider.CloudAccount{
		ID:              cloudAccount.ID,
		Name:            cloudAccount.Name,
		Provider:        cloudAccount.Provider,
		ProviderID:      cloudAccount.ProviderID,
		ProviderDetails: cloudAccount.ProviderDetails,
	}

	cluster := provider.Cluster{
		Name:   clusterName,
		Region: clusterRegion,
	}

	res, err := s.deploymentSpace.ListNamespace(ctx, &cluster, &deploymentSpaceAccount, creds)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (*Service) FetchDeploymentSpaceOptions(_ *gofr.Context, id int64) ([]DeploymentSpaceOptions, error) {
	options := []DeploymentSpaceOptions{
		{
			Name: "gke",
			Path: fmt.Sprintf("/cloud-accounts/%v/deployment-space/clusters", id),
			Type: "type",
		},
		{
			Name: "oke",
			Path: fmt.Sprintf("/cloud-accounts/%v/deployment-space/clusters", id),
			Type: "type",
		},
	}

	return options, nil
}

func (s *Service) FetchCredentials(ctx *gofr.Context, cloudAccountID int64) (interface{}, error) {
	creds, err := s.store.GetCredentials(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	cloudAcc, err := s.store.GetCloudAccountByID(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	cloudAcc.Credentials = creds

	return cloudAcc, nil
}
