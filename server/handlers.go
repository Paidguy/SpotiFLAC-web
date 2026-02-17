package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"spotiflac/backend"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Server represents the HTTP server with all its dependencies
type Server struct {
	sseBroker      *SSEBroker
	downloadPath   string
	dataDir        string
}

// NewServer creates a new server instance
func NewServer(downloadPath, dataDir string) *Server {
	broker := NewSSEBroker()
	go broker.Run()

	return &Server{
		sseBroker:    broker,
		downloadPath: downloadPath,
		dataDir:      dataDir,
	}
}

// getFirstArtist extracts the first artist from a delimited string
func getFirstArtist(artistString string) string {
	if artistString == "" {
		return ""
	}
	delimiters := []string{", ", " & ", " feat. ", " ft. ", " featuring "}
	for _, d := range delimiters {
		if idx := strings.Index(strings.ToLower(artistString), d); idx != -1 {
			return strings.TrimSpace(artistString[:idx])
		}
	}
	return artistString
}

// HandleSSE handles Server-Sent Events for real-time progress updates
func (s *Server) HandleSSE(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	// Create new client
	client := &SSEClient{
		ID:      uuid.New().String(),
		Channel: make(chan []byte, 256),
	}

	s.sseBroker.RegisterClient(client)
	defer s.sseBroker.UnregisterClient(client)

	// Keep connection alive
	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case msg := <-client.Channel:
			// Parse JSON to extract event type
			var eventData map[string]interface{}
			if err := json.Unmarshal(msg, &eventData); err == nil {
				if eventType, ok := eventData["type"].(string); ok {
					// Send as named event for addEventListener
					if _, err := fmt.Fprintf(c.Response(), "event: %s\ndata: %s\n\n", eventType, msg); err != nil {
						return err
					}
				} else {
					// Fallback to default message event
					if _, err := fmt.Fprintf(c.Response(), "data: %s\n\n", msg); err != nil {
						return err
					}
				}
			} else {
				// Not JSON, send as-is
				if _, err := fmt.Fprintf(c.Response(), "data: %s\n\n", msg); err != nil {
					return err
				}
			}
			c.Response().Flush()
		}
	}
}

// HandleHealth handles health check requests
func (s *Server) HandleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

// HandleGetDownloadPath returns the server's download path
func (s *Server) HandleGetDownloadPath(c echo.Context) error {
	return c.JSON(http.StatusOK, DownloadPathResponse{Path: s.downloadPath})
}

// HandleGetSpotifyMetadata handles metadata fetching requests
func (s *Server) HandleGetSpotifyMetadata(c echo.Context) error {
	var req SpotifyMetadataRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
	}

	if req.Timeout <= 0 {
		req.Timeout = 300.0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.Timeout)*time.Second)
	defer cancel()

	result, err := backend.GetFilteredSpotifyData(ctx, req.URL, req.Batch, time.Duration(req.Delay)*time.Second)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to encode metadata: %v", err)})
	}

	return c.String(http.StatusOK, string(jsonData))
}

// HandleGetStreamingURLs handles streaming URL fetching
func (s *Server) HandleGetStreamingURLs(c echo.Context) error {
	spotifyTrackID := c.QueryParam("spotify_track_id")
	region := c.QueryParam("region")

	if spotifyTrackID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Spotify track ID is required"})
	}

	client := backend.NewSongLinkClient()
	songlink, err := client.GetAllURLsFromSpotify(spotifyTrackID, region)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	jsonData, err := json.Marshal(songlink)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to encode response: %v", err)})
	}

	return c.String(http.StatusOK, string(jsonData))
}

// HandleSearchSpotify handles Spotify search requests
func (s *Server) HandleSearchSpotify(c echo.Context) error {
	var req SpotifySearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := backend.SearchSpotify(ctx, req.Query, req.Limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// HandleSearchSpotifyByType handles typed Spotify search requests
func (s *Server) HandleSearchSpotifyByType(c echo.Context) error {
	var req SpotifySearchByTypeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	if req.SearchType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search type is required"})
	}

	if req.Limit <= 0 {
		req.Limit = 50
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := backend.SearchSpotifyByType(ctx, req.Query, req.SearchType, req.Limit, req.Offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// HandleDownloadTrack handles track download requests
func (s *Server) HandleDownloadTrack(c echo.Context) error {
	var req DownloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, DownloadResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// SECURITY: Override client's output directory with server's configured path
	req.OutputDir = s.downloadPath

	if req.Service == "qobuz" && req.SpotifyID == "" {
		return c.JSON(http.StatusBadRequest, DownloadResponse{
			Success: false,
			Error:   "Spotify ID is required for Qobuz",
		})
	}

	if req.Service == "" {
		req.Service = "tidal"
	}

	if req.OutputDir == "" {
		req.OutputDir = s.downloadPath
	}

	if req.AudioFormat == "" {
		req.AudioFormat = "flac"
	}

	if req.FilenameFormat == "" {
		req.FilenameFormat = "{track_number}. {track_name}"
	}

	// Apply folder template if needed from settings
	// The actual folder creation logic should be in backend

	// Handle first artist only if requested
	if req.UseFirstArtistOnly && req.ArtistName != "" {
		req.ArtistName = getFirstArtist(req.ArtistName)
		if req.AlbumArtist != "" {
			req.AlbumArtist = getFirstArtist(req.AlbumArtist)
		}
	}

	// Create download item if ItemID is provided
	if req.ItemID != "" {
		backend.AddToQueue(req.ItemID, req.TrackName, req.ArtistName, req.AlbumName, req.SpotifyID)
		backend.StartDownloadItem(req.ItemID)

		// Set up global progress callback for this download
		backend.SetGlobalProgressCallback(func(itemID string, mbDownloaded, speedMBps float64) {
			// Calculate percentage if we have total size info
			percent := mbDownloaded // Approximate percentage, actual file size unknown until complete

			// Broadcast progress event
			s.sseBroker.BroadcastJSON(map[string]interface{}{
				"type":     "download:progress",
				"item_id":  itemID,
				"status":   "downloading",
				"percent":  percent,
				"speed":    speedMBps,
				"message":  fmt.Sprintf("Downloading: %.2f MB (%.2f MB/s)", mbDownloaded, speedMBps),
			})
		})

		// Broadcast start event
		s.sseBroker.BroadcastJSON(map[string]interface{}{
			"type":    "download:progress",
			"item_id": req.ItemID,
			"status":  "downloading",
			"percent": 0,
			"message": "Starting download...",
		})
	}

	// Perform the download
	var downloadErr error
	var filePath string
	var message string
	success := false

	switch req.Service {
	case "tidal":
		downloader := backend.NewTidalDownloader(req.ApiURL)
		filePath, downloadErr = downloader.Download(req.SpotifyID, req.OutputDir, req.Query, req.FilenameFormat, req.Position > 0, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, req.ServiceURL, req.AllowFallback, req.UseFirstArtistOnly)
	case "qobuz":
		downloader := backend.NewQobuzDownloader()
		filePath, downloadErr = downloader.DownloadTrack(req.SpotifyID, req.OutputDir, req.Query, req.FilenameFormat, req.Position > 0, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.UseAlbumTrackNumber, req.CoverURL, req.EmbedMaxQualityCover, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, req.ServiceURL, req.AllowFallback, req.UseFirstArtistOnly)
	case "amazon":
		downloader := backend.NewAmazonDownloader()
		filePath, downloadErr = downloader.DownloadBySpotifyID(req.SpotifyID, req.OutputDir, req.Query, req.FilenameFormat, "", "", req.Position > 0, req.Position, req.TrackName, req.ArtistName, req.AlbumName, req.AlbumArtist, req.ReleaseDate, req.CoverURL, req.SpotifyTrackNumber, req.SpotifyDiscNumber, req.SpotifyTotalTracks, req.EmbedMaxQualityCover, req.SpotifyTotalDiscs, req.Copyright, req.Publisher, req.ServiceURL, req.UseFirstArtistOnly)
	default:
		downloadErr = fmt.Errorf("unsupported service: %s", req.Service)
	}

	// Clear global callback after download completes
	if req.ItemID != "" {
		backend.SetGlobalProgressCallback(nil)
	}

	// Check if file already exists
	if downloadErr == nil && filePath != "" && strings.HasPrefix(filePath, "EXISTS:") {
		actualPath := strings.TrimPrefix(filePath, "EXISTS:")
		if req.ItemID != "" {
			backend.SkipDownloadItem(req.ItemID, actualPath)
			// Broadcast exists event
			s.sseBroker.BroadcastJSON(map[string]interface{}{
				"type":    "download:progress",
				"item_id": req.ItemID,
				"status":  "exists",
				"percent": 100,
				"message": "File already exists",
			})
		}
		return c.JSON(http.StatusOK, DownloadResponse{
			Success: true,
			Message: "File already exists",
			File:    actualPath,
			ItemID:  req.ItemID,
		})
	}

	if downloadErr != nil {
		if req.AllowFallback && req.ItemID != "" {
			// Return error but don't mark as failed yet - caller will handle fallback
			return c.JSON(http.StatusOK, DownloadResponse{
				Success: false,
				Error:   downloadErr.Error(),
				ItemID:  req.ItemID,
			})
		}

		if req.ItemID != "" {
			backend.FailDownloadItem(req.ItemID, downloadErr.Error())
			// Broadcast failure event
			s.sseBroker.BroadcastJSON(map[string]interface{}{
				"type":    "download:progress",
				"item_id": req.ItemID,
				"status":  "error",
				"message": downloadErr.Error(),
			})
		}

		return c.JSON(http.StatusOK, DownloadResponse{
			Success: false,
			Error:   downloadErr.Error(),
			ItemID:  req.ItemID,
		})
	}

	success = true
	message = "Download completed successfully"

	if req.ItemID != "" {
		// Get file size
		var finalSize float64
		if fileInfo, err := os.Stat(filePath); err == nil {
			finalSize = float64(fileInfo.Size()) / (1024 * 1024) // Convert to MB
		}
		backend.CompleteDownloadItem(req.ItemID, filePath, finalSize)
		// Broadcast completion event
		s.sseBroker.BroadcastJSON(map[string]interface{}{
			"type":    "download:progress",
			"item_id": req.ItemID,
			"status":  "done",
			"message": "Download completed",
			"percent": 100,
		})
	}

	// Add to download history
	historyItem := backend.HistoryItem{
		Title:     req.TrackName,
		Artists:   req.ArtistName,
		Album:     req.AlbumName,
		Quality:   req.AudioFormat,
		Format:    req.AudioFormat,
		Timestamp: time.Now().Unix(),
		Path:      filePath,
		SpotifyID: req.SpotifyID,
	}
	backend.AddHistoryItem(historyItem, "SpotiFLAC")

	return c.JSON(http.StatusOK, DownloadResponse{
		Success: success,
		Message: message,
		File:    filePath,
		ItemID:  req.ItemID,
	})
}

// HandleDownloadLyrics handles lyrics download requests
func (s *Server) HandleDownloadLyrics(c echo.Context) error {
	var req LyricsDownloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, backend.LyricsDownloadResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	if req.SpotifyID == "" {
		return c.JSON(http.StatusBadRequest, backend.LyricsDownloadResponse{
			Success: false,
			Error:   "Spotify ID is required",
		})
	}

	client := backend.NewLyricsClient()
	backendReq := backend.LyricsDownloadRequest{
		SpotifyID:           req.SpotifyID,
		TrackName:           req.TrackName,
		ArtistName:          req.ArtistName,
		AlbumName:           req.AlbumName,
		AlbumArtist:         req.AlbumArtist,
		ReleaseDate:         req.ReleaseDate,
		OutputDir:           req.OutputDir,
		FilenameFormat:      req.FilenameFormat,
		TrackNumber:         req.TrackNumber,
		Position:            req.Position,
		UseAlbumTrackNumber: req.UseAlbumTrackNumber,
		DiscNumber:          req.DiscNumber,
	}

	resp, err := client.DownloadLyrics(backendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, backend.LyricsDownloadResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// HandleDownloadCover handles cover art download requests
func (s *Server) HandleDownloadCover(c echo.Context) error {
	var req CoverDownloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, backend.CoverDownloadResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	if req.CoverURL == "" {
		return c.JSON(http.StatusBadRequest, backend.CoverDownloadResponse{
			Success: false,
			Error:   "Cover URL is required",
		})
	}

	client := backend.NewCoverClient()
	backendReq := backend.CoverDownloadRequest{
		CoverURL:       req.CoverURL,
		TrackName:      req.TrackName,
		ArtistName:     req.ArtistName,
		AlbumName:      req.AlbumName,
		AlbumArtist:    req.AlbumArtist,
		ReleaseDate:    req.ReleaseDate,
		OutputDir:      req.OutputDir,
		FilenameFormat: req.FilenameFormat,
		TrackNumber:    req.TrackNumber,
		Position:       req.Position,
		DiscNumber:     req.DiscNumber,
	}

	resp, err := client.DownloadCover(backendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, backend.CoverDownloadResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// HandleDownloadHeader handles header image download requests
func (s *Server) HandleDownloadHeader(c echo.Context) error {
	var req HeaderDownloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, backend.HeaderDownloadResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	if req.HeaderURL == "" {
		return c.JSON(http.StatusBadRequest, backend.HeaderDownloadResponse{
			Success: false,
			Error:   "Header URL is required",
		})
	}

	if req.ArtistName == "" {
		return c.JSON(http.StatusBadRequest, backend.HeaderDownloadResponse{
			Success: false,
			Error:   "Artist name is required",
		})
	}

	client := backend.NewCoverClient()
	backendReq := backend.HeaderDownloadRequest{
		HeaderURL:  req.HeaderURL,
		ArtistName: req.ArtistName,
		OutputDir:  req.OutputDir,
	}

	resp, err := client.DownloadHeader(backendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, backend.HeaderDownloadResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// HandleDownloadGalleryImage handles gallery image download requests
func (s *Server) HandleDownloadGalleryImage(c echo.Context) error {
	var req GalleryImageDownloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, backend.GalleryImageDownloadResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	if req.ImageURL == "" {
		return c.JSON(http.StatusBadRequest, backend.GalleryImageDownloadResponse{
			Success: false,
			Error:   "Image URL is required",
		})
	}

	if req.ArtistName == "" {
		return c.JSON(http.StatusBadRequest, backend.GalleryImageDownloadResponse{
			Success: false,
			Error:   "Artist name is required",
		})
	}

	client := backend.NewCoverClient()
	backendReq := backend.GalleryImageDownloadRequest{
		ImageURL:   req.ImageURL,
		ArtistName: req.ArtistName,
		ImageIndex: req.ImageIndex,
		OutputDir:  req.OutputDir,
	}

	resp, err := client.DownloadGalleryImage(backendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, backend.GalleryImageDownloadResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// HandleDownloadAvatar handles avatar download requests
func (s *Server) HandleDownloadAvatar(c echo.Context) error {
	var req AvatarDownloadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, backend.AvatarDownloadResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	if req.AvatarURL == "" {
		return c.JSON(http.StatusBadRequest, backend.AvatarDownloadResponse{
			Success: false,
			Error:   "Avatar URL is required",
		})
	}

	if req.ArtistName == "" {
		return c.JSON(http.StatusBadRequest, backend.AvatarDownloadResponse{
			Success: false,
			Error:   "Artist name is required",
		})
	}

	client := backend.NewCoverClient()
	backendReq := backend.AvatarDownloadRequest{
		AvatarURL:  req.AvatarURL,
		ArtistName: req.ArtistName,
		OutputDir:  req.OutputDir,
	}

	resp, err := client.DownloadAvatar(backendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, backend.AvatarDownloadResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// HandleGetDownloadProgress returns the current download progress
func (s *Server) HandleGetDownloadProgress(c echo.Context) error {
	progress := backend.GetDownloadProgress()
	return c.JSON(http.StatusOK, progress)
}

// HandleGetDownloadQueue returns the current download queue
func (s *Server) HandleGetDownloadQueue(c echo.Context) error {
	queue := backend.GetDownloadQueue()
	return c.JSON(http.StatusOK, queue)
}

// HandleClearCompletedDownloads clears completed downloads from the queue
func (s *Server) HandleClearCompletedDownloads(c echo.Context) error {
	backend.ClearAllDownloads()
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleClearAllDownloads clears all downloads from the queue
func (s *Server) HandleClearAllDownloads(c echo.Context) error {
	backend.ClearAllDownloads()
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleCancelAllQueuedItems cancels all queued items
func (s *Server) HandleCancelAllQueuedItems(c echo.Context) error {
	backend.CancelAllQueuedItems()
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleSkipDownloadItem skips a download item
func (s *Server) HandleSkipDownloadItem(c echo.Context) error {
	itemID := c.QueryParam("item_id")
	filePath := c.QueryParam("file_path")

	if itemID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "item_id is required"})
	}

	backend.SkipDownloadItem(itemID, filePath)
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleExportFailedDownloads exports failed downloads
func (s *Server) HandleExportFailedDownloads(c echo.Context) error {
	queue := backend.GetDownloadQueue()
	var failedItems []backend.DownloadItem
	for _, item := range queue.Queue {
		if item.Status == backend.StatusFailed {
			failedItems = append(failedItems, item)
		}
	}

	if len(failedItems) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "No failed downloads to export",
			"data":    "",
		})
	}

	var exportLines []string
	exportLines = append(exportLines, "# Failed Downloads")
	exportLines = append(exportLines, fmt.Sprintf("# Exported: %s", time.Now().Format("2006-01-02 15:04:05")))
	exportLines = append(exportLines, "")

	for _, item := range failedItems {
		exportLines = append(exportLines, fmt.Sprintf("Track: %s", item.TrackName))
		exportLines = append(exportLines, fmt.Sprintf("Artist: %s", item.ArtistName))
		exportLines = append(exportLines, fmt.Sprintf("Album: %s", item.AlbumName))
		exportLines = append(exportLines, fmt.Sprintf("Spotify ID: %s", item.SpotifyID))
		exportLines = append(exportLines, fmt.Sprintf("Error: %s", item.ErrorMessage))
		exportLines = append(exportLines, "---")
		exportLines = append(exportLines, "")
	}

	exportData := strings.Join(exportLines, "\n")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Exported %d failed downloads", len(failedItems)),
		"data":    exportData,
	})
}

// HandleGetDefaults returns default settings
func (s *Server) HandleGetDefaults(c echo.Context) error {
	defaults := map[string]string{
		"downloadPath": s.downloadPath,
		"audioFormat":  "flac",
	}
	return c.JSON(http.StatusOK, defaults)
}

// HandleLoadSettings loads settings from file
func (s *Server) HandleLoadSettings(c echo.Context) error {
	configPath, err := backend.GetFFmpegDir()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	settingsFile := filepath.Join(configPath, "settings.json")
	data, err := os.ReadFile(settingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default settings
			return c.JSON(http.StatusOK, map[string]interface{}{})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, settings)
}

// HandleSaveSettings saves settings to file
func (s *Server) HandleSaveSettings(c echo.Context) error {
	var settings map[string]interface{}
	if err := c.Bind(&settings); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	configPath, err := backend.GetFFmpegDir()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	settingsFile := filepath.Join(configPath, "settings.json")
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := os.WriteFile(settingsFile, data, 0644); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleGetHistory returns download history
func (s *Server) HandleGetHistory(c echo.Context) error {
	history, err := backend.GetHistoryItems("SpotiFLAC")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, history)
}

// HandleDeleteHistory deletes all history
func (s *Server) HandleDeleteHistory(c echo.Context) error {
	if err := backend.ClearHistory("SpotiFLAC"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleDeleteHistoryItem deletes a specific history item
func (s *Server) HandleDeleteHistoryItem(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	if err := backend.DeleteHistoryItem(id, "SpotiFLAC"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleGetFetchHistory returns fetch history
func (s *Server) HandleGetFetchHistory(c echo.Context) error {
	history, err := backend.GetFetchHistoryItems("SpotiFLAC")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, history)
}

// HandleClearFetchHistory clears fetch history
func (s *Server) HandleClearFetchHistory(c echo.Context) error {
	if err := backend.ClearFetchHistory("SpotiFLAC"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleDeleteFetchHistoryItem deletes a specific fetch history item
func (s *Server) HandleDeleteFetchHistoryItem(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	if err := backend.DeleteFetchHistoryItem(id, "SpotiFLAC"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleClearFetchHistoryByType clears fetch history by type
func (s *Server) HandleClearFetchHistoryByType(c echo.Context) error {
	itemType := c.Param("type")
	if itemType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type is required"})
	}

	if err := backend.ClearFetchHistoryByType(itemType, "SpotiFLAC"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleAddFetchHistory adds an item to fetch history
func (s *Server) HandleAddFetchHistory(c echo.Context) error {
	var item backend.FetchHistoryItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := backend.AddFetchHistoryItem(item, "SpotiFLAC"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleCheckTrackAvailability checks if a track is available
func (s *Server) HandleCheckTrackAvailability(c echo.Context) error {
	spotifyTrackID := c.QueryParam("spotify_track_id")

	if spotifyTrackID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Spotify track ID is required"})
	}

	client := backend.NewSongLinkClient()
	availability, err := client.CheckTrackAvailability(spotifyTrackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	jsonData, err := json.Marshal(availability)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to encode response: %v", err)})
	}

	return c.String(http.StatusOK, string(jsonData))
}

// HandleGetPreviewURL gets the preview URL for a track
func (s *Server) HandleGetPreviewURL(c echo.Context) error {
	trackID := c.QueryParam("track_id")

	if trackID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Track ID is required"})
	}

	previewURL, err := backend.GetPreviewURL(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"preview_url": previewURL})
}

// HandleAnalyzeTrack analyzes an audio track
func (s *Server) HandleAnalyzeTrack(c echo.Context) error {
	filePath := c.QueryParam("file_path")

	if filePath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File path is required"})
	}

	analysis, err := backend.AnalyzeTrack(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	jsonData, err := json.Marshal(analysis)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to encode response: %v", err)})
	}

	return c.String(http.StatusOK, string(jsonData))
}

// HandleAnalyzeMultipleTracks analyzes multiple audio tracks
func (s *Server) HandleAnalyzeMultipleTracks(c echo.Context) error {
	var req struct {
		FilePaths []string `json:"file_paths"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if len(req.FilePaths) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File paths are required"})
	}

	var results []map[string]interface{}

	for _, filePath := range req.FilePaths {
		analysis, err := backend.AnalyzeTrack(filePath)
		if err != nil {
			results = append(results, map[string]interface{}{
				"file_path": filePath,
				"error":     err.Error(),
			})
		} else {
			resultMap := map[string]interface{}{
				"file_path": filePath,
			}
			// Convert analysis struct to map
			jsonData, _ := json.Marshal(analysis)
			json.Unmarshal(jsonData, &resultMap)
			results = append(results, resultMap)
		}
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to encode response: %v", err)})
	}

	return c.String(http.StatusOK, string(jsonData))
}

// HandleCheckFFmpegInstalled checks if FFmpeg is installed
func (s *Server) HandleCheckFFmpegInstalled(c echo.Context) error {
	installed, err := backend.IsFFmpegInstalled()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"installed": false,
			"error":     err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]bool{"installed": installed})
}

// HandleIsFFprobeInstalled checks if FFprobe is installed
func (s *Server) HandleIsFFprobeInstalled(c echo.Context) error {
	installed, err := backend.IsFFprobeInstalled()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"installed": false,
			"error":     err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]bool{"installed": installed})
}

// HandleGetFFmpegPath gets the FFmpeg path
func (s *Server) HandleGetFFmpegPath(c echo.Context) error {
	path, err := backend.GetFFmpegPath()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"path": path})
}

// HandleDownloadFFmpeg downloads FFmpeg
func (s *Server) HandleDownloadFFmpeg(c echo.Context) error {
	// This is a placeholder - actual implementation depends on backend
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": false,
		"message": "FFmpeg download not implemented for web server",
	})
}

// HandleConvertAudio converts audio files
func (s *Server) HandleConvertAudio(c echo.Context) error {
	var req ConvertAudioRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	backendReq := backend.ConvertAudioRequest{
		InputFiles:   req.InputFiles,
		OutputFormat: req.OutputFormat,
		Bitrate:      req.Bitrate,
		Codec:        req.Codec,
	}

	results, err := backend.ConvertAudio(backendReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// HandleGetFileSizes gets file sizes
func (s *Server) HandleGetFileSizes(c echo.Context) error {
	var req struct {
		Files []string `json:"files"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	sizes := backend.GetFileSizes(req.Files)
	return c.JSON(http.StatusOK, sizes)
}

// HandleListDirectoryFiles lists files in a directory
func (s *Server) HandleListDirectoryFiles(c echo.Context) error {
	dirPath := c.QueryParam("dir_path")

	if dirPath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Directory path is required"})
	}

	files, err := backend.ListDirectory(dirPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, files)
}

// HandleListAudioFilesInDir lists audio files in a directory
func (s *Server) HandleListAudioFilesInDir(c echo.Context) error {
	dirPath := c.QueryParam("dir_path")

	if dirPath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Directory path is required"})
	}

	files, err := backend.ListAudioFiles(dirPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, files)
}

// HandleReadFileMetadata reads file metadata
func (s *Server) HandleReadFileMetadata(c echo.Context) error {
	filePath := c.QueryParam("file_path")

	if filePath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File path is required"})
	}

	metadata, err := backend.ReadAudioMetadata(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, metadata)
}

// HandlePreviewRenameFiles previews file renaming
func (s *Server) HandlePreviewRenameFiles(c echo.Context) error {
	var req struct {
		Files  []string `json:"files"`
		Format string   `json:"format"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	preview := backend.PreviewRename(req.Files, req.Format)
	return c.JSON(http.StatusOK, preview)
}

// HandleRenameFilesByMetadata renames files by metadata
func (s *Server) HandleRenameFilesByMetadata(c echo.Context) error {
	var req struct {
		Files  []string `json:"files"`
		Format string   `json:"format"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	results := backend.RenameFiles(req.Files, req.Format)
	return c.JSON(http.StatusOK, results)
}

// HandleReadTextFile reads a text file
func (s *Server) HandleReadTextFile(c echo.Context) error {
	filePath := c.QueryParam("file_path")

	if filePath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File path is required"})
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"content": string(content)})
}

// HandleRenameFileTo renames a file
func (s *Server) HandleRenameFileTo(c echo.Context) error {
	var req struct {
		OldPath string `json:"old_path"`
		NewName string `json:"new_name"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	dir := filepath.Dir(req.OldPath)
	ext := filepath.Ext(req.OldPath)
	newPath := filepath.Join(dir, req.NewName+ext)

	if err := os.Rename(req.OldPath, newPath); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleUploadImage uploads an image
func (s *Server) HandleUploadImage(c echo.Context) error {
	filePath := c.QueryParam("file_path")

	if filePath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File path is required"})
	}

	url, err := backend.UploadToSendNow(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

// HandleUploadImageBytes uploads image bytes
func (s *Server) HandleUploadImageBytes(c echo.Context) error {
	var req struct {
		Filename   string `json:"filename"`
		Base64Data string `json:"base64_data"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	imageData, err := base64.StdEncoding.DecodeString(req.Base64Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid base64 data"})
	}

	url, err := backend.UploadBytesToSendNow(req.Filename, imageData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

// HandleReadImageAsBase64 reads an image as base64
func (s *Server) HandleReadImageAsBase64(c echo.Context) error {
	filePath := c.QueryParam("file_path")

	if filePath == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File path is required"})
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	base64Data := base64.StdEncoding.EncodeToString(data)
	return c.JSON(http.StatusOK, map[string]string{"data": base64Data})
}

// HandleCheckFilesExistence checks if files exist
func (s *Server) HandleCheckFilesExistence(c echo.Context) error {
	var req struct {
		OutputDir string                      `json:"output_dir"`
		RootDir   string                      `json:"root_dir"`
		Tracks    []CheckFileExistenceRequest `json:"tracks"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	var results []CheckFileExistenceResult

	for i, track := range req.Tracks {
		// Build the expected file path based on the filename format
		filename := backend.BuildExpectedFilename(
			track.TrackName,
			track.ArtistName,
			track.AlbumName,
			"", // albumArtist
			"", // releaseDate
			track.FilenameFormat,
			"", // playlistName
			"", // playlistOwner
			track.UseAlbumTrackNumber,
			track.Position,
			track.DiscNumber,
			track.UseAlbumTrackNumber,
		)

		// Add extension
		filename += "." + track.Format

		filePath := filepath.Join(req.OutputDir, filename)
		if req.RootDir != "" {
			filePath = filepath.Join(req.RootDir, filename)
		}

		_, err := os.Stat(filePath)
		exists := err == nil

		results = append(results, CheckFileExistenceResult{
			Exists:   exists,
			FilePath: filePath,
			Index:    i,
		})
	}

	return c.JSON(http.StatusOK, results)
}

// HandleCreateM3U8File creates an M3U8 playlist file
func (s *Server) HandleCreateM3U8File(c echo.Context) error {
	var req M3U8Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// SECURITY: Override output directory with server's configured path
	req.OutputDir = s.downloadPath

	// TODO: Implement M3U8 creation - not currently in backend
	return c.JSON(http.StatusNotImplemented, map[string]string{
		"error": "M3U8 playlist creation not yet implemented",
	})
}

// HandleGetOSInfo returns OS information
func (s *Server) HandleGetOSInfo(c echo.Context) error {
	osInfo, err := backend.GetOSInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"os": osInfo})
}

// HandleUploadAudio handles audio file uploads from browser
func (s *Server) HandleUploadAudio(c echo.Context) error {
	// Get the file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file provided"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	// Create uploads directory in server's download path
	uploadsDir := filepath.Join(s.downloadPath, "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create uploads directory"})
	}

	// Create destination file
	dstPath := filepath.Join(uploadsDir, filepath.Base(file.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create destination file"})
	}
	defer dst.Close()

	// Copy file contents
	if _, err := dst.ReadFrom(src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"path":     dstPath,
		"filename": filepath.Base(file.Filename),
	})
}

// HandleOpenFileManager opens the file manager (no-op in web mode)
func (s *Server) HandleOpenFileManager(c echo.Context) error {
	// In web mode, we can't open the file manager on the server
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "skipped",
		"message": "File manager cannot be opened in web mode",
	})
}
