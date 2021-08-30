package seed

import (
	"log"

	"github.com/kurniawanew/pfm-kambing/api/models"
	"gorm.io/gorm"
)

var users = []models.User{
	{
		Nickname: "Kurniawan Eko Wasono",
		Email:    "kurniawanew@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Iwan",
		Email:    "iwan@kurniawanew.com",
		Password: "password",
	},
}

var transaction = []models.Transaction{
	{
		TransactionName:   "Makan",
		TransactionAmount: 15000,
		TransactionDate:   "2021-08-30",
	},
	{
		TransactionName:   "Bensin",
		TransactionAmount: 50000,
		TransactionDate:   "2021-08-25",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().Migrator().DropTable(&models.Transaction{}, &models.User{})
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err1 := db.Debug().AutoMigrate(&models.User{}, &models.Transaction{})
	if err1 != nil {
		log.Fatalf("cannot migrate table: %v", err1)
	}

	for i := range users {
		err2 := db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err2 != nil {
			log.Fatalf("cannot seed users table: %v", err2)
		}
		transaction[i].TransactionUserID = users[i].ID

		err3 := db.Debug().Model(&models.Transaction{}).Create(&transaction[i]).Error
		if err3 != nil {
			log.Fatalf("cannot seed posts table: %v", err3)
		}
	}
}
