package domain

import (
    "Tai/internal/dao"
    _ "github.com/fatih/color"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
    Stream   bool      `json:"stream"`
}

type ChatResponse struct {
    Choices []struct {
        Delta struct {
            Content string `json:"content"`
        } `json:"delta"`
    } `json:"choices"`
}

func ConvertConversationToMessage(conversation *dao.Conversation) Message {
    role := "user"
    if conversation.Role == dao.RoleAI {
        role = "assistant"
    }
    return Message{
        Role:    role,
        Content: conversation.Message,
    }
}

func ConvertConversationsToMessages(conversationList []*dao.Conversation) []Message {
    messages := make([]Message, 0, len(conversationList))
    for _, conversation := range conversationList {
        messages = append(messages, ConvertConversationToMessage(conversation))
    }
    return messages
}
