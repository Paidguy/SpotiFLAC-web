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
