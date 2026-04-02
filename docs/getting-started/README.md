# Getting Started

This guide walks you through setting up the project locally for development.

---

## Prerequisites

<!-- List the tools required to run this project and link to their installation pages. -->

| Tool | Version | Install |
|------|---------|---------|
| <!-- e.g. Node --> | <!-- e.g. >= 22.x --> | <!-- link --> |
| <!-- e.g. Docker --> | <!-- e.g. >= 27.x --> | <!-- link --> |

---

## Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/your-org/your-repo.git
   cd your-repo
   ```

2. **Activate the commit-msg hook** (one-time, after cloning)

   ```bash
   git config core.hooksPath .githooks
   ```

3. **Copy the environment file**

   ```bash
   cp .env.example .env
   ```

   Open `.env` and fill in the required values. See [Environment Variables](#environment-variables) below.

4. **Install dependencies**

   ```bash
   # Replace with your package manager / build tool
   # e.g. npm install / pip install -r requirements.txt / go mod download / mvn install
   ```

5. **Start the development server**

   ```bash
   # Replace with your start command
   # e.g. npm run dev / python manage.py runserver / go run ./cmd/server
   ```

---

## Environment Variables

<!-- Document every variable in .env.example here. -->

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT`   | No       | `3000`  | Port the server listens on |
| <!-- add rows --> | | | |

---

## Common Commands

<!-- Replace with your actual commands. -->

```bash
# Start development server
# Run tests
# Run linter
# Build for production
```

---

## Troubleshooting

**Port already in use**

Find and stop the process occupying the port before starting the server.

**Dependencies not installing**

Make sure your runtime version matches the one listed in [Prerequisites](#prerequisites). Delete any lock files and retry a fresh install.

For other issues, check the [FAQ](../developer-guide/faq.md) or open a [Discussion](../../../discussions).
