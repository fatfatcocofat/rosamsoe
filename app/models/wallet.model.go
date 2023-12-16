package models

import (
	"time"

	"github.com/google/uuid"
)

var WalletCurrencies = []string{"IDR", "USD"}

type Wallet struct {
	ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    *uuid.UUID
	User      User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Address   string     `gorm:"type:varchar(225);uniqueIndex;not null"`
	Balance   float64    `gorm:"type:numeric(10,2);default:0;not null"`
	Currency  string     `gorm:"type:varchar(50);default:'IDR';not null"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"default:null"`
}

type WalletResponse struct {
	ID        *uuid.UUID   `json:"id"`
	User      UserResponse `json:"user"`
	Address   string       `json:"address"`
	Balance   float64      `json:"balance"`
	Currency  string       `json:"currency"`
	CreatedAt *time.Time   `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}

type WalletCreateRequest struct {
	Currency string `json:"currency" validate:"required"`
}

type WalletUpdateRequest struct {
	Currency string `json:"currency"`
}

func WalletFilterRecord(wallet *Wallet) WalletResponse {
	return WalletResponse{
		ID:        wallet.ID,
		User:      UserFilterRecord(&wallet.User),
		Address:   wallet.Address,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
