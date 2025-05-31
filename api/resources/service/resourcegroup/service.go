package resourcegroup

import (
	"strconv"
	"sync"

	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/models"
)

const (
	STOPPED = "STOPPED"
	RUNNING = "RUNNING"
)

type Service struct {
	grpStore RGStore
	resSvc   ResourceService
}

func New(store RGStore, rsSvc ResourceService) *Service {
	return &Service{grpStore: store, resSvc: rsSvc}
}

func (s *Service) GetAllResourceGroups(ctx *gofr.Context, cloudAccID int64) ([]models.ResourceGroupData, error) {
	rsg, err := s.grpStore.GetAllResourceGroups(ctx, cloudAccID)
	if err != nil {
		return nil, &errInternalServer{}
	}

	var (
		resourceGroupData []models.ResourceGroupData
		mu                sync.Mutex
		wg                sync.WaitGroup
	)

	for _, rg := range rsg {
		wg.Add(1)
		go func(rg models.ResourceGroup) {
			defer wg.Done()

			// Get resource IDs for the current resource group
			resIDs, er := s.grpStore.GetResourceIDs(ctx, rg.ID)
			if er != nil {
				return
			}

			var resources []models.Resource

			for i := range resIDs {
				resource, er := s.resSvc.GetByID(ctx, resIDs[i])
				if er != nil {
					return
				}

				resources = append(resources, *resource)
			}

			mu.Lock()
			resourceGroupData = append(resourceGroupData, models.ResourceGroupData{
				ResourceGroup: rg,
				Resources:     resources,
			})
			mu.Unlock()
		}(rg)
	}

	wg.Wait()

	return resourceGroupData, nil
}

func (s *Service) GetResourceGroupByID(ctx *gofr.Context, cloudAccID, id int64) (*models.ResourceGroupData, error) {
	rg, err := s.grpStore.GetResourceGroupByID(ctx, cloudAccID, id)
	if err != nil {
		return nil, &errInternalServer{}
	}

	resIDs, er := s.grpStore.GetResourceIDs(ctx, id)
	if er != nil {
		return nil, &errInternalServer{}
	}

	rg.Status = RUNNING

	var resources []models.Resource

	for i := range resIDs {
		resource, er := s.resSvc.GetByID(ctx, resIDs[i])
		if er != nil {
			return nil, &errInternalServer{}
		}

		if resource.Status == STOPPED {
			rg.Status = STOPPED
		}

		resources = append(resources, *resource)
	}

	return &models.ResourceGroupData{
		ResourceGroup: *rg,
		Resources:     resources,
	}, nil
}

func (s *Service) CreateResourceGroup(ctx *gofr.Context, rg *models.ResourceGroup) (*models.ResourceGroupData, error) {
	id, err := s.grpStore.CreateResourceGroup(ctx, rg)
	if err != nil {
		return nil, &errInternalServer{}
	}

	return s.GetResourceGroupByID(ctx, rg.CloudAccountID, id)
}

func (s *Service) UpdateResourceGroup(ctx *gofr.Context, rg *models.ResourceGroup) (*models.ResourceGroupData, error) {
	if err := s.grpStore.UpdateResourceGroup(ctx, rg); err != nil {
		return nil, &errInternalServer{}
	}

	return s.GetResourceGroupByID(ctx, rg.CloudAccountID, rg.ID)
}

func (s *Service) AddResourceToGroup(ctx *gofr.Context, rgID, resID int64) error {
	res, err := s.resSvc.GetByID(ctx, resID)
	if err != nil {
		return &errInternalServer{}
	}

	if res == nil {
		return gofrHttp.ErrorEntityNotFound{Name: "resource", Value: strconv.FormatInt(resID, 10)}
	}

	grp, err := s.grpStore.GetResourceGroupByID(ctx, res.CloudAccount.ID, rgID)
	if err != nil {
		return &errInternalServer{}
	}

	if grp == nil {
		return gofrHttp.ErrorEntityNotFound{Name: "resource", Value: strconv.FormatInt(rgID, 10)}
	}

	if er := s.grpStore.AddResourceToGroup(ctx, rgID, resID); er != nil {
		return &errInternalServer{}
	}

	return nil
}

func (s *Service) RemoveResourceFromGroup(ctx *gofr.Context, rgID, resID int64) error {
	res, err := s.resSvc.GetByID(ctx, resID)
	if err != nil {
		return &errInternalServer{}
	}

	if res == nil {
		return gofrHttp.ErrorEntityNotFound{Name: "resource", Value: strconv.FormatInt(resID, 10)}
	}

	grp, err := s.grpStore.GetResourceGroupByID(ctx, res.CloudAccount.ID, rgID)
	if err != nil {
		return &errInternalServer{}
	}

	if grp == nil {
		return gofrHttp.ErrorEntityNotFound{Name: "resource-group", Value: strconv.FormatInt(rgID, 10)}
	}

	if er := s.grpStore.RemoveResourceFromGroup(ctx, rgID, resID); er != nil {
		return &errInternalServer{}
	}

	return nil
}

func (s *Service) DeleteResourceGroup(ctx *gofr.Context, cloudAccID, id int64) error {
	// Check if the resource group exists
	grp, err := s.grpStore.GetResourceGroupByID(ctx, cloudAccID, id)
	if err != nil {
		return &errInternalServer{}
	}

	if grp == nil {
		return gofrHttp.ErrorEntityNotFound{Name: "resource group", Value: strconv.FormatInt(id, 10)}
	}

	if err := s.grpStore.DeleteResourceGroup(ctx, id); err != nil {
		return &errInternalServer{}
	}

	return nil
}
