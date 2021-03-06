package admin

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/williamchang80/sea-apd/domain/user"
	"github.com/williamchang80/sea-apd/dto/request/admin"
)

var emptyAdmin = user.User{}
var emptyAdminRequest = admin.Admin{}

// MockUsecase ...
type MockUsecase struct {
	ctrl *gomock.Controller
}

// RegisterAdmin ...
func (m MockUsecase) RegisterAdmin(req admin.Admin) error {
	if req == emptyAdminRequest {
		return errors.New("Cannot register admin")
	}
	return nil
}

// NewMockUsecase ...
func NewMockUsecase(repo *gomock.Controller) *MockUsecase {
	return &MockUsecase{
		ctrl: repo,
	}
}
