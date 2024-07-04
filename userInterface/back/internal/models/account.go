package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeyStore struct {
	Address string `json:"address"`
	Crypto  struct {
		Cipher       string `json:"cipher"`
		Ciphertext   string `json:"ciphertext"`
		Cipherparams struct {
			Iv string `json:"iv"`
		} `json:"cipherparams"`
		Kdf       string `json:"kdf"`
		Kdfparams struct {
			Dklen int    `json:"dklen"`
			N     int    `json:"n"`
			P     int    `json:"p"`
			R     int    `json:"r"`
			Salt  string `json:"salt"`
		} `json:"kdfparams"`
		Mac string `json:"mac"`
	} `json:"crypto"`
	ID      string `json:"id"`
	Version int    `json:"version"`
}

func (ks KeyStore) Value() (driver.Value, error) {
	return json.Marshal(ks)
}

func (ks *KeyStore) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, ks)
}

type Account struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	UserID   string   `json:"user_id"`
	ExtID    string   `json:"ext_id"`
	Keystore KeyStore `json:"-" gorm:"type:jsonb"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return
}
