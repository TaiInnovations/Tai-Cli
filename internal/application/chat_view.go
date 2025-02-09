package application

import (
    "Tai/internal/dao"
    "fmt"
    "github.com/fatih/color"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var chatView = tview.NewTextView()
var chatViewWriter = tview.ANSIWriter(chatView)
var forceDrawChatView = false

func initChatView() {
    chatView.
        SetChangedFunc(func() {
            if forceDrawChatView {
                tviewApp.ForceDraw()
            } else {
                tviewApp.Draw()
            }
        }).
        SetDynamicColors(false).
        SetRegions(true).
        SetWordWrap(true).
        SetBorderPadding(1, 1, 1, 1).
        SetTitleColor(tcell.ColorYellow)
    chatView.SetBorder(true)
    app.GetActiveAiModel()
    UpdateChatViewTitle()
}

func UpdateChatViewTitle() {
    title := fmt.Sprintf(" %s ( Provided by %s )", _activeAiModel.Name, app.GetAiServiceProvider(_activeAiModel.ProviderId).Name)
    chatView.SetTitle(" " + title + " ")
}

func addMessageToChatView(conversation *dao.Conversation) {
    blue := color.New(color.FgBlue, color.Bold)
    red := color.New(color.FgRed, color.Bold)
    green := color.New(color.FgGreen, color.Bold)
    var message string
    if conversation.Role == dao.RoleUser {
        message = blue.Sprintf("You:\n")
    } else if conversation.Role == dao.RoleSystem {
        message = red.Sprintf("System:\n")
    } else {
        message = green.Sprintf("AI:\n")
    }
    message += conversation.Message + "\n\n"
    fmt.Fprintf(chatViewWriter, message)
}
