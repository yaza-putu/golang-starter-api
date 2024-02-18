package database

import "gorm.io/gorm"

type MigrateFunc func(db *gorm.DB) error

var (
	Instance      *gorm.DB      // Instance database
	upMigration   []MigrateFunc // create database migration collections
	downMigration []MigrateFunc // rollback database migration collections
	upSeeders     []MigrateFunc // run seeder data
)

// MigrationRegister store to migration collection
func MigrationRegister(up MigrateFunc, down MigrateFunc) {
	upMigration = append(upMigration, up)
	downMigration = append(downMigration, down)
}

// SeederRegister store to seeder collections
func SeederRegister(seed MigrateFunc) {
	upSeeders = append(upSeeders, seed)
}

// MigrationUp can execute all migration
func MigrationUp() error {
	for i := 0; i < len(upMigration); i++ {
		err := upMigration[i](Instance)
		if err != nil {
			return err
		}
	}
	return nil
}

// MigrationDown can execute all migration
func MigrationDown() error {
	for i := 0; i < len(downMigration); i++ {
		err := downMigration[i](Instance)
		if err != nil {
			return err
		}
	}
	return nil
}

// SeederUp can execute all seeder
func SeederUp() error {
	for i := 0; i < len(upSeeders); i++ {
		err := upSeeders[i](Instance)
		if err != nil {
			return err
		}
	}
	return nil
}

// Mock database connection & query
func Mock(db *gorm.DB) {
	Instance = db
}
