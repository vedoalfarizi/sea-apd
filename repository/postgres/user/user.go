package user

import (
	"github.com/jinzhu/gorm"
	"github.com/williamchang80/sea-apd/domain/user"
)

// UserRepository is contract for repo User
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository ...
func NewUserRepository(db *gorm.DB) user.UserRepository {
	if db != nil {
		db.AutoMigrate(&user.User{})
	}
	return &UserRepository{db: db}
}

// CreateUser ...
func (u *UserRepository) CreateUser(user user.User) error {
	if err := u.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// GetUserByEmail ...
func (u *UserRepository) GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}