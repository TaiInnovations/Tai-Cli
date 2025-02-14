package application

import (
    "Tai/internal/dao"
    "Tai/internal/domain"
    "github.com/rivo/tview"
)

var activeSessionIndex int
var sessions []domain.Session
var sessionListView = tview.NewList()

func initSessionList() {
    sessionListView.SetBorder(true).SetBorderPadding(1, 1, 1, 1)
    sessionListView.SetTitle("  Session List ")
    sessionListView.SetChangedFunc(func(i int, mainText string, secondaryText string, key rune) {
        sessions[i].GetConversationList()
        chatView.Clear()
        for _, conversation := range sessions[i].ConversationList {
            if conversation.Role != dao.RoleSystem {
                addMessageToChatView(conversation)
            }
        }
        chatView.ScrollToEnd()
        activeSessionIndex = i
    })
    var sessionDaos []dao.Session
    db.Find(&sessionDaos)
    for _, sessionDao := range sessionDaos {
        insertSession(&sessionDao)
    }
}

func insertSession(sessionDao *dao.Session) {
    session := domain.NewSessionByDao(sessionDao)
    sessions = append([]domain.Session{session}, sessions...)
    activeSessionIndex = 0
    sessionListView.InsertItem(activeSessionIndex, sessionDao.Name, "", 0, nil)
    sessionListView.SetCurrentItem(activeSessionIndex)
}

func createSession(name string) {
    sessionDao := dao.Session{
        Name: name,
    }
    db.Create(&sessionDao)
    insertSession(&sessionDao)
}

func deleteSession(index int) {
    sessionCount := sessionListView.GetItemCount()
    if sessionCount <= 1 {
        return
    }
    db.Delete(&dao.Session{Id: sessions[index].Id})
    db.Delete(&dao.Conversation{}, "session_id = ?", sessions[index].Id)
    sessions = append(sessions[:index], sessions[index+1:]...)
    sessionListView.RemoveItem(index)
    if index == activeSessionIndex {
        activeSessionIndex = 0
        if index > 0 {
            activeSessionIndex = index - 1
        }
        sessionListView.SetCurrentItem(activeSessionIndex)
    }
}

func changeSessionName(index int, newName string) {
    session := sessions[index]
    sessionDao := dao.Session{
        Id:   session.Id,
        Name: newName,
    }
    db.Save(&sessionDao)
    session.Name = newName
    sessionListView.SetItemText(index, session.Name, "")
}
