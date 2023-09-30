package db

import (
	"book-manager-server/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	db  *gorm.DB
	cfg config.Config
}

func ConnectToDB(cfg config.Config) (*GormDB, error) {
	ConnectionString := fmt.Sprintf("dbname=%s host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Database.DBName,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password)

	// Create the connection to database
	db, err := gorm.Open(postgres.Open(ConnectionString))
	if err != nil {
		return nil, err
	}

	return &GormDB{
		db:  db,
		cfg: cfg,
	}, nil
}

func (gdb *GormDB) CreateNewSchema() error {
	if err := gdb.db.AutoMigrate(&User{}); err != nil {
		log.Print("Could not migrate user table")
		return err
	}

	if err := gdb.db.AutoMigrate(&Book{}); err != nil {
		log.Print("Could not migrate book table")
		return err
	}
	return nil

}
