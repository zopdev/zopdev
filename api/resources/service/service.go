package service

import (
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/models"
)

type Service struct {
	gcp   GCPClient
	http  HTTPClient
	store Store
}

func New(gcp GCPClient, http HTTPClient, store Store) *Service {
	return &Service{gcp: gcp, http: http, store: store}
}

func (s *Service) GetAll(ctx *gofr.Context, id int64, resourceType []string) ([]models.Instance, error) {
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
		err = s.changeSQLState(ctx, ca, resDetails)
		if err != nil {
			ctx.Errorf("failed to change SQL state: %v", err)
		}

		err = s.store.UpdateResource(ctx, &models.Instance{ID: resDetails.ID, Status: getStatus(resDetails.State)})
		if err != nil {
			ctx.Errorf("failed to update resource status: %v", err)
		}

		return nil
	default:
		return gofrHttp.ErrorInvalidParam{Params: []string{"req.Type"}}
	}
}

func getStatus(action ResourceState) string {
	switch action {
	case START:
		return RUNNING
	case SUSPEND:
		return STOPPED
	default:
		return ""
	}
}

func (s *Service) SyncResources(ctx *gofr.Context, id int64) ([]models.Instance, error) {
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

	for i := range ins {
		idx, found := bSearch(res, ins[i].UID)
		if !found {
			// This is true when the resource is present in the cloud but not in the store.
			err = s.store.InsertResource(ctx, &ins[i])
			if err != nil {
				ctx.Errorf("failed to insert resource: %v", err)
			}
		} else {
			// else update the existing resource and mark the resource as visited.
			visited[idx] = true
			err = s.store.UpdateResource(ctx, &ins[i])

			if err != nil {
				ctx.Errorf("failed to update resource: %v", err)
			}
		}
	}

	s.removeStale(ctx, visited, res)

	return s.GetAll(ctx, id, nil)
}

func (s *Service) removeStale(ctx *gofr.Context, visited []bool, res []models.Instance) {
	for i, v := range visited {
		if v {
			continue
		}

		// Remove a resource if it is not visited, i.e., The resource is no longer present on the cloud.
		err := s.store.RemoveResource(ctx, res[i].ID)
		if err != nil {
			ctx.Errorf("failed to remove resource: %v", err)
		}
	}
}

// bSearch performs a binary search on the sorted slice of models.Instance.
func bSearch(res []models.Instance, uid string) (int, bool) {
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
