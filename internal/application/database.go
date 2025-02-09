package application

import (
    "Tai/internal/dao"
    "Tai/internal/database"
    "embed"
    _ "embed"
    "fmt"
    "gorm.io/gorm"
    "io/fs"
    "log"
    "strconv"
    "strings"
)

const DataVersion int = 1

//go:embed config/sql/migrations/v*.sql
var migrationFp embed.FS

func initDB() {
    db := database.GetDB()
    var needInitSchema int
    db.Raw("SELECT COUNT(*) > 0 FROM sqlite_master WHERE type = 'table' AND name = ?", "setting").Scan(&needInitSchema)
    if needInitSchema > 0 {
        log.Println("Initialize database...")
    }

    var err error
    oldVersion := 0
    if needInitSchema > 0 {
        oldVersion, err = strconv.Atoi(dao.GetSettingValue("data_version", "0"))
        if err != nil {
            log.Fatal("Failed to get data version: ", err)
        }
    }

    log.Printf("Current data version: %d\n", oldVersion)

    if oldVersion == DataVersion {
        log.Println("Database is up to date.")
        return
    }

    for version := oldVersion; version <= DataVersion; version++ {
        if version <= 0 {
            continue
        }
        log.Println("Migrating to version", version)
        migrate(db, version)
    }
    log.Println("Database migration completed. current version:", DataVersion)
}

func migrate(db *gorm.DB, version int) {
    sqlBytes, err := migrationFp.ReadFile(fmt.Sprintf("config/sql/migrations/v%d.sql", version))
    if err != nil {
        if _, ok := err.(*fs.PathError); ok {
            return
        }
        log.Fatal(err)
    }
    sqlStatements := strings.Split(string(sqlBytes), ";")
    err = nil
    tx := db.Begin()
    defer func(e error) {
        if e != nil {
            tx.Rollback()
        }
    }(err)
    for _, stmt := range sqlStatements {
        stmt = strings.Trim(stmt, " \t\r\n")
        if stmt == "" {
            continue
        }
        result := tx.Exec(stmt)
        if result.Error != nil {
            tx.Rollback()
            log.Fatalf("Error executing statement: %s\nError: %v", stmt, result.Error)
        }
    }
    tx.Commit()
}
