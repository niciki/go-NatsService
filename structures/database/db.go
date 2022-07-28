package db

import (
	"log"

	ls "github.com/niciki/go-NatsService/structures/localStore"
	so "github.com/niciki/go-NatsService/structures/structOrder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseRecord struct {
	Order_uid string
	Order     so.Order
}

type Database struct {
	db *gorm.DB
}

func InitDb(port string) (Database, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=user password=user dbname=db_wb_l0 port=" + port))
	if err != nil {
		return Database{}, err
	}
	err = db.AutoMigrate(&DatabaseRecord{})
	if err != nil {
		return Database{}, err
	}
	return Database{db}, nil
}

func (d *Database) UploadCache(cache ls.Store) error {
	var databaseData []DatabaseRecord
	err := d.db.Find(&databaseData)
	for _, rec := range databaseData {
		cache.Add(rec.Order)
	}
	return err.Error
}

func (d *Database) AddRecord(rec ...so.Order) error {
	log.Printf("Add to db record: %v", rec)
	for _, r := range rec {
		err := d.db.Create(r)
		if err.Error != nil {
			return err.Error
		}
	}
	return nil
}
