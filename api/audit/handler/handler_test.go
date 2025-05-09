package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/audit/store"
)

func TestHandler_RunAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name          string
		pathParam     string
		expectedError error
		mockResponse  map[string][]*store.Result
		mockError     error
	}{
		{
			name:          "Missing ID",
			pathParam:     "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:          "Invalid ID",
			pathParam:     "abc",
			expectedError: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:         "Success",
			pathParam:    "123",
			mockResponse: map[string][]*store.Result{},
			mockError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/audit/all", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tc.pathParam})

			ctx := &gofr.Context{
				Request: gofrHttp.NewRequest(r),
			}

			if tc.mockResponse != nil || tc.mockError != nil {
				mockService.EXPECT().RunAll(ctx, int64(123)).
					Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.RunAll(ctx)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.mockResponse, resp)
			}
		})
	}
}

func TestHandler_RunByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name          string
		ruleID        string
		cloudAccID    string
		expectedError error
		mockResponse  *store.Result
		mockError     error
	}{
		{
			name:          "Missing ID",
			cloudAccID:    "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:          "Invalid ID",
			cloudAccID:    "abc",
			expectedError: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:          "Missing Rule ID",
			cloudAccID:    "123",
			ruleID:        "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"ruleId"}},
		},
		{
			name:         "Success",
			cloudAccID:   "123",
			ruleID:       "rule-1",
			mockResponse: &store.Result{},
			mockError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/audit/cloud-accounts/{id}/rule/{ruleID}", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tc.cloudAccID, "ruleId": tc.ruleID})

			ctx := &gofr.Context{
				Request: gofrHttp.NewRequest(r),
			}

			if tc.mockResponse != nil || tc.mockError != nil {
				cID, _ := strconv.ParseInt(tc.cloudAccID, 10, 64)
				mockService.EXPECT().RunByID(ctx, tc.ruleID, cID).
					Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.RunByID(ctx)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.mockResponse, resp)
			}
		})
	}
}

func TestHandler_RunByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name          string
		categoryID    string
		cloudAccID    string
		expectedError error
		mockResponse  []*store.Result
		mockError     error
	}{
		{
			name:          "Missing ID",
			cloudAccID:    "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:          "Invalid ID",
			cloudAccID:    "overprovision",
			expectedError: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:          "Missing Category",
			cloudAccID:    "123",
			categoryID:    "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"category"}},
		},
		{
			name:         "Success",
			cloudAccID:   "123",
			categoryID:   "overprovision",
			mockResponse: []*store.Result{},
			mockError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/audit/cloud-account/{id}/category/{category}", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tc.cloudAccID, "category": tc.categoryID})

			ctx := &gofr.Context{
				Request: gofrHttp.NewRequest(r),
			}

			if tc.mockResponse != nil || tc.mockError != nil {
				cID, _ := strconv.ParseInt(tc.cloudAccID, 10, 64)
				mockService.EXPECT().RunByCategory(ctx, tc.categoryID, cID).
					Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.RunByCategory(ctx)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.mockResponse, resp)
			}
		})
	}
}

func TestHandler_GetResultByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name          string
		ruleID        string
		cloudAccID    string
		expectedError error
		mockResponse  *store.Result
		mockError     error
	}{
		{
			name:          "Missing ID",
			cloudAccID:    "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:          "Invalid ID",
			cloudAccID:    "abc",
			expectedError: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:          "Missing Rule ID",
			cloudAccID:    "123",
			ruleID:        "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"ruleId"}},
		},
		{
			name:         "Success",
			cloudAccID:   "123",
			ruleID:       "rule-1",
			mockResponse: &store.Result{},
			mockError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/audit/cloud-accounts/{id}/rule/{ruleId}", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tc.cloudAccID, "ruleId": tc.ruleID})

			ctx := &gofr.Context{
				Request: gofrHttp.NewRequest(r),
			}

			if tc.mockResponse != nil || tc.mockError != nil {
				cID, _ := strconv.ParseInt(tc.cloudAccID, 10, 64)
				mockService.EXPECT().GetResultByID(ctx, cID, tc.ruleID).
					Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.GetResultByID(ctx)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.mockResponse, resp)
			}
		})
	}
}

func TestHandler_GetResultByCategory(t *testing.T) {
	var storeResult = &store.Result{
		ID:             1,
		RuleID:         "rule-1",
		CloudAccountID: 123,

		Result: &store.ResultData{
			Data: []store.Items{
				{
					InstanceName: "instance-1",
					Status:       "danger",
				},
				{
					InstanceName: "instance-1",
					Status:       "compliant",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := New(mockService)

	testCases := []struct {
		name          string
		cloudAccID    string
		Category      string
		expectedError error
		mockResponse  map[string][]*store.Result
		mockError     error
	}{
		{
			name:          "Missing ID",
			cloudAccID:    "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:          "Invalid ID",
			cloudAccID:    "abc",
			expectedError: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:          "Missing Category",
			cloudAccID:    "123",
			Category:      "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"category"}},
		},
		{
			name:       "Success",
			cloudAccID: "123",
			Category:   "overprovision",
			mockResponse: map[string][]*store.Result{
				"overprovision": {
					storeResult,
				},
			},
			mockError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/audit/cloud-accounts/{id}/category", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tc.cloudAccID, "category": tc.Category})

			ctx := &gofr.Context{
				Request: gofrHttp.NewRequest(r),
			}

			if tc.mockResponse != nil || tc.mockError != nil {
				cID, _ := strconv.ParseInt(tc.cloudAccID, 10, 64)
				mockService.EXPECT().GetResultByCategory(ctx, cID).
					Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.GetResultByCategory(ctx)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.mockResponse, resp)
			}
		})
	}
}

func TestHandler_GetAllResults(t *testing.T) {
	var storeResult = &store.Result{
		ID:             1,
		RuleID:         "rule-1",
		CloudAccountID: 123,

		Result: &store.ResultData{
			Data: []store.Items{
				{
					InstanceName: "instance-1",
					Status:       "danger",
				},
				{
					InstanceName: "instance-1",
					Status:       "compliant",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := New(mockService)
	testCases := []struct {
		name          string
		cloudAccID    string
		expectedError error
		mockResponse  map[string][]*store.Result
		mockError     error
	}{
		{
			name:          "Missing ID",
			cloudAccID:    "",
			expectedError: gofrHttp.ErrorMissingParam{Params: []string{"id"}},
		},
		{
			name:          "Invalid ID",
			cloudAccID:    "abc",
			expectedError: gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:       "Success",
			cloudAccID: "123",
			mockResponse: map[string][]*store.Result{
				"overprovision": []*store.Result{storeResult},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/audit/cloud-accounts/{id}/all", nil)
			r = mux.SetURLVars(r, map[string]string{"id": tc.cloudAccID})

			ctx := &gofr.Context{
				Request: gofrHttp.NewRequest(r),
			}

			if tc.mockResponse != nil || tc.mockError != nil {
				cID, _ := strconv.ParseInt(tc.cloudAccID, 10, 64)
				mockService.EXPECT().GetAllResults(ctx, cID).
					Return(tc.mockResponse, tc.mockError)
			}

			resp, err := handler.GetAllResults(ctx)

			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				assert.Equal(t, tc.mockResponse, resp)
			}
		})
	}
}
