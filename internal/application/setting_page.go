package application

import (
    "Tai/internal/dao"
    "Tai/internal/database"
    "encoding/base64"
    "fmt"
    "github.com/atotto/clipboard"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var serviceProviderList = make(map[int]*dao.ServiceProvider)
var aiModelList []*dao.AiModel
var _settingPage *tview.Flex
var form *tview.Form
var aiModelOptions []string
var aiModelCurOption = 0

func initSettingPage() *tview.Flex {
    _settingPage = tview.NewFlex().
        SetDirection(tview.FlexRow)

    var _serviceProviderList []dao.ServiceProvider
    database.GetDB().Order("id asc").Find(&_serviceProviderList)
    for _, serviceProvider := range _serviceProviderList {
        serviceProviderList[serviceProvider.Id] = &serviceProvider
    }

    database.GetDB().Order("id ASC").Find(&aiModelList)

    form = tview.NewForm()
    for _, serviceProvider := range serviceProviderList {
        apiKeyLabel := fmt.Sprintf("%s API Key", serviceProvider.Name)
        apiKey := ""
        if len(serviceProvider.ApiKey) > 0 {
            decodedStr, err := base64.StdEncoding.DecodeString(serviceProvider.ApiKey)
            if err == nil {
                apiKey = string(decodedStr)
            }
        }
        form.AddPasswordField(apiKeyLabel, apiKey, 80, '*', nil)
        apiKeyInputField := form.GetFormItemByLabel(apiKeyLabel).(*tview.InputField)
        apiKeyInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
            if event.Key() == tcell.KeyCtrlV {
                pastedText, _ := clipboard.ReadAll()
                if len(pastedText) > 0 {
                    apiKeyInputField.SetText(pastedText)
                }
            }
            return event
        })
    }
    form.SetItemPadding(1).SetBorder(true).SetTitle(" Setting ")
    form.SetFieldBackgroundColor(tcell.ColorTeal)
    form.SetButtonStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDarkSeaGreen))
    form.SetButtonActivatedStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkSeaGreen).Background(tcell.ColorWhite))

    for i, aiModel := range aiModelList {
        provider := serviceProviderList[aiModel.ProviderId]
        aiModelOptions = append(aiModelOptions, fmt.Sprintf("%s (Provided by %s)", aiModel.Name, provider.Name))
        if aiModel.Id == activeAiModelId {
            aiModelCurOption = i
        }
    }
    form.AddDropDown("Active AI Model", aiModelOptions, aiModelCurOption, nil)

    form.AddButton("Save", nil).AddButton("Back", nil)
    saveBtn := form.GetButton(form.GetButtonIndex("Save"))
    backBtn := form.GetButton(form.GetButtonIndex("Back"))

    saveBtn.SetSelectedFunc(func() {
        for _, serviceProvider := range serviceProviderList {
            apiKey := form.GetFormItemByLabel(fmt.Sprintf("%s API Key", serviceProvider.Name)).(*tview.InputField).GetText()
            serviceProvider.ApiKey = apiKey
            if len(apiKey) > 0 {
                serviceProvider.ApiKey = base64.StdEncoding.EncodeToString([]byte(apiKey))
            }
            database.GetDB().Save(serviceProvider)
        }
        option, _ := form.GetFormItemByLabel("Active AI Model").(*tview.DropDown).GetCurrentOption()
        dao.UpdateSetting("active_ai_model_id", string(aiModelList[option].Id))
        app.SetActiveAiModel(aiModelList[option].Id)
        aiModelCurOption = option
        saveBtn.Blur()
    })
    backBtn.SetSelectedFunc(func() {
        backBtn.Blur()
        tviewPages.SwitchToPage("main")
        UpdateChatViewTitle()
    })

    _settingPage.AddItem(form, 0, 1, false)
    return _settingPage
}
