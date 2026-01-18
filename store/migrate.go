package store

import "gorm.io/gorm"

// autoMigrateCreateOnly performs a "create-only" migration for SQLite:
//   - if a table already exists, skip AutoMigrate to avoid SQLite rebuild via *_temp
//     which may fail on legacy databases.
//   - if a table does not exist, create it via AutoMigrate.
//
// For non-SQLite dialects, it falls back to normal AutoMigrate behavior.
func autoMigrateCreateOnly(db *gorm.DB, models ...interface{}) error {
	if db.Dialector.Name() != "sqlite" {
		return db.AutoMigrate(models...)
	}

	for _, model := range models {
		if db.Migrator().HasTable(model) {
			continue
		}
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}
