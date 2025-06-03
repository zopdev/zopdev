package resourcegroup

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestService_GetAllResourceGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockRGStore(ctrl)
	mockResSvc := NewMockResourceService(ctrl)
	svc := New(mockStore, mockResSvc)
	ctx := &gofr.Context{}

	rg1 := models.ResourceGroup{ID: 1, Status: STOPPED}
	rg2 := models.ResourceGroup{ID: 2, Status: RUNNING}
	r1 := &models.Resource{ID: 10, Status: RUNNING}
	r2 := &models.Resource{ID: 11, Status: STOPPED}

	tests := []struct {
		name        string
		setup       func()
		expected    []models.ResourceGroupData
		expectedErr error
	}{
		{
			name: "success",
			setup: func() {
				mockStore.EXPECT().GetAllResourceGroups(ctx, int64(1)).Return([]models.ResourceGroup{rg1, rg2}, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg1.ID).Return([]int64{10, 11}, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg2.ID).Return([]int64{10}, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(r1, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(r1, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(11)).Return(r2, nil)
			},
			expected: []models.ResourceGroupData{
				{ResourceGroup: rg2, Resources: []models.Resource{*r1}},
				{ResourceGroup: rg1, Resources: []models.Resource{*r1, *r2}},
			},
			expectedErr: nil,
		},
		{
			name: "store error - get all resource groups",
			setup: func() {
				mockStore.EXPECT().GetAllResourceGroups(ctx, int64(1)).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - get resource IDs",
			setup: func() {
				mockStore.EXPECT().GetAllResourceGroups(ctx, int64(1)).Return([]models.ResourceGroup{rg1}, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg1.ID).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "resource service error - get resource IDs - partial success",
			setup: func() {
				mockStore.EXPECT().GetAllResourceGroups(ctx, int64(1)).Return([]models.ResourceGroup{rg1, rg2}, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg1.ID).Return(nil, assert.AnError)
				mockStore.EXPECT().GetResourceIDs(ctx, rg2.ID).Return([]int64{10}, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(r1, nil)
			},
			expectedErr: &errInternalServer{},
			expected: []models.ResourceGroupData{
				{ResourceGroup: rg2, Resources: []models.Resource{*r1}},
			},
		},
		{
			name: "resource service error - get resource by ID",
			setup: func() {
				mockStore.EXPECT().GetAllResourceGroups(ctx, int64(1)).Return([]models.ResourceGroup{rg1}, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg1.ID).Return([]int64{10}, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			got, err := svc.GetAllResourceGroups(ctx, 1)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestService_GetResourceGroupByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockRGStore(ctrl)
	mockResSvc := NewMockResourceService(ctrl)
	svc := New(mockStore, mockResSvc)
	ctx := &gofr.Context{}

	rg := &models.ResourceGroup{ID: 1, Status: STOPPED}
	r1 := &models.Resource{ID: 10, Status: RUNNING}
	r2 := &models.Resource{ID: 11, Status: STOPPED}
	resourceIDs := []int64{10, 11}

	tests := []struct {
		name        string
		setup       func()
		expected    *models.ResourceGroupData
		expectedErr error
	}{
		{
			name: "success",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, int64(1)).Return(resourceIDs, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(r1, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(11)).Return(r2, nil)
			},
			expected: &models.ResourceGroupData{
				ResourceGroup: *rg,
				Resources:     []models.Resource{*r1, *r2},
			},
			expectedErr: nil,
		},
		{
			name: "store error - group not found",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).
					Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - resources not found",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, int64(1)).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "resource service error",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, int64(1)).Return(resourceIDs, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			got, err := svc.GetResourceGroupByID(ctx, 1, 1)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestService_CreateResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockRGStore(ctrl)
	mockResSvc := NewMockResourceService(ctrl)
	svc := New(mockStore, mockResSvc)
	ctx := &gofr.Context{}

	rgCreate := &models.RGCreate{CloudAccountID: 1, ResourceIDs: []int64{10}}
	rg := &models.ResourceGroup{ID: 1}
	resource := &models.Resource{ID: 10}
	resourceIDs := []int64{10}

	tests := []struct {
		name        string
		setup       func()
		expectedErr error
	}{
		{
			name: "success",
			setup: func() {
				mockStore.EXPECT().CreateResourceGroup(ctx, rgCreate).Return(int64(1), nil)
				mockStore.EXPECT().AddResourcesToGroup(ctx, int64(1), rgCreate.ResourceIDs).Return(nil)
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, int64(1)).Return(resourceIDs, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(resource, nil)
			},
			expectedErr: nil,
		},
		{
			name: "store error",
			setup: func() {
				mockStore.EXPECT().CreateResourceGroup(ctx, rgCreate).
					Return(int64(0), assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - add resources to group",
			setup: func() {
				mockStore.EXPECT().CreateResourceGroup(ctx, rgCreate).Return(int64(1), nil)
				mockStore.EXPECT().AddResourcesToGroup(ctx, int64(1), rgCreate.ResourceIDs).Return(assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			_, err := svc.CreateResourceGroup(ctx, rgCreate)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestService_UpdateResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockRGStore(ctrl)
	mockResSvc := NewMockResourceService(ctrl)
	svc := New(mockStore, mockResSvc)
	ctx := &gofr.Context{}

	rgUpdate := &models.RGUpdate{ID: 2, CloudAccountID: 1, ResourceIDs: []int64{10, 11}}
	rg := &models.ResourceGroup{ID: 2, CloudAccountID: 1, Name: "myGroup", Status: RUNNING}
	r1 := &models.Resource{ID: 10}
	r2 := &models.Resource{ID: 11}
	resourceIDs := []int64{10}

	tests := []struct {
		name        string
		setup       func()
		expectedErr error
		expectedOut *models.ResourceGroupData
	}{
		{
			name: "success",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rgUpdate.ID).Return(resourceIDs, nil) // existing resource IDs - 10
				mockStore.EXPECT().UpdateResourceGroup(ctx, rgUpdate).Return(nil)
				mockStore.EXPECT().AddResourcesToGroup(ctx, rgUpdate.ID, []int64{11}).Return(nil) // adding new resource ID - 11
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg.ID).Return([]int64{11, 10}, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(11)).Return(r1, nil)
				mockResSvc.EXPECT().GetByID(ctx, int64(10)).Return(r2, nil)
			},
			expectedOut: &models.ResourceGroupData{
				ResourceGroup: *rg,
				Resources:     []models.Resource{*r1, *r2},
			},
			expectedErr: nil,
		},
		{
			name: "error - adding new resource",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rgUpdate.ID).Return(resourceIDs, nil) // existing resource IDs - 10
				mockStore.EXPECT().UpdateResourceGroup(ctx, rgUpdate).Return(nil)
				mockStore.EXPECT().AddResourcesToGroup(ctx, rgUpdate.ID, []int64{11}).Return(assert.AnError) // adding new resource ID - 11
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "not found",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).Return(nil, nil)
			},
			expectedErr: gofrHttp.ErrorEntityNotFound{Name: "resource group", Value: strconv.FormatInt(rg.ID, 10)},
		},
		{
			name: "store error - get resource group",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).
					Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - get resource IDs",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rgUpdate.ID).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - update resource group",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, rg.CloudAccountID, rgUpdate.ID).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rgUpdate.ID).Return(resourceIDs, nil)
				mockStore.EXPECT().UpdateResourceGroup(ctx, rgUpdate).Return(assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			res, err := svc.UpdateResourceGroup(ctx, rgUpdate)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedOut, res)
		})
	}
}

func TestService_modifyResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockRGStore(ctrl)
	svc := &Service{grpStore: mockStore}
	ctx := &gofr.Context{}

	rgID := int64(1)
	existingResources := []int64{10, 11}
	resourceIDs := []int64{11, 12}

	tests := []struct {
		name        string
		setup       func()
		expectedErr error
	}{
		{
			name: "success - modify resources",
			setup: func() {
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, rgID, int64(10)).Return(nil)
				mockStore.EXPECT().AddResourcesToGroup(ctx, rgID, []int64{12}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "error - remove resource from group",
			setup: func() {
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, rgID, int64(10)).Return(assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			err := svc.modifyResources(ctx, rgID, existingResources, resourceIDs)

			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestService_DeleteResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockRGStore(ctrl)
	mockResSvc := NewMockResourceService(ctrl)
	svc := New(mockStore, mockResSvc)
	ctx := &gofr.Context{}

	rg := &models.ResourceGroup{ID: 1}

	tests := []struct {
		name        string
		setup       func()
		expectedErr error
	}{
		{
			name: "success",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg.ID).Return([]int64{1, 2}, nil)
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, int64(1), int64(1)).Return(nil)
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, int64(1), int64(2)).Return(nil)
				mockStore.EXPECT().DeleteResourceGroup(ctx, int64(1)).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "store error - getting the group",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).
					Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - getting resource IDs",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).
					Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, int64(1)).Return(nil, assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - removing resource from group",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg.ID).Return([]int64{1, 2}, nil)
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, int64(1), int64(1)).Return(assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "store error - delete resource group",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(rg, nil)
				mockStore.EXPECT().GetResourceIDs(ctx, rg.ID).Return([]int64{1, 2}, nil)
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, int64(1), int64(1)).Return(nil)
				mockStore.EXPECT().RemoveResourceFromGroup(ctx, int64(1), int64(2)).Return(nil)
				mockStore.EXPECT().DeleteResourceGroup(ctx, int64(1)).Return(assert.AnError)
			},
			expectedErr: &errInternalServer{},
		},
		{
			name: "not found",
			setup: func() {
				mockStore.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).Return(nil, nil)
			},
			expectedErr: gofrHttp.ErrorEntityNotFound{Name: "resource group", Value: strconv.FormatInt(1, 10)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			err := svc.DeleteResourceGroup(ctx, 1, 1)

			if err != nil {
				if tc.expectedErr == nil {
					t.Errorf("expected no error, got %v", err)
				}

				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
		})
	}
}
