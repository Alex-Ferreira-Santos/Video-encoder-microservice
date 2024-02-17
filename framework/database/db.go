package database

import (
	"log"

	"github.com/alex-ferreira-santos/encoder/domain"
	"github.com/jinzhu/gorm"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "Test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (database *Database) Connect() (*gorm.DB, error) {
	var err error

	if database.Env != "Test" {
		database.Db, err = gorm.Open(database.DbType, database.Dsn)
	}
	if database.Env == "Test" {
		database.Db, err = gorm.Open(database.DbTypeTest, database.DsnTest)
	}

	if err != nil {
		return nil ,err
	}

	if database.Debug {
		database.Db.LogMode(true)
	}

	if database.AutoMigrateDb {
		database.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
	}

	return database.Db, nil
}
