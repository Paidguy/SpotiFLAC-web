<div align="center">

# 🎵 SpotiFLAC WEB

[![Build Status](https://github.com/Paidguy/SpotiFLAC-web/workflows/Build%20and%20Test/badge.svg)](https://github.com/Paidguy/SpotiFLAC-web/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/Paidguy/SpotiFLAC-web?style=social)](https://github.com/Paidguy/SpotiFLAC-web/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/Paidguy/SpotiFLAC-web?style=social)](https://github.com/Paidguy/SpotiFLAC-web/network/members)

**Self-hosted web application to download Spotify tracks in pristine FLAC audio quality**

*Stream Spotify for discovery, download in lossless quality for your library*

[Features](#-features) • [Quick Start](#-quick-start) • [Installation](#-installation) • [Configuration](#-configuration) • [Documentation](#-documentation) • [Support](#-support)

![SpotiFLAC Interface](https://github.com/user-attachments/assets/adbdc056-bace-44a9-8ba6-898b4526b65a)

</div>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Quick Start](#-quick-start)
- [Installation](#-installation)
  - [Docker Compose (Recommended)](#docker-compose-recommended)
  - [Manual Installation](#manual-installation)
- [Configuration](#-configuration)
- [Usage Guide](#-usage-guide)
- [Architecture](#-architecture)
- [API Documentation](#-api-documentation)
- [Reverse Proxy Setup](#-reverse-proxy-setup)
- [Troubleshooting](#-troubleshooting)
- [Development](#-development)
- [Contributing](#-contributing)
- [FAQ](#-faq)
- [Security & Privacy](#-security--privacy)
- [Credits](#-credits)
- [License](#-license)

---

## 🌟 Overview

**SpotiFLAC** is a powerful, self-hosted web application that allows you to download Spotify tracks in high-quality FLAC audio format by automatically sourcing them from Tidal, Qobuz, and Amazon Music. No account required on any of these services!

### Why SpotiFLAC?

- 🎧 **Lossless Quality**: Download tracks in FLAC format, up to 24-bit/192kHz
- 🔒 **Privacy First**: Self-hosted, no data collection, complete control
- 🚀 **Fast & Efficient**: Concurrent downloads, smart queue management
- 🎨 **Modern UI**: Beautiful, responsive interface with dark/light themes
- 🔄 **Real-time Progress**: Live download status with Server-Sent Events
- 📦 **Batch Processing**: Download entire albums and playlists
- 🎵 **Rich Metadata**: Embedded cover art, lyrics, and complete track information

---

## ✨ Features

### Core Functionality

- **🎼 Multiple Audio Sources**
  - Automatically fetches tracks from Tidal, Qobuz, or Amazon Music
  - Fallback mechanism ensures maximum success rate
  - No authentication required on source services

- **📥 Flexible Input Types**
  - Single tracks via Spotify URLs
  - Complete albums with automatic track numbering
  - Full playlists with batch download support
  - Artist discographies

- **🎵 High-Quality Audio**
  - FLAC (lossless) format support
  - Up to 24-bit/192kHz quality depending on source
  - Automatic audio format conversion
  - Bitrate selection for lossy formats

### Advanced Features

- **📊 Real-time Progress Tracking**
  - Live download speed monitoring
  - Server-Sent Events (SSE) for instant updates
  - Visual progress indicators
  - Download queue management

- **🎨 Rich Metadata Embedding**
  - High-resolution cover art
  - Synchronized lyrics (when available)
  - Complete track information (artist, album, year, etc.)
  - Disc and track numbering
  - Album artist and compilation tags

- **📁 Flexible Organization**
  - Customizable folder structure with templates
  - Configurable filename patterns
  - Support for variables: `{artist}`, `{album}`, `{title}`, `{track}`, etc.
  - Automatic file sanitization

- **🔄 Download Management**
  - Queue system with priority support
  - Automatic duplicate detection
  - Retry mechanism for failed downloads
  - Resume capability for interrupted downloads
  - Batch operations (download all, download selected)

- **📜 Download History**
  - Browse past downloads
  - Search and filter capabilities
  - Re-download functionality
  - Export history

- **🎨 User Interface**
  - Modern, responsive design
  - Dark/light theme with automatic system detection
  - Drag-and-drop support
  - Keyboard shortcuts
  - Mobile-friendly

- **🔧 Customization**
  - Extensive settings panel
  - Service preferences
  - Audio quality selection
  - Folder and filename templates
  - Download behavior options

---

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose (recommended), OR
- Go 1.22+, Node.js 20+, pnpm, and ffmpeg

### Docker Compose (Recommended)

The fastest way to get SpotiFLAC running:

1. **Create `docker-compose.yml`**

```yaml
version: "3.9"

services:
  spotiflac:
    image: ghcr.io/paidguy/spotiflac-web:latest
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

2. **Start the application**

```bash
docker compose up -d
```

3. **Access the interface**

Open your browser and navigate to **http://localhost:8080**

4. **Start downloading!**

Simply paste a Spotify URL (track, album, or playlist) and click download.

---

## 💻 Installation

### Docker Compose (Recommended)

#### Step 1: Install Docker

- **macOS**: [Docker Desktop for Mac](https://docs.docker.com/desktop/install/mac-install/)
- **Windows**: [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/)
- **Linux**: [Docker Engine](https://docs.docker.com/engine/install/)

#### Step 2: Create Configuration

Create a `docker-compose.yml` file in your preferred directory:

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
      - ./downloads:/downloads    # Your music will be saved here
      - ./data:/app/data         # App data and settings
    environment:
      - PORT=8080
      - DOWNLOAD_PATH=/downloads
      - DATA_DIR=/app/data
    restart: unless-stopped
```

#### Step 3: Launch

```bash
# Start in background
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

#### Step 4: Access

Navigate to `http://localhost:8080` in your web browser.

---

### Manual Installation

For advanced users who want to build from source.

#### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.22+ | Backend compilation |
| **Node.js** | 20+ | Frontend build |
| **pnpm** | Latest | Package management |
| **ffmpeg** | Latest | Audio processing |

#### Platform-Specific Setup

<details>
<summary><b>macOS (Homebrew)</b></summary>

```bash
# Install all dependencies
brew install go node pnpm ffmpeg

# Verify installations
go version
node --version
pnpm --version
ffmpeg -version
```
</details>

<details>
<summary><b>Ubuntu/Debian</b></summary>

```bash
# Update package lists
sudo apt update

# Install dependencies
sudo apt install -y golang-go nodejs npm ffmpeg

# Install pnpm globally
sudo npm install -g pnpm

# Verify installations
go version
node --version
pnpm --version
ffmpeg -version
```
</details>

<details>
<summary><b>Windows</b></summary>

1. **Install Go**: Download from [go.dev/dl](https://go.dev/dl/) and run installer
2. **Install Node.js**: Download from [nodejs.org](https://nodejs.org/) and run installer
3. **Install pnpm**: Open PowerShell and run:
   ```powershell
   npm install -g pnpm
   ```
4. **Install ffmpeg**:
   - Download from [ffmpeg.org/download.html](https://ffmpeg.org/download.html)
   - Extract and add to PATH
   - Verify with: `ffmpeg -version`
</details>

#### Build Steps

1. **Clone the repository**

```bash
git clone https://github.com/Paidguy/SpotiFLAC-web.git
cd SpotiFLAC-web
```

2. **Build the frontend**

```bash
cd frontend
pnpm install
pnpm run build
cd ..
```

3. **Build the backend**

```bash
go build -o spotiflac .
```

4. **Run the application**

```bash
# Basic usage
./spotiflac

# Custom configuration
DOWNLOAD_PATH=/path/to/music PORT=8080 ./spotiflac
```

5. **Access the interface**

Open `http://localhost:8080` in your browser.

---

## ⚙️ Configuration

### Environment Variables

SpotiFLAC is configured entirely through environment variables:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | HTTP port for the web server | `8080` | `3000` |
| `DOWNLOAD_PATH` | Absolute path for downloaded music | `./downloads` | `/mnt/music` |
| `DATA_DIR` | Directory for settings and database | `./data` | `/var/lib/spotiflac` |
| `ENV` | Environment mode (production/development) | `production` | `development` |

#### Configuration Examples

**Basic Usage:**
```bash
PORT=8080 DOWNLOAD_PATH=./music ./spotiflac
```

**Custom Paths:**
```bash
PORT=3000 \
DOWNLOAD_PATH=/home/user/Music \
DATA_DIR=/home/user/.spotiflac \
./spotiflac
```

**Development Mode:**
```bash
ENV=development \
DOWNLOAD_PATH=./test-downloads \
go run .
```

### Application Settings

Additional settings are managed through the web UI Settings panel:

- **Download Service**: Choose between Tidal, Qobuz, Amazon Music, or Auto
- **Audio Quality**: Select preferred quality level
- **Folder Structure**: Customize directory organization
- **Filename Format**: Define filename patterns with variables
- **Download Behavior**: Configure retry attempts, timeout values, etc.

Settings are persisted to `$DATA_DIR/settings.json` and persist across restarts.

**Note**: The `DOWNLOAD_PATH` cannot be changed from the UI for security reasons.

---

## 📖 Usage Guide

### Basic Workflow

1. **Open SpotiFLAC** in your browser
2. **Paste a Spotify URL** into the search bar
   - Track: `https://open.spotify.com/track/...`
   - Album: `https://open.spotify.com/album/...`
   - Playlist: `https://open.spotify.com/playlist/...`
3. **Browse the results** - view tracks, albums, or playlist contents
4. **Select tracks** - choose individual tracks or select all
5. **Click Download** - watch real-time progress as downloads complete
6. **Find your music** in the configured download directory

### Advanced Features

#### Batch Downloads

- **Download All**: Click "Download All" to queue every track
- **Download Selected**: Select specific tracks and click "Download Selected"
- **Download Albums**: When viewing a playlist, download complete albums

#### Lyrics & Cover Art

- **Download Lyrics**: Optional separate lyrics download for each track
- **Download Covers**: Save high-resolution album art separately
- **Batch Operations**: Download all lyrics or covers at once

#### Search Mode

- Switch to **Search Mode** to search Spotify directly
- Browse tracks, albums, artists, and playlists
- Click any result to load full details

#### Download Queue

- View active downloads in the queue dialog
- Monitor progress, speed, and status
- Retry failed downloads
- Cancel pending downloads
- Export failed downloads for later

#### History

- Browse complete download history
- Search by track, artist, or album
- Re-download previous tracks
- Clear history selectively or entirely

---

## 🏗️ Architecture

### Technology Stack

**Backend:**
- **Language**: Go 1.22+
- **Framework**: Echo v4 (HTTP server)
- **Database**: bbolt (embedded key-value store)
- **Audio Processing**: ffmpeg (via CLI)

**Frontend:**
- **Framework**: React 19
- **Build Tool**: Vite 7
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS 4
- **UI Components**: Radix UI
- **State Management**: React Hooks

**Infrastructure:**
- **Container**: Docker + Docker Compose
- **Deployment**: Multi-stage Docker build
- **Base Images**: golang:alpine, node:alpine

### System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Web Browser                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │         React Frontend (SPA)                     │  │
│  │  • UI Components  • State Management             │  │
│  │  • API Client     • Real-time Updates (SSE)      │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────┬───────────────────────────────────┘
                      │ HTTP/SSE
┌─────────────────────▼───────────────────────────────────┐
│              Go Backend Server                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  HTTP Server (Echo)                              │  │
│  │  • REST API Endpoints                            │  │
│  │  • SSE Event Broker                              │  │
│  │  • Static File Serving                           │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Core Services                                    │  │
│  │  • Metadata Service (Spotify API)                │  │
│  │  • Download Service (Tidal/Qobuz/Amazon)         │  │
│  │  • Queue Manager                                 │  │
│  │  • Progress Tracker                              │  │
│  └──────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Data Layer                                       │  │
│  │  • bbolt Database (history, cache)               │  │
│  │  • JSON Settings Store                           │  │
│  │  • File System (downloads)                       │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────┬───────────────────────────────────┘
                      │
         ┌────────────┼────────────┐
         │            │            │
    ┌────▼───┐   ┌───▼────┐   ┌──▼─────┐
    │ Tidal  │   │ Qobuz  │   │ Amazon │
    │  API   │   │  API   │   │  Music │
    └────────┘   └────────┘   └────────┘
```

### Project Structure

```
SpotiFLAC-web/
├── backend/               # Core Go business logic
│   ├── amazon.go         # Amazon Music integration
│   ├── qobuz.go          # Qobuz integration
│   ├── tidal.go          # Tidal integration
│   ├── metadata.go       # Spotify metadata fetching
│   ├── progress.go       # Download queue & progress
│   ├── filename.go       # File naming logic
│   ├── lyrics.go         # Lyrics fetching
│   └── cover.go          # Cover art handling
├── server/               # HTTP server layer
│   ├── handlers.go       # API endpoint handlers
│   ├── sse.go           # Server-Sent Events broker
│   └── types.go         # Request/response types
├── frontend/            # React application
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── hooks/       # Custom React hooks
│   │   ├── lib/         # Utilities and API client
│   │   └── types/       # TypeScript type definitions
│   └── dist/           # Build output (embedded in binary)
├── docs/               # Documentation
│   ├── README.md       # Documentation index
│   └── development/    # Development docs
├── main.go            # Application entry point
├── Dockerfile         # Container build definition
└── docker-compose.yml # Orchestration config
```

---

## 📡 API Documentation

### REST Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/health` | Health check |
| `POST` | `/api/metadata` | Fetch Spotify metadata |
| `POST` | `/api/download` | Queue a track download |
| `GET` | `/api/download-queue` | Get queue status |
| `GET` | `/api/events` | SSE stream for real-time updates |
| `GET` | `/api/settings` | Load application settings |
| `POST` | `/api/settings` | Save application settings |
| `GET` | `/api/history` | Get download history |
| `DELETE` | `/api/history` | Clear download history |
| `POST` | `/api/lyrics` | Download lyrics file |
| `POST` | `/api/cover` | Download cover art |
| `POST` | `/api/search` | Search Spotify |

### API Examples

<details>
<summary><b>Fetch Track Metadata</b></summary>

```bash
curl -X POST http://localhost:8080/api/metadata \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://open.spotify.com/track/...",
    "batch": true,
    "delay": 1.0,
    "timeout": 300
  }'
```
</details>

<details>
<summary><b>Download a Track</b></summary>

```bash
curl -X POST http://localhost:8080/api/download \
  -H "Content-Type: application/json" \
  -d '{
    "service": "auto",
    "track_name": "Song Title",
    "artist_name": "Artist Name",
    "album_name": "Album Name",
    "spotify_id": "..."
  }'
```
</details>

<details>
<summary><b>Get Download Queue</b></summary>

```bash
curl http://localhost:8080/api/download-queue
```
</details>

For complete API documentation, see [`server/handlers.go`](server/handlers.go).

---

## 🔐 Reverse Proxy Setup

For production deployments behind a reverse proxy.

### nginx Configuration

SpotiFLAC uses Server-Sent Events (SSE) for real-time progress updates. Your nginx configuration **must** support SSE:

```nginx
server {
    listen 80;
    server_name spotiflac.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;

        # Standard proxy headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Required for Server-Sent Events (SSE)
        proxy_set_header Connection '';
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 86400s;
        chunked_transfer_encoding on;
    }
}
```

**SSL/HTTPS Configuration:**

```nginx
server {
    listen 443 ssl http2;
    server_name spotiflac.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        # Same proxy configuration as above
        proxy_pass http://localhost:8080;
        # ... (rest of config)
    }
}
```

### Caddy Configuration

Caddy automatically handles SSE correctly:

```caddy
spotiflac.example.com {
    reverse_proxy localhost:8080
}
```

### Apache Configuration

```apache
<VirtualHost *:80>
    ServerName spotiflac.example.com

    ProxyPreserveHost On
    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/

    # SSE Support
    ProxyPass / http://localhost:8080/ disablereuse=On
</VirtualHost>
```

---

## 🔧 Troubleshooting

### Common Issues

<details>
<summary><b>ffmpeg not found</b></summary>

**Symptom**: Downloads fail with "ffmpeg not found" error

**Solution**: Install ffmpeg and ensure it's in your PATH

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
# Download from https://ffmpeg.org/download.html
# Add to PATH environment variable
```

Verify installation:
```bash
ffmpeg -version
```
</details>

<details>
<summary><b>Port already in use</b></summary>

**Symptom**: "address already in use" error on startup

**Solution**: Change the port or kill the conflicting process

```bash
# Check what's using port 8080
lsof -i :8080  # Linux/macOS
netstat -ano | findstr :8080  # Windows

# Use a different port
PORT=3000 ./spotiflac
```
</details>

<details>
<summary><b>No files appearing in download directory</b></summary>

**Symptom**: Downloads complete but files don't appear

**Solution**:
1. Verify `DOWNLOAD_PATH` is correct and exists
2. Check write permissions on the download directory
3. For Docker: verify volume mount in `docker-compose.yml`
4. Check server logs for errors

```bash
# Check permissions
ls -la /path/to/downloads

# Docker volume inspection
docker compose exec spotiflac ls -la /downloads
```
</details>

<details>
<summary><b>Progress bar not updating</b></summary>

**Symptom**: Download starts but progress stays at 0%

**Solution**:
1. Check reverse proxy configuration (see [Reverse Proxy Setup](#-reverse-proxy-setup))
2. Verify browser console for EventSource errors
3. Test direct access (bypass proxy)
4. Ensure no firewall blocking SSE connections
</details>

<details>
<summary><b>Download fails with service error</b></summary>

**Symptom**: "Failed to download from Tidal/Qobuz/Amazon" error

**Solution**:
1. Verify track is available on the selected service
2. Try a different service or "Auto" mode
3. Check internet connectivity
4. Some regions have limited availability - try VPN
5. Service API may be temporarily unavailable
</details>

<details>
<summary><b>Settings not persisting</b></summary>

**Symptom**: Settings reset after restart

**Solution**:
1. Ensure `DATA_DIR` is writable
2. For Docker: verify data volume mount
3. Check server logs for write errors
4. Ensure correct permissions on settings file

```bash
# Check settings file
ls -la $DATA_DIR/settings.json

# Docker data volume check
docker compose exec spotiflac ls -la /app/data
```
</details>

### Debug Mode

Enable debug logging for troubleshooting:

```bash
ENV=development ./spotiflac
```

Check logs:
```bash
# Docker
docker compose logs -f

# Direct execution
# Logs output to stdout
```

---

## 🛠️ Development

### Local Development Setup

1. **Clone and install dependencies**

```bash
git clone https://github.com/Paidguy/SpotiFLAC-web.git
cd SpotiFLAC-web
cd frontend && pnpm install && cd ..
go mod download
```

2. **Start development servers**

```bash
# Terminal 1: Frontend dev server (with hot reload)
cd frontend
pnpm run dev
# Runs on http://localhost:5173

# Terminal 2: Backend server
ENV=development DOWNLOAD_PATH=./test-downloads go run .
# Runs on http://localhost:8080
```

3. **Development workflow**

- Frontend changes auto-reload at `http://localhost:5173`
- Backend changes require restart (`Ctrl+C` and rerun)
- Frontend proxies API requests to backend

### Building from Source

```bash
# Build frontend
cd frontend
pnpm install
pnpm run build
cd ..

# Build backend (embeds frontend)
go build -o spotiflac .

# Run
./spotiflac
```

### Running Tests

```bash
# Frontend tests
cd frontend
pnpm test

# Go tests
go test ./...
```

### Code Quality

```bash
# Frontend linting
cd frontend
pnpm lint

# Go formatting
go fmt ./...

# Go linting (requires golangci-lint)
golangci-lint run
```

---

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines.

### How to Contribute

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes**
4. **Test thoroughly**
5. **Commit your changes**: `git commit -m 'Add amazing feature'`
6. **Push to the branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Contribution Guidelines

- Follow existing code style and conventions
- Write clear commit messages
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting PR

### Development Resources

- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [Echo Framework](https://echo.labstack.com/)
- [Vite Documentation](https://vitejs.dev/)

---

## ❓ FAQ

<details>
<summary><b>Is this legal?</b></summary>

SpotiFLAC is provided for **educational and private use only**. Users are responsible for ensuring their use complies with local laws and service terms. The tool itself doesn't circumvent DRM or access premium content - it sources publicly available audio streams.
</details>

<details>
<summary><b>Do I need accounts on Tidal, Qobuz, or Amazon Music?</b></summary>

No! SpotiFLAC uses public APIs and doesn't require authentication on any source service.
</details>

<details>
<summary><b>What audio quality can I expect?</b></summary>

Quality depends on the source service:
- **Tidal**: Up to 24-bit/96kHz FLAC
- **Qobuz**: Up to 24-bit/192kHz FLAC
- **Amazon Music**: Up to 24-bit/48kHz FLAC

Actual quality varies by track availability.
</details>

<details>
<summary><b>Can I download private playlists?</b></summary>

No, only public Spotify playlists, albums, and tracks can be downloaded. Private playlists are not accessible.
</details>

<details>
<summary><b>Why do some downloads fail?</b></summary>

Downloads may fail if:
- Track is not available on any source service
- Regional restrictions apply
- Source service API is temporarily unavailable
- Network connectivity issues

Try using "Auto" mode to attempt all services, or try a different region with a VPN.
</details>

<details>
<summary><b>Can I run this on a Raspberry Pi?</b></summary>

Yes! SpotiFLAC runs well on Raspberry Pi 4 or newer. Use the Docker image for easiest setup. ARM64 builds are available.
</details>

<details>
<summary><b>How much disk space do I need?</b></summary>

FLAC files are typically 20-40MB per song. A 100-song playlist requires ~3GB. Plan accordingly for large libraries.
</details>

<details>
<summary><b>Can I contribute to the project?</b></summary>

Absolutely! See the [Contributing](#-contributing) section for guidelines.
</details>

---

## 🔒 Security & Privacy

### Privacy Commitment

- **No telemetry**: SpotiFLAC doesn't collect or transmit user data
- **No analytics**: No tracking, no metrics, no phone-home features
- **Local-only**: All data stays on your server
- **Self-hosted**: You control everything

### Security Considerations

- **Path validation**: All file paths are validated to prevent traversal attacks
- **Input sanitization**: User inputs are sanitized to prevent injection
- **CORS configured**: Proper Cross-Origin Resource Sharing for security
- **No authentication**: For simplicity, use reverse proxy auth if needed

### Recommendations

For production deployments:

1. **Use HTTPS**: Set up SSL/TLS certificates
2. **Firewall**: Restrict access to trusted networks
3. **Reverse proxy auth**: Add authentication layer (Basic Auth, OAuth, etc.)
4. **Regular updates**: Keep dependencies and Docker images updated
5. **Backup data**: Regularly backup your `DATA_DIR`

---

## 🎨 Credits

### Project Maintainer

**[@Paidguy](https://github.com/Paidguy)** - Current maintainer and web version developer

This web version includes:
- Complete web-based interface migration
- Enhanced UI/UX with modern React components
- Docker containerization support
- Comprehensive documentation
- Bug fixes and stability improvements
- New features and functionality

### Original Creator

**[@afkarxyz](https://github.com/afkarxyz)** - Original SpotiFLAC concept and implementation

This project builds upon the excellent foundation created by afkarxyz.

### Special Thanks

- **API Credits**:
  - [hifi-api](https://github.com/binimum/hifi-api) - Tidal integration
  - [dabmusic.xyz](https://dabmusic.xyz) - Qobuz integration
  - [squid.wtf](https://squid.wtf) - Additional Qobuz support
  - [jumo-dl](https://jumo-dl.pages.dev/) - Qobuz integration

- **Open Source Libraries**:
  - [Echo](https://github.com/labstack/echo) - Go web framework
  - [React](https://github.com/facebook/react) - UI framework
  - [Vite](https://github.com/vitejs/vite) - Build tool
  - [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss) - Styling
  - [Radix UI](https://github.com/radix-ui/primitives) - UI components
  - [ffmpeg](https://github.com/FFmpeg/FFmpeg) - Audio processing

---

## 📄 License

This project is dual-licensed:

**Primary License**: MIT License

Copyright (c) 2026 [@Paidguy](https://github.com/Paidguy) (Web version maintainer)

Original work Copyright (c) 2026 [@afkarxyz](https://github.com/afkarxyz)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

See the [LICENSE](LICENSE) file for complete details.

---

## ⚠️ Disclaimer

**Important Legal Notice**

This software is provided for **educational and private use only**. The developers do not condone or encourage copyright infringement.

**SpotiFLAC** is a third-party tool and is **not affiliated with, endorsed by, or connected to** Spotify, Tidal, Qobuz, Amazon Music, or any other streaming service.

### Your Responsibilities

By using this software, you acknowledge and agree that:

1. You are solely responsible for ensuring your use complies with local laws and regulations
2. You must read and adhere to the Terms of Service of all respective platforms
3. You accept full responsibility for any legal consequences resulting from misuse
4. The software is provided "as is" without warranty of any kind
5. The authors assume no liability for damages, bans, or legal issues arising from use

### Ethical Use

Please:
- ✅ Support artists by purchasing music and concert tickets
- ✅ Use streaming services for discovery and support
- ✅ Respect copyright and intellectual property rights
- ✅ Use this tool responsibly and legally

---

## 💝 Support

### Support the Maintainer

If you find this web version useful, consider supporting [@Paidguy](https://github.com/Paidguy):

[![GitHub Sponsor](https://img.shields.io/badge/Sponsor-%E2%9D%A4-red?style=for-the-badge&logo=github)](https://github.com/sponsors/Paidguy)

### Support the Original Creator

Show appreciation to the original creator [@afkarxyz](https://github.com/afkarxyz):

[![Ko-fi](https://img.shields.io/badge/Support%20on%20Ko--fi-72a5f2?style=for-the-badge&logo=ko-fi&logoColor=white)](https://ko-fi.com/afkarxyz)

---

## ⭐ Star History

If this project helps you, please consider giving it a star! ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=Paidguy/SpotiFLAC-web&type=Date)](https://star-history.com/#Paidguy/SpotiFLAC-web&Date)

---

## 🔗 Links

- **GitHub Repository**: [Paidguy/SpotiFLAC-web](https://github.com/Paidguy/SpotiFLAC-web)
- **Issues**: [Report a bug or request a feature](https://github.com/Paidguy/SpotiFLAC-web/issues)
- **Discussions**: [Join the community](https://github.com/Paidguy/SpotiFLAC-web/discussions)
- **Docker Hub**: [SpotiFLAC Images](https://github.com/Paidguy/SpotiFLAC-web/pkgs/container/spotiflac-web)

---

<div align="center">

**Made with ❤️ by [@Paidguy](https://github.com/Paidguy)**

**Built upon the original work by [@afkarxyz](https://github.com/afkarxyz)**

[⬆ Back to Top](#-spotiflac)

</div>
