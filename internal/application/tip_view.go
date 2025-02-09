package application

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var tipView = tview.NewTextView()

func initTipView() {
    tipView.SetSize(1, 0).
        SetTextColor(tcell.ColorMediumSeaGreen).
        SetText("Tips: Use /help to get help.")
    tipView.SetBorderPadding(0, 1, 1, 0)
}
