package application

import (
    "Tai/internal/dao"
    "Tai/internal/domain"
    "bufio"
    "bytes"
    "encoding/base64"
    "encoding/json"
    "github.com/fatih/color"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "io"
    "log"
    "net/http"
    "strings"
)

var inputView = tview.NewTextArea()
var inputViewRows = 6

func initInputView() {
    green := color.New(color.FgGreen, color.Bold)
    magenta := color.New(color.FgMagenta, color.Bold)
    inputView.SetLabel("> ").
        SetPlaceholder("Press Enter to send messages.\nUse \"/help\" to check command usages.").
        SetPlaceholderStyle(tcell.StyleDefault.Foreground(tcell.ColorGray)).
        SetSize(inputViewRows, 0)
    inputView.SetBorder(true).SetBorderPadding(0, 0, 1, 0)

    inputView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        inputView.SetDisabled(true)
        defer func() {
            inputView.SetDisabled(false)
        }()
        var conversation *dao.Conversation
        userMessage := strings.Trim(inputView.GetText(), "\t\n\r")
        if event.Key() == tcell.KeyCtrlJ {
            inputView.SetText(inputView.GetText()+"\n", true)
            return nil
        } else if event.Key() == tcell.KeyEnter {
            if len(userMessage) == 0 {
                return nil
            }
            inputView.SetText("", true)
            if strings.HasPrefix(userMessage, "/") {
                commands := strings.SplitN(userMessage, " ", 2)
                directive := strings.ToLower(strings.TrimLeft(commands[0], "/"))
                systemMessage := ""
                switch directive {
                case "help", "helps":
                    systemMessage += "/help: Show command usages.\n" +
                        "/setting or F2: Go to setting window.\n" +
                        "/new <name>: Create a new session.\n" +
                        "/rename <new name>: Rename the current session.\n" +
                        "/delete or /del: Delete current session.\n" +
                        "F1: Go to chat window.\n" +
                        "Ctrl+J: New line.\n" +
                        "Ctrl+Enter: New line. (Unsupported on certain operating system.)\n" +
                        "/exit or /quit: Exit the application."
                case "setting":
                    tviewPages.SwitchToPage("setting")
                case "new":
                    sessionName := "New Chat"
                    if len(commands) >= 2 {
                        sessionName = commands[1]
                    }
                    createSession(sessionName)
                case "rename":
                    if len(commands) < 2 {
                        systemMessage += "Usage: /rename <new name>"
                    } else {
                        sessionName := commands[1]
                        changeSessionName(activeSessionIndex, sessionName)
                        systemMessage += "Session name changed successfully."
                    }
                case "delete", "del":
                    sessionsCount := len(sessions)
                    if sessionsCount > 1 {
                        deleteSession(activeSessionIndex)
                    } else if sessionsCount > 0 {
                        oldActiveSessionIndex := activeSessionIndex
                        createSession("New Chat")
                        deleteSession(oldActiveSessionIndex + 1)
                    } else {
                        createSession("New Chat")
                    }
                case "exit", "quit":
                    tviewApp.Stop()
                    return nil
                default:
                    systemMessage += "Unknown command. Use /help to check supported commands."
                }
                inputView.SetText("", true)
                if len(systemMessage) > 0 {
                    conversation = &dao.Conversation{
                        Role:    dao.RoleSystem,
                        Message: magenta.Sprintf(systemMessage),
                    }
                    addMessageToChatView(conversation)
                    sessions[activeSessionIndex].ConversationList = append(sessions[activeSessionIndex].ConversationList, conversation)
                }
            } else {
                conversation = dao.InsertConversation(sessions[activeSessionIndex].Id, dao.RoleUser, userMessage)
                addMessageToChatView(conversation)
                sessions[activeSessionIndex].ConversationList = append(sessions[activeSessionIndex].ConversationList, conversation)
                green.Fprintf(chatViewWriter, "AI:\n")
                responseConversation, chatResponseErrExist := chatWithModel()
                if !chatResponseErrExist {
                    responseConversation = dao.InsertConversation(conversation.SessionId, dao.RoleAI, responseConversation.Message)
                }
                sessions[activeSessionIndex].ConversationList = append(sessions[activeSessionIndex].ConversationList, responseConversation)
            }
            chatView.ScrollToEnd()
            return nil
        }
        return event
    })
}

func chatWithModel() (*dao.Conversation, bool) {
    forceDrawChatView = true
    chatResponseErrExist := false
    defer func() {
        forceDrawChatView = false
    }()
    conversation := &dao.Conversation{
        SessionId: sessions[activeSessionIndex].Id,
        Role:      dao.RoleAI,
        Message:   "",
    }
    aiModel := app.GetActiveAiModel()
    var conversationList []*dao.Conversation
    for _, c := range sessions[activeSessionIndex].ConversationList {
        if c.Role != dao.RoleSystem {
            conversationList = append(conversationList, c)
        }
    }
    red := color.New(color.FgRed)

    serviceProvider := app.GetAiServiceProvider(aiModel.ProviderId)

    chatReq := domain.ChatRequest{
        Model:    aiModel.Name,
        Messages: domain.ConvertConversationsToMessages(conversationList),
        Stream:   true,
    }

    jsonData, err := json.Marshal(chatReq)
    if err != nil {
        appendChatResponseMessage(red.Sprintf(err.Error()), conversation)
        return conversation, chatResponseErrExist
    }

    req, err := http.NewRequest("POST", serviceProvider.Url, bytes.NewBuffer(jsonData))
    if err != nil {
        appendChatResponseMessage(red.Sprintf(err.Error()), conversation)
        return conversation, chatResponseErrExist
    }
    apiKey := ""
    if len(serviceProvider.ApiKey) > 0 {
        decodedStr, err := base64.StdEncoding.DecodeString(serviceProvider.ApiKey)
        if err == nil {
            apiKey = string(decodedStr)
        }
    }

    if len(apiKey) <= 0 {
        appendChatResponseMessage(red.Sprintln("\nError: Api key is required. Please go to the setting window to set the api key."), conversation)
        return conversation, chatResponseErrExist
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        appendChatResponseMessage(red.Sprintf("\nError: %v\n", err), conversation)
        return conversation, chatResponseErrExist
    }
    defer resp.Body.Close()

    var messages []domain.Message

    if resp.StatusCode == http.StatusOK {

        reader := bufio.NewReader(resp.Body)
        aiResponse := ""

        for {
            line, err := reader.ReadString('\n')
            if err == io.EOF {
                break
            }
            if err != nil {
                appendChatResponseMessage(red.Sprintf("Error reading response: %v\n", err), conversation)
                break
            }

            line = strings.TrimSpace(line)
            if !strings.HasPrefix(line, "data: ") {
                continue
            }

            line = strings.TrimPrefix(line, "data: ")
            if line == "[DONE]" {
                appendChatResponseMessage("\n", conversation)
                conversation.Message += red.Sprintln()
                messages = append(messages, domain.Message{
                    Role:    "assistant",
                    Content: aiResponse,
                })
                break
            }

            var chatResp domain.ChatResponse
            if err := json.Unmarshal([]byte(line), &chatResp); err != nil {
                appendChatResponseMessage(red.Sprintf("Error parsing JSON: %v\n", err), conversation)
                chatResponseErrExist = true
                continue
            }

            log.Println(chatResp)
            if len(chatResp.Choices) > 0 {
                content := chatResp.Choices[0].Delta.Content
                aiResponse += content
                appendChatResponseMessage(content, conversation)
                continue
            }
        }
    } else {
        appendChatResponseMessage(red.Sprintf("Error: No response. (status code: %d, error message)\n", resp.StatusCode), conversation)
        chatResponseErrExist = true
    }
    appendChatResponseMessage("\n", conversation)
    conversation.Message += red.Sprintln()
    return conversation, chatResponseErrExist
}

func appendChatResponseMessage(message string, conversation *dao.Conversation) {
    chatViewWriter.Write([]byte(message))
    conversation.Message += message
}
