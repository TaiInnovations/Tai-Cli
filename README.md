# Free Gemini

使用 OpenRouter 和 Google AI Studio 免费访问 Gemini API 的简单工具。

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

2. **配置 OpenRouter**
   - 访问 [OpenRouter 集成页面](https://openrouter.ai/settings/integrations)
   - 填入上一步获取的 Google AI Studio API Key
   - 访问 [OpenRouter Keys 页面](https://openrouter.ai/settings/keys)
   - 创建并获取 OpenRouter API Key

3. **运行构建**
```bash
./build.sh
```

## 注意事项
- 请确保妥善保管您的 API Keys
- 使用前请阅读相关服务条款

## 贡献
欢迎提交 Issues 和 Pull Requests！