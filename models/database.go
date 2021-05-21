package models

import (
	"fmt"
	"os"

	"github.com/Jarover/sqre/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB //база данных

func init() {

	e := godotenv.Load(utils.GetDir() + "/.env") //Загрузить файл .env
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	//fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	//db.Debug().AutoMigrate(&Account{}, &Contact{}) //Миграция базы данных
}

// возвращает дескриптор объекта DB
func GetDB() *gorm.DB {
	return db
}
