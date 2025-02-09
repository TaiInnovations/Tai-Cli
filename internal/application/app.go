package application

import (
    "Tai/internal/dao"
    "Tai/internal/database"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "gorm.io/gorm"
    "log"
    "strconv"
)

var db *gorm.DB
var tviewApp *tview.Application
var tviewPages *tview.Pages
var activeAiModelId int
var _activeAiModel *dao.AiModel

type App struct {
    db                  *gorm.DB
    serviceProviderList map[int]*dao.ServiceProvider
}

var app = &App{}

func Run() {
    db = database.GetDB()
    defer database.CloseDB()
    initDB()

    app.serviceProviderList = make(map[int]*dao.ServiceProvider)
    activeAiModelId, _ = strconv.Atoi(dao.GetSettingValue("active_ai_model_id", "1"))

    var err error

    tviewApp = tview.NewApplication()
    tviewPages = tview.NewPages()

    log.Println("Initializing application...")
    mainPage := initMainPage()
    tviewPages.AddPage("main", mainPage, true, true)
    settingPage := initSettingPage()
    tviewPages.AddPage("setting", settingPage, true, false)

    tviewApp.EnableMouse(true)

    tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyF1:
            tviewPages.SwitchToPage("main")
            UpdateChatViewTitle()
        case tcell.KeyF2:
            tviewPages.SwitchToPage("setting")
        case tcell.KeyCtrlC:
            return nil
        case tcell.KeyEscape:
            tviewApp.Stop()
            return nil
        }
        return event
    })

    if err = tviewApp.SetRoot(tviewPages, true).Run(); err != nil {
        panic(err)
    }
}

func (app *App) GetAiServiceProvider(id int) *dao.ServiceProvider {
    var serviceProvider *dao.ServiceProvider
    var exists bool
    if serviceProvider, exists = app.serviceProviderList[id]; !exists {
        serviceProvider = &dao.ServiceProvider{}
        database.GetDB().Where("id = ?", id).First(serviceProvider)
        app.serviceProviderList[id] = serviceProvider
    }
    return serviceProvider
}

func (app *App) SetActiveAiModel(id int) *dao.AiModel {
    value := strconv.Itoa(id)
    dao.UpdateSetting("active_ai_model_id", value)
    activeAiModelId = id
    _activeAiModel = dao.GetAiModelById(id)
    return _activeAiModel
}

func (app *App) GetActiveAiModel() *dao.AiModel {
    if _activeAiModel == nil {
        _activeAiModel = dao.GetAiModelById(activeAiModelId)
    }
    return _activeAiModel
}
