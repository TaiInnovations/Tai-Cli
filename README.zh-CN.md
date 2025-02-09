# Tai

使用 OpenRouter 和 Google AI Studio / DeepSeek 免费访问 Gemini / DeepSeek API 的简单工具。

## 快速开始

### 安装方法

通过 GitHub 安装脚本一键安装：
```bash
bash <(curl -Lk https://github.com/tasselx/free-gemini/releases/download/latest/install.sh)
```

### 配置步骤

1. **获取 Google API Key**
   - 访问 [Google AI Studio](https://aistudio.google.com/app/apikey)
   - 创建并复制 API Key

2. **获取 DeepSeek API Key**
   - 访问 [DeepSeek](https://platform.deepseek.com/api_keys)
   - 创建并复制 API Key

3. **配置 OpenRouter**
   - 访问 [OpenRouter 集成页面](https://openrouter.ai/settings/integrations)
   - 填入上一步获取的 Google AI Studio / DeepSeek API Key
   - 访问 [OpenRouter Keys 页面](https://openrouter.ai/settings/keys)
   - 创建并获取 OpenRouter API Key

4. **运行构建**
```bash
# 全平台
make
# 指定平台
make windows_amd64 linux_amd64
```

## 注意事项
- 请确保妥善保管您的 API Keys
- 使用前请阅读相关服务条款

## 贡献
欢迎提交 Issues 和 Pull Requests！