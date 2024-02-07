package repo

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func New(dbpath string) *repo {
	db, err := initDB(dbpath)
	if err != nil {
		log.Fatal("error connecting to db")
	}
	return &repo{
		db: db,
	}
}

func initDB(dbpath string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return
	}
	db.AutoMigrate(&Terminal{})
	return
}

func (r *repo) GetTerminal(terminal string) (addr string, err error) {
	var t Terminal
	if err = r.db.Where("terminal = ?", terminal).First(&t).Error; err != nil {
		return
	}
	return t.Addr, nil
}
