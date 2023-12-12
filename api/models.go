package api

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	ID          int    `gorm:"primaryKey" json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	SSN         string `json:"ssn"`
	Title       string `json:"job_title"`
	Company     string `json:"company"`
	SecondaryId string `json:"secondary_id"`
}

func DbConn(
	host string, port int, user string,
	password string, dbname string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Profile{})
	if err != nil {
		panic(err)
	}
	return db
}
