package lib

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	ORM *gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(config Config, logger Logger) Database {
	mc := postgres.Config{
		DSN: config.Database.DSN(),
	}

	db, err := gorm.Open(postgres.New(mc), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		SkipDefaultTransaction: true,
		// disable foreign keys
		// specifying foreign keys does not create real foreign key constraints in postgresql
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
			//TablePrefix:   config.Database.TablePrefix + "_",
		},
		// query all fields, and in some cases "*" does not walk the index
		QueryFields: true,
	})

	if err != nil {
		logger.Zap.Fatalf("Error to open database[%s] connection: %v", mc.DSN, err)
	}

	if config.Log.Level == "debug" {
		db = db.Debug()
	}

	logger.Zap.Info("Database connection established")
	return Database{
		ORM: db,
	}
}
