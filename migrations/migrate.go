package migrations

import (
	"log"
	"os"
	"time"

	"github.com/Cprime50/RegulerClub/database"
	"github.com/Cprime50/RegulerClub/models"
)

func Migrate() {
	log.Printf("Migrations Started")
	startTime := time.Now()
	database.Db.AutoMigrate(&models.Role{})
	database.Db.AutoMigrate(&models.User{})
	database.Db.AutoMigrate(&models.OTP{})
	database.Db.AutoMigrate(&models.OTPVerificationType{})

	err := seedData() // default data being added into the database upon migration
	if err != nil {
		log.Fatal(err)
	}
	log.Println("seeding data complete")
	elapsed := time.Since(startTime)
	log.Printf("Migrate completed in %s", elapsed)

}

var userModel models.User

// adding some default user data and roles into the db
func seedData() error {

	var roles = []models.Role{{ID: 1, Name: "admin", Description: "Administrator role"}, {ID: 2, Name: "mod", Description: "Moderator role"}, {ID: 3, Name: "prime user", Description: "prime user"}, {ID: 4, Name: "user", Description: "Normal user role"}}
	var user = &models.User{ID: 1, Username: os.Getenv("ADMIN_USERNAME"), Email: os.Getenv("ADMIN_EMAIL"), Password: userModel.SetPassword(os.Getenv("ADMIN_PASSWORD")), RoleID: 1, Verified: true}
	var OtpType = []models.OTPVerificationType{{ID: 1, Type: "Email verification", Ttl: 1}, {ID: 2, Type: "Password Reset", Ttl: 10}}
	database.Db.Save(&roles)
	database.Db.Save(&user)
	database.Db.Save(&OtpType)

	return nil

}
