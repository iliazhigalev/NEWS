package database

import (
	"fmt"
	"log"
	"news/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("DB_USER: %s, DB_PASSWORD: %s, DB_NAME: %s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

		log.Fatal("FAILED to connect to database. \n", err)
		os.Exit(2)
	}

	// сообщение о том, что мы подключились
	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	//воспользуемся функцией автоматической миграции, чтобы создать таблицы из нашей модели
	log.Println("running migrations")
	db.AutoMigrate(&models.Article{})

	DB = Dbinstance{
		Db: db,
	}

}
