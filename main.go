package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"spotiflac/backend"
	"spotiflac/server"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed all:frontend/dist
var frontendDist embed.FS

func main() {
	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	downloadPath := os.Getenv("DOWNLOAD_PATH")
	if downloadPath == "" {
		downloadPath = "./downloads"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	env := os.Getenv("ENV")
	isDev := env == "development"

	// Create directories if they don't exist
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize history database
	if err := backend.InitHistoryDB("SpotiFLAC"); err != nil {
		log.Printf("Failed to init history DB: %v", err)
	}
	defer backend.CloseHistoryDB()

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware for development
	if isDev {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
	}

	// Create server instance
	srv := server.NewServer(downloadPath, dataDir)

	// API routes
	api := e.Group("/api")

	// Health check
	api.GET("/health", srv.HandleHealth)

	// Metadata and search
	api.POST("/metadata", srv.HandleGetSpotifyMetadata)
	api.GET("/streaming-urls", srv.HandleGetStreamingURLs)
	api.POST("/search", srv.HandleSearchSpotify)
	api.POST("/search-by-type", srv.HandleSearchSpotifyByType)

	// Download operations
	api.POST("/download", srv.HandleDownloadTrack)
	api.POST("/lyrics", srv.HandleDownloadLyrics)
	api.POST("/cover", srv.HandleDownloadCover)
	api.POST("/header", srv.HandleDownloadHeader)
	api.POST("/gallery-image", srv.HandleDownloadGalleryImage)
	api.POST("/avatar", srv.HandleDownloadAvatar)

	// Download queue and progress
	api.GET("/download-progress", srv.HandleGetDownloadProgress)
	api.GET("/download-queue", srv.HandleGetDownloadQueue)
	api.POST("/clear-completed", srv.HandleClearCompletedDownloads)
	api.POST("/clear-all", srv.HandleClearAllDownloads)
	api.POST("/cancel-queued", srv.HandleCancelAllQueuedItems)
	api.POST("/skip-item", srv.HandleSkipDownloadItem)
	api.GET("/export-failed", srv.HandleExportFailedDownloads)

	// Settings
	api.GET("/settings", srv.HandleLoadSettings)
	api.POST("/settings", srv.HandleSaveSettings)
	api.GET("/defaults", srv.HandleGetDefaults)
	api.GET("/download-path", srv.HandleGetDownloadPath)

	// History
	api.GET("/history", srv.HandleGetHistory)
	api.DELETE("/history", srv.HandleDeleteHistory)
	api.DELETE("/history/:id", srv.HandleDeleteHistoryItem)
	api.GET("/fetch-history", srv.HandleGetFetchHistory)
	api.POST("/fetch-history", srv.HandleAddFetchHistory)
	api.DELETE("/fetch-history", srv.HandleClearFetchHistory)
	api.DELETE("/fetch-history/:id", srv.HandleDeleteFetchHistoryItem)
	api.DELETE("/fetch-history/type/:type", srv.HandleClearFetchHistoryByType)

	// Track availability and preview
	api.GET("/track-availability", srv.HandleCheckTrackAvailability)
	api.GET("/preview-url", srv.HandleGetPreviewURL)

	// Audio analysis
	api.GET("/analyze-track", srv.HandleAnalyzeTrack)
	api.POST("/analyze-tracks", srv.HandleAnalyzeMultipleTracks)

	// FFmpeg
	api.GET("/ffmpeg/installed", srv.HandleCheckFFmpegInstalled)
	api.GET("/ffprobe/installed", srv.HandleIsFFprobeInstalled)
	api.GET("/ffmpeg/path", srv.HandleGetFFmpegPath)
	api.POST("/ffmpeg/download", srv.HandleDownloadFFmpeg)

	// Audio conversion
	api.POST("/convert-audio", srv.HandleConvertAudio)

	// File operations
	api.POST("/file-sizes", srv.HandleGetFileSizes)
	api.GET("/list-directory", srv.HandleListDirectoryFiles)
	api.GET("/list-audio-files", srv.HandleListAudioFilesInDir)
	api.GET("/read-metadata", srv.HandleReadFileMetadata)
	api.POST("/preview-rename", srv.HandlePreviewRenameFiles)
	api.POST("/rename-files", srv.HandleRenameFilesByMetadata)
	api.GET("/read-text-file", srv.HandleReadTextFile)
	api.POST("/rename-file", srv.HandleRenameFileTo)
	api.POST("/check-files-existence", srv.HandleCheckFilesExistence)
	api.POST("/create-m3u8", srv.HandleCreateM3U8File)

	// Image operations
	api.POST("/upload-image", srv.HandleUploadImage)
	api.POST("/upload-image-bytes", srv.HandleUploadImageBytes)
	api.GET("/read-image-base64", srv.HandleReadImageAsBase64)

	// Audio file upload
	api.POST("/upload-audio", srv.HandleUploadAudio)

	// System info
	api.GET("/os-info", srv.HandleGetOSInfo)
	api.POST("/files/open", srv.HandleOpenFileManager)

	// Server-Sent Events for real-time progress
	api.GET("/events", srv.HandleSSE)

	// Serve embedded frontend static files
	frontendFS, err := fs.Sub(frontendDist, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to get frontend filesystem: %v", err)
	}

	// Serve static files
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(frontendFS))))

	// Start server
	address := fmt.Sprintf(":%s", port)
	log.Printf("SpotiFLAC web server starting on http://localhost%s", address)
	log.Printf("Download path: %s", downloadPath)
	log.Printf("Data directory: %s", dataDir)

	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
