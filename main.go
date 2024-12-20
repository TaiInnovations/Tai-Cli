package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// 可用模型列表
var availableModels = []string{
	"google/gemini-2.0-flash-exp:free", // 默认模型
	"google/gemini-exp-1206:free",
	"google/gemini-exp-1121:free",
	"google/learnlm-1.5-pro-experimental:free",
	"google/gemini-exp-1114:free",
	"google/gemini-2.0-flash-thinking-exp:free",
}

var currentModel string

// Message 定义聊天消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 定义请求结构
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse 定义响应结构
type ChatResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func init() {
	// 设置控制台输出编码
	if runtime.GOOS == "windows" {
		// Windows 系统特殊处理
		cmd := exec.Command("chcp", "65001")
		cmd.Run()
	}
}

func loadAPIKey() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("OPENROUTER_API_KEY")
	}

	keyFile := filepath.Join(homeDir, ".gemini_api_key")
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return os.Getenv("OPENROUTER_API_KEY")
	}

	return strings.TrimSpace(string(data))
}

func saveAPIKey(apiKey string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	keyFile := filepath.Join(homeDir, ".gemini_api_key")
	return os.WriteFile(keyFile, []byte(apiKey), 0600)
}

func validateAPIKey(apiKey string) bool {
	if apiKey == "" {
		return false
	}
	if strings.Contains(apiKey, "{") || strings.Contains(apiKey, "}") {
		return false
	}
	return true
}

func chatWithModel(messages []Message, apiKey string, model string) (*http.Response, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	chatReq := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	return client.Do(req)
}

func showModels() {
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("\n%s\n", yellow("可用模型列表："))

	for i, model := range availableModels {
		if model == currentModel {
			fmt.Printf("%d. %s (当前使用)\n", i+1, green(model))
		} else {
			fmt.Printf("%d. %s\n", i+1, model)
		}
	}
}

func main() {
	// 确保控制台输出使用 UTF-8 编码
	os.Stdout.WriteString("\xEF\xBB\xBF") // 添加 UTF-8 BOM

	currentModel = availableModels[0]
	var messages []Message

	blue := color.New(color.FgBlue, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed)

	blue.Println("\n欢迎使用 AI 聊天!")
	green.Printf("当前使用的模型: %s\n", currentModel)
	fmt.Println("\n命令提示：")
	fmt.Println("• 输入 'quit' 或 'exit' 结束对话")
	fmt.Println("• 输入 'new' 或 'clear' 开启新会话")
	fmt.Println("• 输入 'models' 查看可用模型")
	fmt.Println("• 输入 'switch <数字>' 切换模型 (例如: switch 2)")
	fmt.Println("• 输入 'setkey' 设置 API key\n")

	apiKey := loadAPIKey()
	if !validateAPIKey(apiKey) {
		fmt.Println("警告: API key 无效！")
		fmt.Println("请使用 'setkey <your-api-key>' 命令设置 API key")
		apiKey = ""
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		green.Print("你: ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()

		// 处理命令
		if strings.HasPrefix(userInput, "setkey ") {
			parts := strings.SplitN(userInput, " ", 2)
			if len(parts) == 2 {
				newKey := strings.TrimSpace(parts[1])
				if validateAPIKey(newKey) {
					apiKey = newKey
					if err := saveAPIKey(apiKey); err != nil {
						red.Printf("保存 API key 失败: %v\n", err)
					} else {
						fmt.Println("API key 已保存到本地文件")
					}
				} else {
					fmt.Println("无效的 API key！")
				}
			}
			continue
		}

		// 验证 API key
		if !validateAPIKey(apiKey) && !strings.HasPrefix(userInput, "setkey") {
			fmt.Println("请先设置有效的 API key！使用命令: setkey <your-api-key>")
			continue
		}

		// 处理其他命令
		switch strings.ToLower(userInput) {
		case "quit", "exit":
			fmt.Println("再见！")
			return
		case "new", "clear":
			messages = nil
			fmt.Println("\n已开启新会话！")
			continue
		case "models":
			showModels()
			continue
		case "setkey":
			fmt.Print("请输入新的 API key: ")
			scanner.Scan()
			newKey := strings.TrimSpace(scanner.Text())
			if validateAPIKey(newKey) {
				apiKey = newKey
				if err := saveAPIKey(apiKey); err != nil {
					red.Printf("保存 API key 失败: %v\n", err)
				} else {
					fmt.Println("API key 已保存到本地文件")
				}
			} else {
				fmt.Println("无效的 API key！")
			}
			continue
		}

		if strings.HasPrefix(strings.ToLower(userInput), "switch ") {
			// 处理切换模型命令
			var modelNum int
			if _, err := fmt.Sscanf(userInput[7:], "%d", &modelNum); err == nil {
				if modelNum > 0 && modelNum <= len(availableModels) {
					currentModel = availableModels[modelNum-1]
					green.Printf("\n已切换到模型: %s\n", currentModel)
				} else {
					fmt.Println("\n无效的模型编号！")
				}
			} else {
				fmt.Println("\n无效的命令格式！请使用 'switch <数字>' 格式。")
			}
			continue
		}

		// 添加用户消息
		messages = append(messages, Message{
			Role:    "user",
			Content: userInput,
		})

		// 发送请求并处理响应
		resp, err := chatWithModel(messages, apiKey, currentModel)
		if err != nil {
			red.Printf("\nError: %v\n", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			blue.Print("\nAI: ")

			reader := bufio.NewReader(resp.Body)
			aiResponse := ""

			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					red.Printf("Error reading response: %v\n", err)
					break
				}

				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "data: ") {
					continue
				}

				line = strings.TrimPrefix(line, "data: ")
				if line == "[DONE]" {
					fmt.Println()
					messages = append(messages, Message{
						Role:    "assistant",
						Content: aiResponse,
					})
					break
				}

				var chatResp ChatResponse
				if err := json.Unmarshal([]byte(line), &chatResp); err != nil {
					red.Printf("Error parsing JSON: %v\n", err)
					continue
				}

				if len(chatResp.Choices) > 0 {
					content := chatResp.Choices[0].Delta.Content
					aiResponse += content
					fmt.Print(content)
				}
			}
		} else {
			red.Printf("\nError: 未收到有效响应 (状态码: %d)\n", resp.StatusCode)
		}

		resp.Body.Close()
	}
}
