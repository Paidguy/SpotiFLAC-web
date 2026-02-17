import type {
	SpotifyMetadataResponse,
	DownloadRequest,
	DownloadResponse,
	HealthResponse,
	LyricsDownloadRequest,
	LyricsDownloadResponse,
	CoverDownloadRequest,
	CoverDownloadResponse,
	HeaderDownloadRequest,
	HeaderDownloadResponse,
	GalleryImageDownloadRequest,
	GalleryImageDownloadResponse,
	AvatarDownloadRequest,
	AvatarDownloadResponse,
	AnalysisResult,
	ConvertAudioRequest,
	ConvertAudioResponse,
} from "@/types/api";

// Base API URL - empty string means same origin
const API_BASE = "";

// Helper function to make API requests
async function apiRequest<T>(
	endpoint: string,
	options?: RequestInit
): Promise<T> {
	const response = await fetch(`${API_BASE}${endpoint}`, {
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
		...options,
	});

	if (!response.ok) {
		const errorText = await response.text();
		throw new Error(
			`API request failed: ${response.status} ${response.statusText} - ${errorText}`
		);
	}

	// Check if response is JSON or text
	const contentType = response.headers.get("content-type");
	if (contentType && contentType.includes("application/json")) {
		return response.json();
	} else {
		return response.text() as T;
	}
}

export async function fetchSpotifyMetadata(
	url: string,
	batch: boolean = true,
	delay: number = 1.0,
	timeout: number = 300.0
): Promise<SpotifyMetadataResponse> {
	const req = {
		url,
		batch,
		delay,
		timeout,
	};

	const jsonString = await apiRequest<string>("/api/metadata", {
		method: "POST",
		body: JSON.stringify(req),
	});

	return JSON.parse(jsonString);
}

export async function downloadTrack(
	request: DownloadRequest
): Promise<DownloadResponse> {
	return apiRequest<DownloadResponse>("/api/download", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

export async function checkHealth(): Promise<HealthResponse> {
	return apiRequest<HealthResponse>("/api/health");
}

export async function downloadLyrics(
	request: LyricsDownloadRequest
): Promise<LyricsDownloadResponse> {
	return apiRequest<LyricsDownloadResponse>("/api/lyrics", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

export async function downloadCover(
	request: CoverDownloadRequest
): Promise<CoverDownloadResponse> {
	return apiRequest<CoverDownloadResponse>("/api/cover", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

export async function downloadHeader(
	request: HeaderDownloadRequest
): Promise<HeaderDownloadResponse> {
	return apiRequest<HeaderDownloadResponse>("/api/header", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

export async function downloadGalleryImage(
	request: GalleryImageDownloadRequest
): Promise<GalleryImageDownloadResponse> {
	return apiRequest<GalleryImageDownloadResponse>("/api/gallery-image", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

export async function downloadAvatar(
	request: AvatarDownloadRequest
): Promise<AvatarDownloadResponse> {
	return apiRequest<AvatarDownloadResponse>("/api/avatar", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

// Search APIs
export async function SearchSpotify(query: string, limit: number = 10): Promise<any> {
	return apiRequest<any>("/api/search", {
		method: "POST",
		body: JSON.stringify({ query, limit }),
	});
}

export async function SearchSpotifyByType(
	query: string,
	search_type: string,
	limit: number = 50,
	offset: number = 0
): Promise<any[]> {
	return apiRequest<any[]>("/api/search-by-type", {
		method: "POST",
		body: JSON.stringify({ query, search_type, limit, offset }),
	});
}

// Download queue and progress
export async function GetDownloadProgress(): Promise<any> {
	return apiRequest<any>("/api/download-progress");
}

export async function GetDownloadQueue(): Promise<any> {
	return apiRequest<any>("/api/download-queue");
}

export async function ClearCompletedDownloads(): Promise<void> {
	await apiRequest<void>("/api/clear-completed", { method: "POST" });
}

export async function ClearAllDownloads(): Promise<void> {
	await apiRequest<void>("/api/clear-all", { method: "POST" });
}

export async function ExportFailedDownloads(): Promise<{ success: boolean; message: string; data: string }> {
	return apiRequest<any>("/api/export-failed");
}

export async function CancelAllQueuedItems(): Promise<void> {
	await apiRequest<void>("/api/cancel-queued", { method: "POST" });
}

// History APIs
export async function GetDownloadHistory(): Promise<any[]> {
	return apiRequest<any[]>("/api/history");
}

export async function ClearDownloadHistory(): Promise<void> {
	await apiRequest<void>("/api/history", { method: "DELETE" });
}

export async function DeleteDownloadHistoryItem(id: string): Promise<void> {
	await apiRequest<void>(`/api/history/${id}`, { method: "DELETE" });
}

export async function GetFetchHistory(): Promise<any[]> {
	return apiRequest<any[]>("/api/fetch-history");
}

export async function AddFetchHistory(item: any): Promise<void> {
	await apiRequest<void>("/api/fetch-history", {
		method: "POST",
		body: JSON.stringify(item),
	});
}

export async function ClearFetchHistory(): Promise<void> {
	await apiRequest<void>("/api/fetch-history", { method: "DELETE" });
}

export async function DeleteFetchHistoryItem(id: string): Promise<void> {
	await apiRequest<void>(`/api/fetch-history/${id}`, { method: "DELETE" });
}

export async function ClearFetchHistoryByType(itemType: string): Promise<void> {
	await apiRequest<void>(`/api/fetch-history/type/${itemType}`, { method: "DELETE" });
}

// Track availability and preview
export async function CheckTrackAvailability(spotifyTrackId: string): Promise<any> {
	const response = await fetch(`/api/track-availability?spotify_track_id=${encodeURIComponent(spotifyTrackId)}`);
	if (!response.ok) throw new Error("Failed to check track availability");
	return response.text();
}

export async function GetPreviewURL(trackId: string): Promise<string> {
	const response = await fetch(`/api/preview-url?track_id=${encodeURIComponent(trackId)}`);
	if (!response.ok) throw new Error("Failed to get preview URL");
	const data = await response.json();
	return data.preview_url || "";
}

export async function GetStreamingURLs(spotifyTrackId: string, region: string = ""): Promise<any> {
	const url = `/api/streaming-urls?spotify_track_id=${encodeURIComponent(spotifyTrackId)}${region ? `&region=${encodeURIComponent(region)}` : ""}`;
	const response = await fetch(url);
	if (!response.ok) throw new Error("Failed to get streaming URLs");
	return response.text();
}

// Audio analysis
export async function analyzeAudioFile(filePath: string): Promise<AnalysisResult> {
	const response = await fetch(`/api/analyze-track?file_path=${encodeURIComponent(filePath)}`);
	if (!response.ok) {
		const errorText = await response.text();
		throw new Error(`Failed to analyze audio: ${errorText}`);
	}
	const jsonString = await response.text();
	return JSON.parse(jsonString);
}

// FFmpeg
export async function CheckFFmpegInstalled(): Promise<{ installed: boolean }> {
	return apiRequest<{ installed: boolean }>("/api/ffmpeg/installed");
}

export async function DownloadFFmpeg(): Promise<any> {
	return apiRequest<any>("/api/ffmpeg/download", { method: "POST" });
}

// Audio conversion
export async function convertAudio(request: ConvertAudioRequest): Promise<ConvertAudioResponse> {
	return apiRequest<ConvertAudioResponse>("/api/convert-audio", {
		method: "POST",
		body: JSON.stringify(request),
	});
}

// File operations
export async function SelectFolder(): Promise<string> {
	// In web mode, return empty or server path
	return "";
}

export async function OpenFolder(path: string): Promise<void> {
	await apiRequest<void>("/api/files/open", {
		method: "POST",
		body: JSON.stringify({ path }),
	});
}

// System info
export async function GetOSInfo(): Promise<string> {
	const response = await fetch("/api/os-info");
	if (!response.ok) throw new Error("Failed to get OS info");
	const data = await response.json();
	return data.os || "";
}

// Image operations
export async function UploadImage(filePath: string): Promise<string> {
	const response = await fetch("/api/upload-image", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ file_path: filePath }),
	});
	if (!response.ok) throw new Error("Failed to upload image");
	const data = await response.json();
	return data.url || "";
}

export async function UploadImageBytes(filename: string, base64Data: string): Promise<string> {
	return apiRequest<any>("/api/upload-image-bytes", {
		method: "POST",
		body: JSON.stringify({ filename, base64_data: base64Data }),
	}).then(data => data.url || "");
}

// Placeholder functions for file selection (web mode doesn't support native dialogs)
export async function SelectFile(): Promise<string> {
	return "";
}

export async function SelectAudioFiles(): Promise<string[]> {
	return [];
}

export async function SelectImageVideo(): Promise<string[]> {
	return [];
}
