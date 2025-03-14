# AI-Powered Git Commit Message Generator

Welcome to the AI-Powered Git Commit Message Generator! This tool leverages advanced AI models to generate meaningful and conventional commit messages based on your code changes. Whether you're using OpenAI, Hugging Face, or a local model, this tool aims to streamline your commit process and improve your commit history.

## Features

- **AI-Driven Commit Messages**: Automatically generate commit messages using AI models.
- **Multiple AI Providers**: Supports OpenAI, Hugging Face, and local models.
- **Customizable Configuration**: Easily switch between AI providers and configure settings.
- **Interactive Approval**: Review, edit, or reject suggested commit messages before committing.

## Getting Started

### Installation

1. **Clone the Repository**:
   ```bash
   git clone git@github.com:AlexThuku/GitCommitAI-.git
   cd ai-git-commit-generator
   ```

2. **Build the Project**:
   ```bash
   go build -o git-msg cmd/git-msg/main.go
   ```

3. **Set Up Configuration**:
   Create a `git-msg.yaml` file in your home directory or current directory with the following content:
   ```yaml
   model_provider: "huggingface" # or "openai", "local"
   huggingface_token: "your_huggingface_token_here"
   huggingface_model: "mistralai/Mistral-7B-Instruct-v0.2"
   openai_api_key: "your_openai_api_key_here"
   openai_model: "gpt-4o"
   local_endpoint: "http://localhost:8000/generate"
   ```

### Usage

1. **Navigate to Your Git Repository**:
   ```bash
   cd path/to/your/repo
   ```

2. **Stage Your Changes**:
   ```bash
   git add .
   ```

3. **Generate a Commit Message**:
   ```bash
   ../path/to/git-msg generate
   ```

4. **Review and Approve**:
   - Accept, edit, or reject the suggested commit message.

### Environment Variables

Instead of using a configuration file, you can set environment variables:

```bash
export OPENAI_API_KEY="your_openai_api_key_here"
export HUGGINGFACE_TOKEN="your_huggingface_token_here"
```

## Troubleshooting

- **Configuration Errors**: Ensure your `git-msg.yaml` file is correctly formatted and free of control characters.
- **API Errors**: Verify your API keys and model IDs are correct and have the necessary permissions.
- **Model Selection**: Experiment with different models if the output isn't as expected.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.



## Contact

For questions or support, please contact alexmwangithuku001@gmail.com.

---

Thank you for using the AI-Powered Git Commit Message Generator! We hope it enhances your development workflow.