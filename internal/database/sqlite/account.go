package sqlite

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SchemaAccount struct {
	ID         uint32 `gorm:"primaryKey"`
	Email      string `gorm:"size:128,unique"`
	Password   string `gorm:"size:128"`
	IsVerified bool
	IsBlocked  bool
	CreatedAt  time.Time `gorm:"autoCreateTime:true"`
	UpdatedAt  time.Time `gorm:"autoCreateTime:true,autoUpdateTime:true"`
}

func (m *SqliteDB) IsEmailExists(email string) (isExists bool, _ error) {
	tx := m.db.First(&isExists, "email = ?", email)
	return isExists, tx.Error
}

func (m *SqliteDB) CreateAccount(c context.Context, email, password string) error {
	return gorm.G[SchemaAccount](m.db).Create(c, &SchemaAccount{
		Email:    email,
		Password: password,
	})
}

func (m *SqliteDB) LoadAccount(c context.Context, email, password string) (_ *SchemaAccount, e error) {
	// todo - to sync.Pool ??
	acc := SchemaAccount{}

	acc, e = gorm.G[SchemaAccount](m.db).Where("email = ? and password = ?", email, password).First(c)
	return &acc, e
}
