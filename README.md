# Tai

A simple tool to access Gemini / DeepSeek API for free using OpenRouter and Google AI Studio / DeepSeek.

[中文文档](README.zh-CN.md)

## Quick Start

### Installation

Install via script:
```bash
bash <(curl -Lk https://github.com/TaiInnovations/Tai-Cli/releases/download/latest/install.sh)
```

### Setup

1. **Get Google API Key**
   - Visit [Google AI Studio](https://aistudio.google.com/app/apikey)
   - Create and copy API Key

2. **Get DeepSeek API Key**
   - Visit [DeepSeek](https://platform.deepseek.com/api_keys)
   - Create and copy API Key

3. **Configure OpenRouter**
   - Go to [OpenRouter Integration](https://openrouter.ai/settings/integrations)
   - Enter your Google AI Studio / DeepSeek API Key
   - Visit [OpenRouter Keys](https://openrouter.ai/settings/keys)
   - Create and get OpenRouter API Key

4. **Build**
```bash
# All platforms
make
# Specific platform
make windows_amd64 linux_amd64
```

## Notes
- Keep your API Keys secure
- Read the terms of service before use

## Contributing
Issues and Pull Requests are welcome!