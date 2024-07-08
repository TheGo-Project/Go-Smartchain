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

func (r Repository) CreateUser(db *gorm.DB, user models.User) (models.User, error) {
	err := db.Create(&user).Error

	slog.Info("Created user: ", "user", user)

	return user, err
}

func (r Repository) GetUsersCount(db *gorm.DB) (int64, error) {
	var count int64

	err := db.Model(&models.User{}).Count(&count).Error

	return count, err
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
		slog.Error("GetAccountByID", "error", result.Error)
	}

	return account, result.Error
}

func (r Repository) GetAccountsByExtID(db *gorm.DB, extID string) (models.Account, error) {
	var account models.Account
	result := db.Where("ext_id = ?", extID).Find(&account)

	if result.Error != nil {
		slog.Error("GetAccountsByExtID", "error", result.Error)
	}

	return account, result.Error
}

func (r Repository) DeleteAllAccounts(db *gorm.DB) error {
	return db.Where("deleted_at IS NULL").Delete(&models.Account{}).Error
}

func (r Repository) GetParam(db *gorm.DB, key string) (string, error) {
	var param models.Param
	result := db.Where("key = ?", key).Find(&param)

	if result.Error != nil {
		slog.Error("GetParam", "error", result.Error)
	}

	return param.Value, result.Error
}

func (r Repository) ParamExists(db *gorm.DB, key string) (bool, error) {
	var exists bool
	err := db.Model(models.Param{}).
		Select("count(*) > 0").
		Where("key = ?", key).
		Find(&exists).
		Error

	return exists, err
}

func (r Repository) SetParam(db *gorm.DB, key string, value string) (string, error) {
	slog.Info("SetParam", "key", key, "value", value)

	exists, err := r.ParamExists(db, key)
	if err != nil {
		return "", err
	}
	slog.Info("SetParam", "exists", exists)

	if !exists {
		param := models.Param{Key: key, Value: value}

		result := db.Create(&param)
		if result.Error != nil {
			slog.Info("SetParam", "error", result.Error)

			return "", result.Error
		}
	} else {
		err = db.Model(&models.Param{}).Where("key = ?", key).Update("value", value).Error
		if err != nil {
			slog.Error("SetParam", "error", err, "key", key, "value", value)

			return "", err
		}
	}
	return value, nil
}
