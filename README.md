# SpotiFLAC

[![Build Status](https://github.com/Paidguy/SpotiFLAC-web/workflows/Build%20and%20Test/badge.svg)](https://github.com/Paidguy/SpotiFLAC-web/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Self-hosted web application to download Spotify tracks in high-quality FLAC audio from Tidal, Qobuz, and Amazon Music — no account required.**

Stream Spotify for discovery, download in lossless quality for your library.

![SpotiFLAC Interface](https://github.com/user-attachments/assets/adbdc056-bace-44a9-8ba6-898b4526b65a)

---

## Features

- **Multiple streaming sources**: Automatically fetch tracks from Tidal, Qobuz, or Amazon Music
- **Supported input types**: Single tracks, albums, and playlists via Spotify URLs
- **High-quality audio**: FLAC (lossless), up to 24-bit/192kHz depending on source
- **Rich metadata**: Embedded cover art, lyrics, track info, album info
- **Real-time progress**: Server-Sent Events (SSE) for live download progress updates
- **Flexible organization**: Customizable folder structure and filename templates
- **Download queue**: Track progress, retry failures, skip duplicates
- **Download history**: Browse and search past downloads
- **Dark/light theme**: Automatic system theme detection
- **Self-hosted**: Full control over your music library, no external services

---

## Quick Start with Docker Compose (Recommended)

The easiest way to run SpotiFLAC is with Docker Compose:

### 1. Create `docker-compose.yml`

```yaml
version: "3.9"

services:
  spotiflac:
    image: ghcr.io/paidguy/spotiflac-web:latest
    # Or build from source:
    # build: .
    container_name: spotiflac
    ports:
      - "8080:8080"
    volumes:
      - ./downloads:/downloads
      - ./data:/app/data
    environment:
      - PORT=8080
      - DOWNLOAD_PATH=/downloads
      - DATA_DIR=/app/data
    restart: unless-stopped
```

### 2. Start the server

```bash
docker compose up -d
```

### 3. Open your browser

Navigate to **http://localhost:8080**

That's it! Paste a Spotify URL and start downloading.

---

## Manual Installation / Build from Source

### Prerequisites

- **Go** 1.22 or later
- **Node.js** 20 or later
- **pnpm** package manager
- **ffmpeg** (required for audio conversion and metadata embedding)

### Installation Steps

#### 1. Install dependencies

**macOS (Homebrew)**:
```bash
brew install go node pnpm ffmpeg
```

**Ubuntu/Debian**:
```bash
sudo apt update
sudo apt install golang-go nodejs npm ffmpeg
sudo npm install -g pnpm
```

**Windows**:
- Download Go: https://go.dev/dl/
- Download Node.js: https://nodejs.org/
- Install pnpm: `npm install -g pnpm`
- Download ffmpeg: https://ffmpeg.org/download.html (add to PATH)

#### 2. Clone the repository

```bash
git clone https://github.com/Paidguy/SpotiFLAC-web.git
cd SpotiFLAC-web
```

#### 3. Build the frontend

```bash
cd frontend
pnpm install
pnpm run build
cd ..
```

#### 4. Build the backend

```bash
go build -o spotiflac .
```

#### 5. Run the server

```bash
DOWNLOAD_PATH=/path/to/music PORT=8080 ./spotiflac
```

Open **http://localhost:8080** in your browser.

---

## Configuration

### Environment Variables

All configuration is done via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP port the server listens on | `8080` |
| `DOWNLOAD_PATH` | Absolute path where downloaded music is saved | `./downloads` |
| `DATA_DIR` | Directory for settings and history database | `./data` |
| `ENV` | Set to `development` for CORS and debug mode | (production) |

**Example**:
```bash
PORT=3000 DOWNLOAD_PATH=/mnt/music DATA_DIR=/var/lib/spotiflac ./spotiflac
```

### Application Settings

Additional settings (downloader service, audio quality, folder structure, filename format, etc.) are configured through the **Settings** panel in the web UI. These settings are persisted to `$DATA_DIR/settings.json`.

**Note**: The download path is set via the `DOWNLOAD_PATH` environment variable and cannot be changed from the UI for security reasons.

---

## Reverse Proxy Setup (nginx)

If running SpotiFLAC behind a reverse proxy, you **must** configure it to support Server-Sent Events (SSE) for real-time progress updates.

### nginx Configuration

```nginx
server {
    listen 80;
    server_name spotiflac.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Required for Server-Sent Events (real-time progress)
        proxy_set_header Connection '';
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 86400s;
        chunked_transfer_encoding on;
    }
}
```

**Key settings for SSE**:
- `proxy_buffering off` - Disables response buffering
- `proxy_cache off` - Disables response caching
- `proxy_set_header Connection ''` - Allows keep-alive connections
- `chunked_transfer_encoding on` - Enables streaming responses

### Caddy Configuration

```caddy
spotiflac.example.com {
    reverse_proxy localhost:8080
}
```

Caddy automatically handles SSE correctly.

---

## Troubleshooting

### ffmpeg not found

**Problem**: Downloads fail with "ffmpeg not found" error.

**Solution**: Install ffmpeg and ensure it's in your system PATH:

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html and add to PATH
```

### Port already in use

**Problem**: Server fails to start with "address already in use" error.

**Solution**: Change the port using the `PORT` environment variable:

```bash
PORT=3000 ./spotiflac
```

Or check what's using port 8080:

```bash
# Linux/macOS
lsof -i :8080

# Windows
netstat -ano | findstr :8080
```

### No files appearing in download path

**Problem**: Downloads complete but no files appear.

**Solution**:
1. Check that `DOWNLOAD_PATH` is set correctly and the directory exists
2. Verify the application has write permissions to the download directory
3. For Docker: ensure the volume mount is correct in `docker-compose.yml`

### Progress bar not updating

**Problem**: Download starts but progress stays at 0%.

**Solution**:
1. Check if a reverse proxy is buffering SSE events (see nginx config above)
2. Verify the browser console for EventSource connection errors
3. Try accessing the server directly (without proxy) to isolate the issue

### Download fails with service error

**Problem**: "Failed to download from Tidal/Qobuz/Amazon" error.

**Solution**:
1. Check that the track is available on the selected service
2. Try a different service (Auto mode tries multiple services)
3. Some regions may have limited availability - try using a VPN
4. Check if the service API is experiencing issues

### Settings not persisting

**Problem**: Settings reset after server restart.

**Solution**:
1. Ensure `DATA_DIR` is set correctly and writable
2. For Docker: verify the data volume mount in `docker-compose.yml`
3. Check server logs for file write errors

---

## API Endpoints

SpotiFLAC exposes a REST API on `/api/*`:

- `GET /api/health` - Health check
- `POST /api/metadata` - Get track/album/playlist metadata from Spotify
- `POST /api/download` - Download a track
- `GET /api/download-queue` - Get current download queue status
- `GET /api/events` - Server-Sent Events for real-time progress
- `GET /api/settings` - Load settings
- `POST /api/settings` - Save settings
- `GET /api/history` - Get download history

See `server/handlers.go` for the complete API reference.

---

## Development

### Running in development mode

```bash
# Terminal 1: Start frontend dev server
cd frontend
pnpm install
pnpm run dev

# Terminal 2: Start backend
ENV=development DOWNLOAD_PATH=./test-downloads go run .
```

The frontend dev server runs on `http://localhost:5173` with hot reload.

### Project Structure

```
spotiflac/
├── backend/           # Go backend logic (downloaders, metadata, etc.)
├── frontend/          # React frontend (Vite + TypeScript)
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── hooks/       # React hooks
│   │   ├── lib/         # Utilities and API client
│   │   └── types/       # TypeScript types
│   └── dist/          # Built frontend (embedded in Go binary)
├── server/            # HTTP handlers and SSE broker
├── main.go            # Server entrypoint
├── Dockerfile         # Docker build configuration
└── docker-compose.yml # Docker Compose setup
```

---

## Disclaimer

This project is for **educational and private use only**. The developers do not condone or encourage copyright infringement.

**SpotiFLAC** is a third-party tool and is not affiliated with, endorsed by, or connected to Spotify, Tidal, Qobuz, Amazon Music, or any other streaming service.

You are solely responsible for:

1. Ensuring your use of this software complies with your local laws
2. Reading and adhering to the Terms of Service of the respective platforms
3. Any legal consequences resulting from the misuse of this tool

The software is provided "as is", without warranty of any kind. The authors assume no liability for any damages, bans, or legal issues arising from its use.

---

## API Credits

- **Tidal**: [hifi-api](https://github.com/binimum/hifi-api)
- **Qobuz**: [dabmusic.xyz](https://dabmusic.xyz), [squid.wtf](https://squid.wtf), [jumo-dl](https://jumo-dl.pages.dev/)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Support

If you find this project useful, consider supporting the original author:

[![Ko-fi](https://img.shields.io/badge/Support%20on%20Ko--fi-72a5f2?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/afkarxyz)

---

## Star History

If this project helps you, please consider giving it a star! ⭐

You'll receive notifications for all new releases.
