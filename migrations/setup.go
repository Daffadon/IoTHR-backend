package migrations

import (
	"IoTHR-backend/db"
	"IoTHR-backend/models"
)

func Init() {
	db := db.GetDB()
	db.AutoMigrate(&models.User{})
}
