package service

import (
	"github.com/zopdev/zopdev/api/resources/client"
	"github.com/zopdev/zopdev/api/resources/providers/models"
	"github.com/zopdev/zopdev/api/resources/store"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
)

type Service struct {
	gcp   GCPClient
	http  HTTPClient
	store Store
}

func New(gcp GCPClient, http HTTPClient, store Store) *Service {
	return &Service{gcp: gcp, http: http, store: store}
}

func (s *Service) GetAll(ctx *gofr.Context, id int64, resourceType []string) ([]store.Resource, error) {
	res, err := s.store.GetResources(ctx, id, resourceType)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) ChangeState(ctx *gofr.Context, resDetails ResourceDetails) error {
	ca, err := s.http.GetCloudCredentials(ctx, resDetails.CloudAccID)
	if err != nil {
		return err
	}

	switch resDetails.Type {
	case SQL:
		return s.changeSQLState(ctx, ca, resDetails)
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.Type"}}
	}
}

func (s *Service) SyncResources(ctx *gofr.Context, id int64) ([]store.Resource, error) {
	ca, err := s.http.GetCloudCredentials(ctx, id)
	if err != nil {
		return nil, err
	}

	ins, err := s.getAllInstances(ctx, ca)
	if err != nil {
		ctx.Errorf("failed to get all instances: %v", err)
		return nil, err
	}

	res, err := s.store.GetResources(ctx, id, nil)
	if err != nil {
		ctx.Errorf("failed to get existing resources: %v", err)
		return nil, err
	}

	visited := make([]bool, len(res))

	for _, i := range ins {
		idx, found := bSearch(res, i.UID)
		if !found {
			// This is true when the resource is present in the cloud but not in the store.
			err = s.store.InsertResource(ctx, getStoreResource(&i, ca))
			if err != nil {
				ctx.Errorf("failed to insert resource: %v", err)
			}
		} else {
			// else update the existing resource and mark the resource as visited.
			visited[idx] = true
			err = s.store.UpdateResource(ctx, getStoreResource(&i, ca))

			if err != nil {
				ctx.Errorf("failed to update resource: %v", err)
			}
		}
	}

	for i, v := range visited {
		if v {
			continue
		}

		// Remove a resource if it is not visited, i.e., The resource is no longer present on the cloud.
		err = s.store.RemoveResource(ctx, res[i].ID)
		if err != nil {
			ctx.Errorf("failed to remove resource: %v", err)
		}
	}

	return s.GetAll(ctx, id, nil)
}

func bSearch(res []store.Resource, uid string) (int, bool) {
	l, r := 0, len(res)-1

	for l <= r {
		mid := l + (r-l)/2
		if res[mid].UID == uid {
			return mid, true
		}

		if res[mid].UID < uid {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	return -1, false
}

func getStoreResource(res *models.Instance, ca *client.CloudAccount) store.Resource {
	return store.Resource{
		UID:            res.UID,
		Name:           res.Name,
		Type:           store.ResourceType(res.Type),
		State:          res.Status,
		CloudAccountID: ca.ID,
		CloudProvider:  ca.Provider,
	}
}
