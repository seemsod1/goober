package initializers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	models "help/models/entities"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectToDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("SSL_MODE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	DB = db
}

func SyncDB() {
	DB.AutoMigrate(
		&models.UserRole{},
		&models.User{},
		&models.Country{},
		&models.RentLocation{},
		&models.CarType{},
		&models.CarBrand{},
		&models.CarModel{},
		&models.CarPurpose{},
		&models.Car{},
		&models.CarAssignment{},
		&models.RentInfo{},
		&models.UserHistory{},
		&models.CarHistory{},
	)
}
