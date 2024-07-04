package repository

import (
	"log/slog"

	"github.com/rosty-git/Smartchain-backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
}

func MakeRepository() Repository {
	return Repository{}
}

func (r Repository) Create(db *gorm.DB, user models.User) (models.User, error) {
	err := db.Create(&user).Error

	slog.Info("Created user: ", "user", user)

	return user, err
}

func (r Repository) GetUserByEmail(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)

	return user, result.Error
}

func (r Repository) GetAccountsByUserID(db *gorm.DB, id string) ([]models.Account, error) {
	var accounts []models.Account
	result := db.Where("user_id = ?", id).Find(&accounts)

	if result.Error != nil {
		slog.Error("GetAccountsByUserID: ", "error", result.Error)
	}

	return accounts, result.Error
}

func (r Repository) GetAccountByID(db *gorm.DB, id string) (models.Account, error) {
	var account models.Account
	result := db.Where("id = ?", id).Find(&account)

	if result.Error != nil {
		slog.Error("GetAccountByID: ", "error", result.Error)
	}

	return account, result.Error
}

func (r Repository) GetAccountsByExtID(db *gorm.DB, extID string) (models.Account, error) {
	var account models.Account
	result := db.Where("ext_id = ?", extID).Find(&account)

	if result.Error != nil {
		slog.Error("GetAccountByID: ", "error", result.Error)
	}

	return account, result.Error
}
