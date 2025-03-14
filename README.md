GitCommitAI is an AI-powered command-line tool written in Go that generates meaningful Git commit messages by analyzing your latest changes (git diff). It supports both OpenAI’s GPT API and locally hosted models, adheres to the Conventional Commits format, and offers an interactive workflow to review or edit suggestions before committing.
Features
AI-Driven: Generates context-aware commit messages using OpenAI’s GPT or a local model via FastAPI.

Conventional Commits: Suggests messages like feat: add login endpoint or fix: resolve bug in parser.

Interactive: Review, edit, or accept AI-generated messages with a smooth CLI experience.

Fast & Portable: Built in Go for performance and single-binary distribution.

Configurable: Supports YAML/JSON/TOML configuration for API keys, model selection, and more.

Extensible: Modular design with robust error handling and logging.

Installation
Prerequisites
Go 1.21+ (for building from source)

Git (any recent version)

Optional: Python 3.9+ and FastAPI (for local model support)

Via Homebrew (macOS)
bash

brew tap AlexThuku/GitCommitAI-
brew install GitCommitAI-

Via Linux Package
For Debian/Ubuntu:
bash

wget https://github.com/AlexThuku/GitCommitAI-/releases/latest/download/GitCommitAI-.deb
sudo dpkg -i GitCommitAI-.deb

For Fedora/RHEL:
bash

wget https://github.com/AlexThuku/GitCommitAI-/releases/latest/download/GitCommitAI-.rpm
sudo rpm -i GitCommitAI-.rpm

From Source
bash

git clone https://github.com/AlexThuku/GitCommitAI-.git
cd GitCommitAI-
go build -o GitCommitAI- ./cmd/GitCommitAI-
sudo mv GitCommitAI- /usr/local/bin/

Binary Download
Grab the latest pre-built binary from the Releases page and move it to your PATH (e.g., /usr/local/bin/).
Usage
Basic Command
Generate a commit message for your staged changes:
bash

GitCommitAI- generate

Example output:

Analyzing changes...
Suggested commit: "feat: implement user login endpoint"
[a]ccept, [e]dit, [r]eject?

Flags
--model: Choose the AI model (openai or local, default: openai).

--edit: Force edit mode for the suggested message.

--config: Specify a custom config file (default: ./config.yaml).

Example:
bash

GitCommitAI- generate --model local --edit

Git Integration
Stage your changes first:
bash

git add .
GitCommitAI- generate

Configuration
GitCommitAI- uses a configuration file (default: config.yaml) to manage settings. Create one in your project root or home directory:
Sample config.yaml
yaml

model:
  type: "openai"          # or "local"
  openai:
    api_key: "your-openai-key-here"
  local:
    endpoint: "http://localhost:8000/generate"
logging:
  level: "info"           # debug, info, warn, error

OpenAI: Set OPENAI_API_KEY as an environment variable or in the config file.

Local Model: Run a FastAPI server (see Local Model Setup (#local-model-setup)) and point to its endpoint.

Environment Variables
Override config settings with:
bash

export OPENAI_API_KEY="your-key"
export GitCommitAI-_MODEL="local"
GitCommitAI- generate

Local Model Setup
For local inference, deploy a Python FastAPI backend:
Install dependencies:
bash

cd backend
python -m venv venv
source venv/bin/activate
pip install fastapi uvicorn your-ai-model-lib

Start the server:
bash

uvicorn main:app --host 0.0.0.0 --port 8000

Configure GitCommitAI- to use model.type: local and the correct endpoint.

Project Structure

GitCommitAI-/
├── cmd/           # CLI entry point
├── internal/      # Private packages (ai, config, git)
├── pkg/           # Reusable utilities (optional)
├── backend/       # Python FastAPI backend (optional)
├── config.yaml    # Sample config
├── go.mod         # Go module file
└── README.md

Contributing
We welcome contributions! To get started:
Fork the repo.

Create a feature branch (git checkout -b feat/my-feature).

Commit changes following Conventional Commits.

Submit a pull request.

See CONTRIBUTING.md for details.
Development
Run tests:
bash

go test ./...

Build locally:
bash

go build -o GitCommitAI- ./cmd/GitCommitAI-

License
MIT License (LICENSE) - feel free to use, modify, and distribute.
Acknowledgements
Built with Go, OpenAI, and FastAPI.

Inspired by tools like aicommits and OpenCommit.

Contact
Questions? Open an issue or reach out at alexmwangithuku001@gmail.com
