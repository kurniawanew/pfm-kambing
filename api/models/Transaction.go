package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                uint32    `gorm:"primary_key;auto_increment" json:"id"`
	TransactionName   string    `gorm:"size:255;not null" json:"transaction_name"`
	TransactionAmount int       `gorm:"not null" json:"transaction_amount"`
	TransactionDate   string    `gorm:"size:255;not null" json:"transaction_date"`
	TransactionUser   User      `json:"transaction_user"`
	TransactionUserID uint32    `sql:"type:int REFERENCES users(id)" json:"transaction_user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (t *Transaction) Prepare() {

	t.ID = 0
	t.TransactionName = html.EscapeString(strings.TrimSpace(t.TransactionName))
	t.TransactionDate = html.EscapeString(strings.TrimSpace(t.TransactionDate))
	t.TransactionAmount = int(t.TransactionAmount)
	t.TransactionUser = User{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Transaction) Validate() error {

	if t.TransactionName == "" {
		return errors.New("Required Transaction Name")
	}
	if t.TransactionAmount < 0 {
		return errors.New("Required Transaction Amount")
	}
	if t.TransactionDate == "" {
		return errors.New("Required Transaction Date")
	}
	if _, e := time.Parse("2006-01-02", t.TransactionDate); e != nil {
		return errors.New("Transaction Date must be YYYY-MM-DD")
	}
	return nil
}

func (t *Transaction) SaveTransaction(db *gorm.DB) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.TransactionUserID).Take(&t.TransactionUser).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}
