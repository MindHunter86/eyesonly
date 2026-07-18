package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDB struct {
	db *gorm.DB
}

func Connect() (_ *SqliteDB, e error) {
	var db *gorm.DB
	if db, e = gorm.Open(sqlite.Open("test.db"), &gorm.Config{}); e != nil {
		return
	}

	// connection settings setup
	// todo - Connection Pool -https://gorm.io/docs/generic_interface.html

	// !! TODO - TO DELETE
	db.Select("1")

	return
}
