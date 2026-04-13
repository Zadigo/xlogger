package backend

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UrlModel struct {
	ID uint `gorm:"primaryKey"`
	// Unique URL or path to identify the log
	Url string `gorm:"uniqueIndex,not null"`
	// Soft deleted at timestamp
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Updated at timestamp
	UpdatedAt time.Time
	// Created at timestamp
	CreatedAt time.Time
}

type LogModel struct {
	ID uint `gorm:"primaryKey"`
}

func MigrateTables() {}

func CreateTables() {}

// DbConnection creates and returns a connection to the database
func DbConnection() (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open("xlogger.sqlite"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// Create tables if they do not exist
	err = conn.AutoMigrate(&UrlModel{}, &LogModel{})

	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	return conn, err
}

func CreateLogSql() {}

func DeleteLogSql() {}

func GetLogSql() {}

func GetAllLogsSql() {}

func UpdateLogSql() {}
