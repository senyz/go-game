package migration

import (
    "embed"
    "io/fs"
    "log"

    "gorm.io/gorm"
)

//go:embed *.sql
var migrationFiles embed.FS

type Migrator struct {
    db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
    return &Migrator{db: db}
}

func (m *Migrator) Run() error {
    // Получаем все файлы миграций
    files, err := fs.ReadDir(migrationFiles, ".")
    if err != nil {
        return err
    }

    for _, file := range files {
        content, err := migrationFiles.ReadFile(file.Name())
        if err != nil {
            return err
        }

        // Выполняем миграцию
        if err := m.db.Exec(string(content)).Error; err != nil {
            log.Printf("Migration %s failed: %v", file.Name(), err)
            return err
        }
        log.Printf("Migration %s applied successfully", file.Name())
    }
    return nil
}
