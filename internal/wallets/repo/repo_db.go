package repo

import (
	"fmt"

	"github.com/knstch/subtrack-libs/log"
	"gorm.io/gorm"
)

type DBRepo struct {
	lg *log.Logger
	db *gorm.DB
}

func (r *DBRepo) NewDBRepo(db *gorm.DB) *DBRepo {
	if db == nil {
		db = r.db.Session(&gorm.Session{NewDB: true})
	}
	return &DBRepo{
		db: db,
		lg: r.lg,
	}
}

func (r *DBRepo) Transaction(fn func(st Repository) error) error {
	db := r.db.Session(&gorm.Session{NewDB: true})
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := fn(r.NewDBRepo(tx)); err != nil {
			return fmt.Errorf("fn: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("db.Transaction: %w", err)
	}
	return nil
}

func NewDBRepo(lg *log.Logger, db *gorm.DB) *DBRepo {
	return &DBRepo{
		lg: lg,
		db: db,
	}
}
