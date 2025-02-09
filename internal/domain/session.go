package domain

import (
    "Tai/internal/dao"
)

type Session struct {
    Id               int
    Name             string
    ConversationList []*dao.Conversation
}

func NewSession(id int, name string) Session {
    session := Session{
        Id:               id,
        Name:             name,
        ConversationList: nil,
    }
    return session
}

func NewSessionByDao(record *dao.Session) Session {
    return NewSession(record.Id, record.Name)
}

func (s *Session) InitConversationList() []*dao.Conversation {
    s.ConversationList = dao.GetConversationListBySessionId(s.Id)
    return s.ConversationList
}

func (s *Session) GetConversationList() []*dao.Conversation {
    if s.ConversationList == nil {
        s.InitConversationList()
    }
    return s.ConversationList
}
