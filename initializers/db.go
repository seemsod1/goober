package initializers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	entities "help/models/entities"
	"log"
	"os"
	"time"
)

func ConnectToDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("SSL_MODE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	return db
}

func SyncDB(DB *gorm.DB) {
	DB.AutoMigrate(
		&entities.UserRole{},
		&entities.User{},
		&entities.City{},
		&entities.RentLocation{},
		&entities.CarType{},
		&entities.CarBrand{},
		&entities.CarModel{},
		&entities.CarPurpose{},
		&entities.Car{},
		&entities.CarAssignment{},
		&entities.RentInfo{},
		&entities.UserHistory{},
		&entities.CarHistory{},
	)
}

func Migration(DB *gorm.DB) {
	roles := []*entities.UserRole{
		{Name: "Admin"},
		{Name: "Manager"},
		{Name: "User"},
	}
	DB.Create(roles)

	users := []*entities.User{
		{Name: "Admin", BirthDate: time.Date(2004, 12, 21, 12, 54, 12, 2, time.Local), Email: "mr.vadym007@gmail.com", Password: "123", Phone: "380123456781", RoleId: 1},
		{Name: "Manager", BirthDate: time.Date(2004, 12, 21, 12, 54, 12, 2, time.Local), Email: "ex1@", Password: "123", Phone: "380123456782", RoleId: 2},
		{Name: "Manager", BirthDate: time.Date(2004, 12, 21, 12, 54, 12, 2, time.Local), Email: "ex2@", Password: "123", Phone: "380123456783", RoleId: 2},
		{Name: "Manager", BirthDate: time.Date(2004, 12, 21, 12, 54, 12, 2, time.Local), Email: "ex3@", Password: "123", Phone: "380123456784", RoleId: 2},
		{Name: "Manager", BirthDate: time.Date(2004, 12, 21, 12, 54, 12, 2, time.Local), Email: "ex4@", Password: "123", Phone: "380123456785", RoleId: 2},
	}
	DB.Create(users)

	cities := []*entities.City{
		{Name: "Nizhyn"},
		{Name: "Kyiv"},
	}
	DB.Create(cities)
	types := []*entities.CarType{
		{Name: "Sedan"},
		{Name: "Micro"},
		{Name: "Hatchback"},
		{Name: "Universal"},
		{Name: "Liftback"},
		{Name: "Coupe"},
		{Name: "Cabriolet"},
		{Name: "Roadster"},
		{Name: "Targa"},
		{Name: "Limousine"},
		{Name: "Muscle car"},
		{Name: "Sport car"},
		{Name: "Supercar"},
		{Name: "SUV"},
		{Name: "Crossover"},
		{Name: "Pickup"},
		{Name: "Van"},
		{Name: "Minivan"},
		{Name: "Minibus"},
		{Name: "Campervan"},
	}
	DB.Create(types)
	locations := []*entities.RentLocation{
		{CityId: 1, FullAddress: "Shevchenko, 24", UserId: 2},
		{CityId: 1, FullAddress: "Vatutina, 12", UserId: 3},
		{CityId: 2, FullAddress: "Vasilkivska, 53", UserId: 4},
		{CityId: 2, FullAddress: "Khreshchatyk, 87", UserId: 5},
	}
	DB.Create(locations)

	brands := []*entities.CarBrand{
		{Name: "Audi"},
		{Name: "BMW"},
		{Name: "Chevrolet"},
		{Name: "Citroen"},
		{Name: "Dacia"},
		{Name: "Daewoo"},
		{Name: "Dodge"},
		{Name: "Fiat"},
		{Name: "Ford"},
		{Name: "Honda"},
		{Name: "Hyundai"},
		{Name: "Infiniti"},
		{Name: "Jaguar"},
		{Name: "Jeep"},
		{Name: "Kia"},
		{Name: "Lada"},
		{Name: "Lamborghini"},
		{Name: "Land Rover"},
		{Name: "Lexus"},
		{Name: "Mazda"},
		{Name: "Mercedes-Benz"},
		{Name: "Mini"},
		{Name: "Mitsubishi"},
		{Name: "Nissan"},
		{Name: "Opel"},
		{Name: "Peugeot"},
		{Name: "Porsche"},
		{Name: "Renault"},
		{Name: "Rolls-Royce"},
		{Name: "Saab"},
		{Name: "Seat"},
		{Name: "Skoda"},
		{Name: "Smart"},
		{Name: "Subaru"},
		{Name: "Suzuki"},
		{Name: "Tesla"},
		{Name: "Toyota"},
		{Name: "Volkswagen"},
		{Name: "Volvo"},
	}
	DB.Create(brands)

	models := []*entities.CarModel{
		{Name: "A1", BrandId: 1},
		{Name: "A3", BrandId: 1},
		{Name: "A4", BrandId: 1},
		{Name: "A5", BrandId: 1},
		{Name: "A6", BrandId: 1},
		{Name: "X5", BrandId: 2},
		{Name: "X6", BrandId: 2},
	}
	DB.Create(models)

	purposes := []*entities.CarPurpose{
		{Name: "Business"},
		{Name: "Family"},
		{Name: "Sport"},
		{Name: "Off-road"},
		{Name: "City"},
		{Name: "Economy"},
		{Name: "Premium"},
		{Name: "Luxury"},
	}
	DB.Create(purposes)

	cars := []*entities.Car{
		{TypeId: 1, ModelId: 3, Bags: 2, Passengers: 5, Year: 2021, Plate: "AA1234AA", Price: 113, Color: "Black", LocationId: 1},
		{TypeId: 14, ModelId: 7, Bags: 4, Passengers: 5, Year: 2021, Plate: "AA1235AA", Price: 200, Color: "Red", LocationId: 2},
	}
	DB.Create(cars)

	assignments := []*entities.CarAssignment{
		{CarId: 1, PurposeId: 1},
		{CarId: 1, PurposeId: 2},
		{CarId: 1, PurposeId: 5},
		{CarId: 1, PurposeId: 6},
		{CarId: 2, PurposeId: 1},
		{CarId: 2, PurposeId: 2},
		{CarId: 2, PurposeId: 5},
		{CarId: 2, PurposeId: 6},
	}
	DB.Create(assignments)

}
