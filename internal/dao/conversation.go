package dao

import (
    "Tai/internal/database"
    "log"
    "time"
)

type Role int

const (
    RoleUser   Role = 1
    RoleAI          = 2
    RoleSystem      = 3
)

type Conversation struct {
    dao       *Dao
    Id        int `gorm:"primaryKey;autoIncrement"`
    SessionId int `gorm:"index:idx_session_id"`
    Role      Role
    Message   string
    CreatedAt time.Time
}

func (Conversation) TableName() string {
    return "conversation"
}

func GetConversationListBySessionId(sessionId int) []*Conversation {
    var conversationList []Conversation
    database.GetDB().Where("session_id = ?", sessionId).Find(&conversationList)
    var result []*Conversation
    for _, conversation := range conversationList {
        if conversation.Role != RoleSystem {
            result = append(result, &conversation)
        }
    }
    return result
}

func InsertConversation(sessionId int, role Role, message string) *Conversation {
    conversation := &Conversation{
        SessionId: sessionId,
        Role:      role,
        Message:   message,
    }
    if role != RoleSystem {
        result := database.GetDB().Create(conversation)
        if result.Error != nil {
            log.Fatal("Failed to create a new chat history: ", result.Error)
        }
    }
    return conversation
}
