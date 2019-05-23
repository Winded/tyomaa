package db

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/winded/tyomaa/backend/util"
)

var (
	errSystemStartingUp = errors.New("pq: the database system is starting up")
)

var Instance *gorm.DB

func initSqlite() (*gorm.DB, error) {
	return gorm.Open("sqlite3", ":memory:")
}

func initPostgres() (*gorm.DB, error) {
	host := util.EnvOrDefault("DB_HOST", "localhost")
	port := util.EnvOrDefault("DB_PORT", "5432")
	dbname := util.EnvOrDefault("DB_DATABASE", "tyomaa")
	user := util.EnvOrDefault("DB_USER", "tyomaa")
	password := util.EnvOrDefault("DB_PASSWORD", "tyomaa")

	params := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable intervalstyle=iso_8601",
		host, port, dbname, user, password)

	tries := 0
	for {
		conn, err := gorm.Open("postgres", params)

		if err != nil && tries < 5 {
			fmt.Println("Connection failed, retrying...")
			time.Sleep(500 * time.Millisecond)
			tries++
			continue
		}

		return conn, err
	}
}

func Init() (instance *gorm.DB, err error) {
	inMemory, _ := strconv.Atoi(os.Getenv("DB_IN_MEMORY"))
	if inMemory <= 0 {
		instance, err = initPostgres()
	} else {
		instance, err = initSqlite()
	}

	if err == nil {
		Instance = instance
	}

	return instance, err
}

func AutoMigrate(db *gorm.DB) {
	var (
		user  User
		entry TimeEntry
	)

	db.AutoMigrate(&user)
	db.AutoMigrate(&entry)
}
