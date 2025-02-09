package dao

import (
    "Tai/internal/database"
    "log"
    "time"
)

type Session struct {
    dao       *Dao
    Id        int `gorm:"primaryKey;autoIncrement"`
    Name      string
    CreatedAt time.Time
}

func (Session) TableName() string {
    return "session"
}

func InsertSession() {
    InsertSessionWithName("New Chat")
}

func InsertSessionWithName(name string) *Session {
    session := &Session{Name: name}
    result := database.GetDB().Create(session)
    if result.Error != nil {
        log.Fatal("Failed to create a new session: ", result.Error)
    }
    return session
}
