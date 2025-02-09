package database

import (
    "log"
    "os"
    "path/filepath"
    "sync"

    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

var (
    db   *gorm.DB
    once sync.Once
)

func getDataDir() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal("Failed to get home directory:", err)
    }
    return filepath.Join(homeDir, "TaiInnovation/TaiCli/data")
}

func getDataPath() string {
    return filepath.Join(getDataDir(), "data.db")
}

func GetDB() *gorm.DB {
    once.Do(func() {
        var err error

        if err := os.MkdirAll(getDataDir(), os.ModePerm); err != nil {
            log.Fatal("Failed to create directory:", err)
        }

        db, err = gorm.Open(sqlite.Open(getDataPath()), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect database:", err)
        }
    })
    return db
}

func CloseDB() {
    db, err := GetDB().DB()
    if err != nil {
        log.Fatal("Failed to close DB connection: ", err)
    }
    db.Close()
}
