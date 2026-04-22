package migration

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type Migrator struct {
	db      *gorm.DB
	path    string
	dialect string
}

func NewMigrator(db *gorm.DB, migrationPath, dialect string) *Migrator {
	return &Migrator{
		db:      db,
		path:    migrationPath,
		dialect: dialect,
	}
}

func (m *Migrator) Run() error {
	dialectPath := filepath.Join(m.path, m.dialect)

	files, err := os.ReadDir(dialectPath)
	if err != nil {
		log.Printf("Failed to read migration directory %s: %v", dialectPath, err)
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			content, err := os.ReadFile(filepath.Join(dialectPath, file.Name()))
			if err != nil {
				log.Printf("Failed to read file %s: %v", file.Name(), err)
				return err
			}

			// Выполняем каждый файл как отдельный запрос
			if err := m.db.Exec(string(content)).Error; err != nil {
				log.Printf("Migration %s failed: %v", file.Name(), err)
				return err
			}
			log.Printf("Migration %s applied successfully", file.Name())
		}
	}
	return nil
}
