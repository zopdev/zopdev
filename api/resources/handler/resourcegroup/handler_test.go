package resourcegroup

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestHandler_GetAllResourceGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
}
