package application

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var componentFocusSwitchList = []tview.Primitive{sessionListView, inputView}

func initMainPage() *tview.Flex {
    initChatView()
    initSessionList()
    initInputView()
    initTipView()

    // main mainPage
    mainPage := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(sessionListView, 0, 1, false).
        AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
            AddItem(chatView, 0, 1, false).
            AddItem(inputView, inputViewRows, 0, true).
            AddItem(tipView, 2, 0, false),
            0, 5, true,
        )
    mainPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyTab:
            for k, v := range componentFocusSwitchList {
                if tviewApp.GetFocus() == v {
                    nextIndex := (k + 1) % len(componentFocusSwitchList)
                    tviewApp.SetFocus(componentFocusSwitchList[nextIndex])
                    break
                }
            }
            return nil
        case tcell.KeyCtrlC:
            return nil
        case tcell.KeyEscape:
            tviewApp.Stop()
            return nil
        }
        return event
    })
    return mainPage
}
