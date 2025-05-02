package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

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
	//nolint:gocritic //addition of more providers
	switch strings.ToUpper(cloudAccount.Provider) {
	case gcp:
		err := fetchGCPProviderDetails(ctx, cloudAccount)
		if err != nil {
			return nil, err
		}
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
		ctx.Logger.Error(err.Error())
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

func (s *Service) FetchDeploymentSpace(ctx *gofr.Context, cloudAccountID int) (interface{}, error) {
	cloudAccount, err := s.store.GetCloudAccountByID(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	credentials, err := s.store.GetCredentials(ctx, cloudAccount.ID)
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

	clusters, err := s.deploymentSpace.ListAllClusters(ctx, &deploymentSpaceAccount, credentials)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (s *Service) ListNamespaces(ctx *gofr.Context, id int, clusterName, clusterRegion string) (interface{}, error) {
	cloudAccount, err := s.store.GetCloudAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	credentials, err := s.store.GetCredentials(ctx, cloudAccount.ID)
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

	res, err := s.deploymentSpace.ListNamespace(ctx, &cluster, &deploymentSpaceAccount, credentials)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (*Service) FetchDeploymentSpaceOptions(_ *gofr.Context, id int) ([]DeploymentSpaceOptions, error) {
	options := []DeploymentSpaceOptions{
		{
			Name: "gke",
			Path: fmt.Sprintf("/cloud-accounts/%v/deployment-space/clusters", id),
			Type: "type"},
	}

	return options, nil
}

func (s *Service) FetchCredentials(ctx *gofr.Context, cloudAccountID int64) (interface{}, error) {
	credentials, err := s.store.GetCredentials(ctx, cloudAccountID)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}
