package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/zopdev/zopdev/api/audit/client"
	"github.com/zopdev/zopdev/api/audit/store"
)

var errMock = errors.New("some internal error")

//nolint:funlen // Test function is long due to multiple test cases
func TestService_RunByID(t *testing.T) {
	ctx, ctrl, mockStore, mockRule, mock := InitlizeTests(t)
	defer ctrl.Finish()

	service := New(mockStore)

	// Mock rule registration
	service.rules["rule-1"] = mockRule
	evalTime := time.Now()
	resData := &store.ResultData{
		Data: []store.Items{
			{InstanceName: "instance-1", Status: "compliant"},
		},
	}
	mockRes := &store.Result{
		ID:             1,
		Result:         resData,
		CloudAccountID: 123,
		RuleID:         "rule-1",
		EvaluatedAt:    evalTime,
	}

	testCases := []struct {
		name           string
		ruleID         string
		cloudAccID     int64
		expectedError  error
		expectedResult *store.Result
		mockCalls      func()
	}{
		{
			name:          "Rule Not Found",
			ruleID:        "non-existent-rule",
			cloudAccID:    123,
			expectedError: gofrHttp.ErrorEntityNotFound{Name: "Rule", Value: "non-existent-rule"},
		},
		{
			name:          "Cloud Credentials Error",
			ruleID:        "rule-1",
			cloudAccID:    123,
			expectedError: errMock,
			mockCalls: func() {
				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(nil, errMock)
			},
		},
		{
			name:          "Error Creating Result Data",
			ruleID:        "rule-1",
			cloudAccID:    123,
			expectedError: errMock,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name:           "Error Updating Result Data",
			ruleID:         "rule-1",
			cloudAccID:     123,
			expectedError:  errMock,
			expectedResult: mockRes,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(resData.Data, nil)
				mockStore.EXPECT().UpdateResult(ctx, gomock.Any()).
					Return(errMock)
			},
		},
		{
			name:          "Error Executing Rule",
			ruleID:        "rule-1",
			cloudAccID:    123,
			expectedError: errMock,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(nil, errMock)
			},
		},
		{
			name:           "Success",
			ruleID:         "rule-1",
			cloudAccID:     123,
			expectedResult: mockRes,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(resData.Data, nil)
				mockStore.EXPECT().UpdateResult(ctx, gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockCalls != nil {
				tc.mockCalls()
			}

			result, err := service.RunByID(ctx, tc.ruleID, tc.cloudAccID)

			assert.Equal(t, tc.expectedError, err, "Error mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expectedResult, result, "Result mismatch for test case: %s", tc.name)
		})
	}
}

//nolint:funlen // Test function is long due to multiple test cases
func TestService_RunByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mock := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))
	mockStore := NewMockStore(ctrl)
	mockRule := NewMockRule(ctrl)
	service := New(mockStore)

	// Mock rule registration
	service.rules["rule-1"] = mockRule
	service.categoryRuleMap["overprovision"] = []Rule{mockRule}

	ctx := &gofr.Context{
		Container: mockContainer,
	}
	evalTime := time.Now()
	resData := &store.ResultData{
		Data: []store.Items{
			{InstanceName: "instance-1", Status: "compliant"},
		},
	}
	mockRes := &store.Result{
		ID:             1,
		Result:         resData,
		CloudAccountID: 123,
		RuleID:         "rule-1",
		EvaluatedAt:    evalTime,
	}

	testCases := []struct {
		name           string
		category       string
		cloudAccID     int64
		expectedError  error
		expectedResult []*store.Result
		mockCalls      func()
	}{
		{
			name:          "error from cloud-account client",
			category:      "overprovion",
			cloudAccID:    123,
			expectedError: errMock,
			mockCalls: func() {
				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(nil, errMock)
			},
		},
		{
			name:          "category not found",
			category:      "non-existent-category",
			cloudAccID:    123,
			expectedError: gofrHttp.ErrorEntityNotFound{Name: "Category", Value: "non-existent-category"},
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
			},
		},
		{
			name:           "error creating result data",
			category:       "overprovision",
			cloudAccID:     123,
			expectedResult: []*store.Result{},
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name:          "error executing rule",
			category:      "overprovision",
			cloudAccID:    123,
			expectedError: nil,
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().GetName().Return("rule-1")
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(nil, errMock)
			},
		},
		{
			name:           "Success",
			category:       "overprovision",
			cloudAccID:     123,
			expectedResult: []*store.Result{mockRes},
			mockCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().GetName().Return("rule-1")
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(resData.Data, nil)
				mockStore.EXPECT().UpdateResult(ctx, gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockCalls != nil {
				tc.mockCalls()
			}

			results, err := service.RunByCategory(ctx, tc.category, tc.cloudAccID)

			assert.Equal(t, tc.expectedError, err, "Error mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expectedResult, results, "Result mismatch for test case: %s", tc.name)
		})
	}
}

//nolint:funlen // Test function is long due to multiple test cases
func TestService_RunAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContainer, mock := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))
	mockStore := NewMockStore(ctrl)
	mockRule := NewMockRule(ctrl)
	service := New(mockStore)

	// Mock rule registration
	service.rules = map[string]Rule{
		"rule-1": mockRule,
	}
	ctx := &gofr.Context{
		Container: mockContainer,
	}
	evalTime := time.Now()
	resData := &store.ResultData{
		Data: []store.Items{
			{InstanceName: "instance-1", Status: "compliant"},
		},
	}
	mockRes := &store.Result{
		ID:             1,
		Result:         resData,
		CloudAccountID: 123,
		RuleID:         "rule-1",
		EvaluatedAt:    evalTime,
	}

	testCases := []struct {
		name           string
		cloudAccID     int64
		expectedError  error
		expectedResult map[string][]*store.Result
		expectedCalls  func()
	}{
		{
			name:          "error from cloud-account client",
			cloudAccID:    123,
			expectedError: errMock,
			expectedCalls: func() {
				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(nil, errMock)
			},
		},
		{
			name:           "error creating result data",
			cloudAccID:     123,
			expectedError:  nil,
			expectedResult: map[string][]*store.Result{},
			expectedCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name:          "error executing rule",
			cloudAccID:    123,
			expectedError: errMock,
			expectedCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(nil, errMock)
			},
		},
		{
			name:       "Success",
			cloudAccID: 123,
			expectedResult: map[string][]*store.Result{
				"overprovision": []*store.Result{mockRes},
			},
			expectedCalls: func() {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"data": {"id": 123, "name": "Test Cloud Account"}}`))),
				}

				mock.HTTPService.EXPECT().Get(ctx, "cloud-accounts/123/credentials", nil).
					Return(resp, nil)
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().CreatePending(ctx, gomock.Any()).
					Return(mockRes, nil)
				mockRule.EXPECT().Execute(ctx, &client.CloudAccount{ID: 123, Name: "Test Cloud Account"}).
					Return(resData.Data, nil)
				mockStore.EXPECT().UpdateResult(ctx, gomock.Any()).
					Return(nil)
				mockRule.EXPECT().GetCategory().Return("overprovision").MaxTimes(4)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedCalls != nil {
				tc.expectedCalls()
			}

			results, err := service.RunAll(ctx, tc.cloudAccID)

			assert.Equal(t, tc.expectedError, err, "Error mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expectedResult, results, "Result mismatch for test case: %s", tc.name)
		})
	}
}

func TestService_GetResultById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockStore(ctrl)
	service := New(mockStore)

	ctx := &gofr.Context{}

	evalTime := time.Now()
	resData := &store.ResultData{
		Data: []store.Items{
			{InstanceName: "instance-1", Status: "compliant"},
		},
	}
	mockRes := &store.Result{
		ID:             1,
		Result:         resData,
		CloudAccountID: 123,
		RuleID:         "rule-1",
		EvaluatedAt:    evalTime,
	}

	testCases := []struct {
		name           string
		ruleID         string
		cloudAccID     int64
		expectedError  error
		expectedResult *store.Result
		mockCalls      func()
	}{
		{
			name:          "error getting result",
			ruleID:        "rule-1",
			cloudAccID:    123,
			expectedError: errMock,
			mockCalls: func() {
				mockStore.EXPECT().GetLastRun(ctx, int64(123), "rule-1").
					Return(nil, errMock)
			},
		},
		{
			name:          "result not found",
			ruleID:        "rule-1",
			cloudAccID:    123,
			expectedError: gofrHttp.ErrorEntityNotFound{Name: "Result", Value: "rule-1"},
			mockCalls: func() {
				mockStore.EXPECT().GetLastRun(ctx, int64(123), "rule-1").
					Return(nil, nil)
			},
		},
		{
			name:           "Success",
			ruleID:         "rule-1",
			cloudAccID:     123,
			expectedResult: mockRes,
			mockCalls: func() {
				mockStore.EXPECT().GetLastRun(ctx, int64(123), "rule-1").
					Return(mockRes, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockCalls != nil {
				tc.mockCalls()
			}

			result, err := service.GetResultByID(ctx, tc.cloudAccID, tc.ruleID)

			assert.Equal(t, tc.expectedError, err, "Error mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expectedResult, result, "Result mismatch for test case: %s", tc.name)
		})
	}
}

func TestService_GetResultsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := NewMockStore(ctrl)
	mockRule := NewMockRule(ctrl)
	service := New(mockStore)
	service.rules = map[string]Rule{
		"rule-1": mockRule,
	}
	ctx := &gofr.Context{}
	evalTime := time.Now()
	resData := &store.ResultData{
		Data: []store.Items{
			{InstanceName: "instance-1", Status: "compliant"},
		},
	}
	mockRes := &store.Result{
		ID:             1,
		Result:         resData,
		CloudAccountID: 123,
		RuleID:         "rule-1",
		EvaluatedAt:    evalTime,
	}

	testCases := []struct {
		name           string
		cloudAccID     int64
		expectedError  error
		expectedResult map[string][]*store.Result
		mockCalls      func()
	}{
		{
			name:          "error getting results",
			cloudAccID:    123,
			expectedError: errMock,
			mockCalls: func() {
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().GetLastRun(ctx, int64(123), "rule-1").
					Return(nil, errMock)
			},
		},
		{
			name:           "results not found",
			cloudAccID:     123,
			expectedResult: map[string][]*store.Result{},
			mockCalls: func() {
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().GetLastRun(ctx, int64(123), "rule-1").
					Return(nil, nil)
			},
		},
		{
			name:           "Success",
			cloudAccID:     123,
			expectedResult: map[string][]*store.Result{"overprovision": []*store.Result{mockRes}},
			mockCalls: func() {
				mockRule.EXPECT().GetName().Return("rule-1")
				mockStore.EXPECT().GetLastRun(ctx, int64(123), "rule-1").
					Return(mockRes, nil)
				mockRule.EXPECT().GetCategory().Return("overprovision").MaxTimes(4)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockCalls != nil {
				tc.mockCalls()
			}

			results, err := service.GetAllResults(ctx, tc.cloudAccID)

			assert.Equal(t, tc.expectedError, err, "Error mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expectedResult, results, "Result mismatch for test case: %s", tc.name)
		})
	}
}

func InitlizeTests(t *testing.T) (*gofr.Context, *gomock.Controller, *MockStore, *MockRule, *container.Mocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockContainer, mock := container.NewMockContainer(t, container.WithMockHTTPService("cloud-account"))
	mockStore := NewMockStore(ctrl)
	mockRule := NewMockRule(ctrl)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	return ctx, ctrl, mockStore, mockRule, mock
}
