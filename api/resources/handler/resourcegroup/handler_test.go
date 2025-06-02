package resourcegroup

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/resources/models"
)

func TestHandler_GetAllResourceGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mSvc := NewMockService(ctrl)
	h := New(mSvc)
	ctx := &gofr.Context{Context: context.Background()}
	sampleRes := []models.ResourceGroupData{
		{
			ResourceGroup: models.ResourceGroup{
				ID:             1,
				Name:           "Test Group",
				CloudAccountID: 1,
				Status:         "RUNNING",
			},
			Resources: []models.Resource{
				{ID: 1, Name: "MySQL Database", Type: "SQL",
					CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
					Status:       "RUNNING",
				},
				{ID: 2, Name: "PSQL", Type: "SQL",
					CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
					Status:       "RUNNING"},
			},
		},
	}

	testCases := []struct {
		name       string
		cloudAccID string
		expErr     error
		expRes     any
		mockCalls  []*gomock.Call
	}{
		{
			name:       "success",
			cloudAccID: "1",
			expErr:     nil,
			expRes:     sampleRes,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().GetAllResourceGroups(ctx, int64(1)).
					Return(sampleRes, nil),
			},
		},
		{
			name:       "missing cloud account ID",
			cloudAccID: "",
			expErr:     gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCalls:  nil,
		},
		{
			name:       "invalid cloud account ID",
			cloudAccID: "invalid",
			expErr:     gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCalls:  nil,
		},
		{
			name:       "service error",
			cloudAccID: "1",
			expErr:     assert.AnError,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().GetAllResourceGroups(ctx, int64(1)).
					Return(nil, assert.AnError),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/resources", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tc.cloudAccID})
			ctx.Request = gofrHttp.NewRequest(req)
			res, err := h.GetAllResourceGroups(ctx)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestHandler_GetResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mSvc := NewMockService(ctrl)
	h := New(mSvc)
	ctx := &gofr.Context{Context: context.Background()}
	sampleRes := &models.ResourceGroupData{
		ResourceGroup: models.ResourceGroup{
			ID:             1,
			Name:           "Test Group",
			CloudAccountID: 1,
			Status:         "RUNNING",
		},
		Resources: []models.Resource{
			{ID: 1, Name: "MySQL Database", Type: "SQL",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Status:       "RUNNING",
			},
			{ID: 2, Name: "PSQL", Type: "SQL",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Status:       "RUNNING"},
		},
	}

	testCases := []struct {
		name      string
		accID     string
		rgID      string
		expErr    error
		expRes    any
		mockCalls []*gomock.Call
	}{
		{
			name:   "success",
			accID:  "1",
			rgID:   "1",
			expErr: nil,
			expRes: sampleRes,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).
					Return(sampleRes, nil),
			},
		},
		{
			name:      "missing cloud account ID",
			accID:     "",
			rgID:      "1",
			expErr:    gofrHttp.ErrorMissingParam{Params: []string{"id"}},
			mockCalls: nil,
		},
		{
			name:      "invalid cloud account ID",
			accID:     "invalid",
			rgID:      "1",
			expErr:    gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
			mockCalls: nil,
		},
		{
			name:      "missing resource group ID",
			accID:     "1",
			rgID:      "",
			expErr:    gofrHttp.ErrorMissingParam{Params: []string{"rgId"}},
			mockCalls: nil,
		},
		{
			name:      "invalid resource group ID",
			accID:     "1",
			rgID:      "invalid",
			expErr:    gofrHttp.ErrorInvalidParam{Params: []string{"rgId"}},
			mockCalls: nil,
		},
		{
			name:   "service error",
			accID:  "1",
			rgID:   "1",
			expErr: assert.AnError,
			expRes: nil,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().GetResourceGroupByID(ctx, int64(1), int64(1)).
					Return(nil, assert.AnError),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cloud-account/{id}/resourcegroup/{rgID}", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tc.accID, "rgID": tc.rgID})
			ctx.Request = gofrHttp.NewRequest(req)
			res, err := h.GetResourceGroup(ctx)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestHandler_CreateResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mSvc := NewMockService(ctrl)
	h := New(mSvc)
	ctx := &gofr.Context{Context: context.Background()}
	createReq := &models.RGCreate{
		Name:           "myGroup",
		Description:    "sample group",
		CloudAccountID: 1,
		ResourceIDs:    []int64{1, 2},
	}
	sampleRes := &models.ResourceGroupData{
		ResourceGroup: models.ResourceGroup{
			ID:             1,
			Name:           "myGroup",
			Description:    "sample group",
			CloudAccountID: 1,
			Status:         "RUNNING",
		},
		Resources: []models.Resource{
			{ID: 1, Name: "MySQL Database", Type: "SQL",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Status:       "RUNNING",
			},
			{ID: 2, Name: "PSQL", Type: "SQL",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Status:       "RUNNING"},
		},
	}

	testCases := []struct {
		name      string
		accID     string
		body      string
		expErr    error
		expRes    any
		mockCalls []*gomock.Call
	}{
		{
			name:   "success",
			accID:  "1",
			body:   `{"name":"myGroup", "description" : "sample group", "resource_ids":[1, 2]}`,
			expRes: sampleRes,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().CreateResourceGroup(ctx, createReq).
					Return(sampleRes, nil),
			},
		},
		{
			name:   "missing cloud account ID",
			accID:  "",
			expErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:   "invalid cloud account ID",
			accID:  "invalid",
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:   "invalid bind",
			accID:  "1",
			body:   `{`,
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"body"}},
		},
		{
			name:   "service error",
			accID:  "1",
			body:   `{"name":"myGroup", "description" : "sample group", "resource_ids":[1, 2]}`,
			expErr: assert.AnError,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().CreateResourceGroup(ctx, createReq).
					Return(nil, assert.AnError),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost,
				"/cloud-account/{id}/resourcegroup", bytes.NewBufferString(tc.body))
			req = mux.SetURLVars(req, map[string]string{"id": tc.accID})

			req.Header.Set("Content-Type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			res, err := h.CreateResourceGroup(ctx)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestHandler_UpdateResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mSvc := NewMockService(ctrl)
	h := New(mSvc)
	ctx := &gofr.Context{Context: context.Background()}
	updateReq := &models.RGUpdate{
		ID:             1,
		Name:           "myGroup",
		Description:    "sample group",
		CloudAccountID: 1,
		ResourceIDs:    []int64{2, 4},
	}
	sampleRes := &models.ResourceGroupData{
		ResourceGroup: models.ResourceGroup{
			ID:             1,
			Name:           "myGroup",
			Description:    "sample group",
			CloudAccountID: 1,
			Status:         "RUNNING",
		},
		Resources: []models.Resource{
			{ID: 2, Name: "PSQL", Type: "SQL",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Status:       "RUNNING"},
			{ID: 4, Name: "MySQL Database", Type: "SQL",
				CloudAccount: models.CloudAccount{ID: 1, Type: "GCP"},
				Status:       "RUNNING"},
		},
	}

	testCases := []struct {
		name      string
		accID     string
		groupID   string
		body      string
		expErr    error
		expRes    any
		mockCalls []*gomock.Call
	}{
		{
			name:    "success",
			accID:   "1",
			groupID: "1",
			body:    `{"name":"myGroup", "description" : "sample group", "resource_ids":[2, 4]}`,
			expRes:  sampleRes,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().UpdateResourceGroup(ctx, updateReq).
					Return(sampleRes, nil),
			},
		},
		{
			name:   "missing cloud account ID",
			accID:  "",
			expErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:   "invalid cloud account ID",
			accID:  "invalid",
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:    "missing resource group ID",
			accID:   "1",
			groupID: "",
			expErr:  gofrHttp.ErrorMissingParam{Params: []string{"rgId"}},
		},
		{
			name:    "invalid resource group ID",
			accID:   "1",
			groupID: "invalid",
			expErr:  gofrHttp.ErrorInvalidParam{Params: []string{"rgId"}},
		},
		{
			name:    "invalid bind",
			accID:   "1",
			groupID: "1",
			body:    `{`,
			expErr:  gofrHttp.ErrorInvalidParam{Params: []string{"body"}},
		},
		{
			name:    "service error",
			accID:   "1",
			groupID: "1",
			body:    `{"name":"myGroup", "description" : "sample group", "resource_ids":[2, 4]}`,
			expErr:  assert.AnError,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().UpdateResourceGroup(ctx, updateReq).
					Return(nil, assert.AnError),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost,
				"/cloud-account/{id}/resourcegroup", bytes.NewBufferString(tc.body))
			req = mux.SetURLVars(req, map[string]string{"id": tc.accID, "rgID": tc.groupID})

			req.Header.Set("Content-Type", "application/json")

			ctx.Request = gofrHttp.NewRequest(req)

			res, err := h.UpdateResourceGroup(ctx)

			assert.Equal(t, tc.expErr, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestHandler_DeleteResourceGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mSvc := NewMockService(ctrl)
	h := New(mSvc)
	ctx := &gofr.Context{Context: context.Background()}

	testCases := []struct {
		name      string
		accID     string
		groupID   string
		expErr    error
		mockCalls []*gomock.Call
	}{
		{
			name:    "success",
			accID:   "1",
			groupID: "1",
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().DeleteResourceGroup(ctx, int64(1), int64(1)).
					Return(nil),
			},
		},
		{
			name:   "missing cloud account ID",
			accID:  "",
			expErr: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:   "invalid cloud account ID",
			accID:  "invalid",
			expErr: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:    "missing resource group ID",
			accID:   "1",
			groupID: "",
			expErr:  gofrHttp.ErrorMissingParam{Params: []string{"rgId"}},
		},
		{
			name:    "invalid resource group ID",
			accID:   "1",
			groupID: "invalid",
			expErr:  gofrHttp.ErrorInvalidParam{Params: []string{"rgId"}},
		},
		{
			name:    "service error",
			accID:   "1",
			groupID: "1",
			expErr:  assert.AnError,
			mockCalls: []*gomock.Call{
				mSvc.EXPECT().DeleteResourceGroup(ctx, int64(1), int64(1)).
					Return(assert.AnError),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost,
				"/cloud-account/{id}/resourcegroup", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"id": tc.accID, "rgID": tc.groupID})

			ctx.Request = gofrHttp.NewRequest(req)

			res, err := h.DeleteResourceGroup(ctx)

			assert.Equal(t, tc.expErr, err)
			assert.Nil(t, res)
		})
	}

}
