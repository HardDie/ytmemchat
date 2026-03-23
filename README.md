# ytmemchat 📺🤖

**ytmemchat** is a high-performance Go-based bridge that connects **YouTube Live Chat** to your **OBS Overlay**. It transforms your chat into an interactive experience by triggering visual/audio memes (alerts) via custom commands and providing real-time, cross-platform Text-to-Speech (TTS).

---

## ✨ Features

* **YouTube Integration**: Real-time chat polling with auto-throttling and history skipping.
* **Smart Alerts**: Trigger media files (GIFs, Videos, Sounds) using custom tokens (e.g., `@wow`).
* **Cross-Platform TTS**: Native support for Windows (PowerShell), macOS (`say`), and Linux (`espeak`).
* **Websocket Overlay**: A modern, transparent OBS overlay with random positioning, scale/volume control, and auto-reconnect.
* **Developer Webhooks**: Manually inject messages or interrupt the TTS queue via HTTP API.

---

## 🚀 Quick Start

### 1. Prerequisites
* **Go 1.21+** (to build the binary).
* **Google Cloud API Key**: [Create one here](https://console.cloud.google.com/) and enable the **YouTube Data API v3**.
* **System Dependencies**:
    * **Windows**: None (uses built-in PowerShell).
    * **macOS**: None (uses built-in `say`).
    * **Linux**: `sudo apt install espeak`

### 2. Installation
```bash
git clone [https://github.com/HardDie/ytmemchat.git](https://github.com/HardDie/ytmemchat.git)
cd ytmemchat
go build -o ytmemchat ./cmd/main.go
```

### 3. Configuration

Copy the template and fill in your credentials:
```bash
cp env.example .env
```
Key variables:
- `YOUTUBE_API_KEY`: Your Google Cloud API key.
- `YOUTUBE_STREAM_ID`: The unique ID in your YouTube stream URL.
- `ALERTS_TOKEN`: The prefix character for commands (default: @).

### 4. Setup Alert Commands

Create a `commands.yaml` in your project root to map chat words to media files:
```yaml
commands:
  - name: "jump"
    file: "mario_jump.mp3"
    volume: 0.5
  - name: "dance"
    file: "dancing_cat.gif"
    scale: 1.2
```

### 🎥 OBS Integration

1. Add Source: Create a new Browser Source in OBS.
1. URL: Set it to `http://localhost:8080`.
1. Size: Set to your canvas size (e.g., `1920x1080`).
1. Interaction: Click "Interact" on the source and click the "Enable Audio & Connect" button. This is required by browsers to allow autoplaying audio.

### 🛠 Architecture & Workflow

1. Ingestion: The app polls YouTube Chat or listens for Webhook POST requests.
1. Processing:
    - If a message starts with your token (e.g., @wow), it searches for a matching command in commands.yaml.
    - If found, it broadcasts a visual/audio alert to the overlay.
    - If no alert is triggered and TTS is enabled, it synthesizes the message text into audio.
1. Delivery: The internal WebSocket server pushes the payload to the OBS browser source for rendering.

### 📡 Webhook API (For Testing)

Inject a fake chat message:
```bash
curl -X POST http://localhost:8080/webhook/ \
     -d '{"message": "@jump"}'
```

Interrupt current TTS:
```bash
curl -X POST http://localhost:8080/interrupt/
```

### 📂 Project Structure

- `cmd/main.go`: Application entry point and service orchestration.
- `internal/alerts`: Logic for parsing commands.yaml and matching tokens.
- `internal/clients/youtube`: YouTube API iterator and message normalization.
- `internal/server`: The HTTP/WebSocket hub and embedded OBS frontend.
- `internal/tts`: Cross-platform audio synthesis drivers.
- `internal/webhook`: HTTP handlers for external control.