# Hyprlander 🏔️

> **An intelligent ReAct agent for Hyprland configuration management**

Hyprlander is a command-line tool that makes it easier than ever to manage and customize your Hyprland window manager. Just type what you want to change in plain language and Hyprlander figures it out and updates your config for you.

Inspired by tools like `gemini-cli` and `claude-code`, it uses a smart reasoning framework to interpret your requests and take action, so you don’t have to dig through config files or memorize syntax.

## ✨ Features

- **ReAct Agent Architecture**: Uses a reasoning-and-action cycle to understand your intent and apply configuration changes.
- **Conversational Interface**: Describe what you want in plain language — no special syntax or commands required.
- **AI-Powered Understanding**: Integrates with the Gemini API to interpret your requests.
- **Autonomous Actions**: Automatically reads, edits, and updates your Hyprland config files.
- **Context Awareness**: Understands your current Hyprland setup before making any changes.

## 🚀 Quick Start

### Prerequisites

- Linux environment with Hyprland window manager installed
- Go 1.25.1 or later
- Gemini API key (get one from [Google AI Studio](https://aistudio.google.com/))

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/saat-sy/hyprlander.git
   cd hyprlander
   ```

2. **Build the project:**
   ```bash
   go build -o hyprlander .
   ```

3. **Install globally (optional):**
   ```bash
   sudo mv hyprlander /usr/local/bin/
   ```

### Initial Setup

Before using Hyprlander, you need to initialize it with your Gemini API key:

```bash
hyprlander init
```

You'll be prompted to enter your Gemini API key. The key will be securely stored in `~/.cache/.hyprlander/secrets.ini`.

## 🛠️ Usage

### Basic Commands

```bash
# Get help
hyprlander --help

# Initialize Hyprlander
hyprlander init

# Start a conversation with your Hyprland setup
hyprlander prompt "make my desktop more minimalist"
hyprlander prompt "I'm having screen tearing issues"
hyprlander prompt "optimize for gaming performance"
```

### How It Works (ReAct Framework)

1. **Reasoning**: Agent analyzes your request and current Hyprland configuration
2. **Action**: Performs necessary file operations (read config, backup, modify)
3. **Observation**: Validates changes and checks for conflicts
4. **Iteration**: Continues reasoning-action cycles until task is complete

The agent can:
- Read and understand your existing `hyprland.conf`
- Research Hyprland documentation and best practices
- Make incremental changes with validation at each step
- Explain what it's doing and why
- Rollback changes if something goes wrong

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Test thoroughly**
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/hyprlander.git
cd hyprlander

# Install dependencies
go mod tidy

# Run the application
go run main.go

# Build for testing
go build -o hyprlander .
```

## 🛣️ Roadmap

- [x] ✅ Basic CLI structure
- [x] ✅ Initialization system
- [x] ✅ API key storage
- [ ] 🚧 ReAct agent core implementation
- [ ] 🚧 Hyprland configuration parser
- [ ] 🚧 Safe configuration modification with rollback
- [ ] 📋 Configuration versioning and history
- [ ] 📋 Integration with Hyprland community configs and themes
- [ ] 📋 Multi-modal input (screenshots for visual feedback)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🧩 Acknowledgments

- [Hyprland](https://hyprland.org/) - The amazing dynamic tiling Wayland compositor
- [Cobra](https://github.com/spf13/cobra) - Powerful CLI framework for Go
- [Google Gemini](https://ai.google.dev/) - AI model powering the ReAct agent
- [ReAct Paper](https://arxiv.org/abs/2210.03629) - Inspiration for the reasoning and acting framework
- [Gemini CLI](https://github.com/replit/gemini-cli) & [Claude Code](https://claude.ai/) - Similar agent architectures

## 📞 Support

- 🐛 **Bug Reports**: [Open an issue](https://github.com/saat-sy/hyprlander/issues)
- 💡 **Feature Requests**: [Start a discussion](https://github.com/saat-sy/hyprlander/discussions)
