package user

import (
	"errors"
	auth2 "github.com/williamchang80/sea-apd/common/auth"
	"github.com/williamchang80/sea-apd/common/constants/role"
	"github.com/williamchang80/sea-apd/domain/user"
	"github.com/williamchang80/sea-apd/dto/request/auth"
)

type UserUsecase struct {
	repo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return UserUsecase{repo: repo}
}

func convertRegisterRequestToUserDomain(request auth.RegisterUserRequest ) user.User {
	userRoles := role.GetUserRoles()
	return user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: auth2.HashPassword(request.Password),
		Role:     userRoles["CUSTOMER"],
	}
}

func (u UserUsecase) CreateUser(request auth.RegisterUserRequest) error {
	if request.Password != request.PasswordConfirmation {
		return errors.New("password and confirmation password must be same")
	}
	user := convertRegisterRequestToUserDomain(request)
	if err := u.repo.CreateUser(user); err != nil {
		return errors.New("email must be unique")
	}
	return nil
}