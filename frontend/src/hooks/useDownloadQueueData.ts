import { useEffect, useState } from "react";

type DownloadStatus = "queued" | "downloading" | "completed" | "failed" | "skipped";

interface DownloadItem {
    id: string;
    track_name: string;
    artist_name: string;
    album_name: string;
    spotify_id: string;
    status: DownloadStatus;
    progress: number;
    total_size: number;
    speed: number;
    start_time: number;
    end_time: number;
    error_message: string;
    file_path: string;
}

interface DownloadQueueInfo {
    is_downloading: boolean;
    queue: DownloadItem[];
    current_speed: number;
    total_downloaded: number;
    session_start_time: number;
    queued_count: number;
    completed_count: number;
    failed_count: number;
    skipped_count: number;
}

export function useDownloadQueueData() {
    const [queueInfo, setQueueInfo] = useState<DownloadQueueInfo>({
        is_downloading: false,
        queue: [],
        current_speed: 0,
        total_downloaded: 0,
        session_start_time: 0,
        queued_count: 0,
        completed_count: 0,
        failed_count: 0,
        skipped_count: 0,
    });
    useEffect(() => {
        const fetchQueue = async () => {
            try {
                const response = await fetch("/api/download-queue");
                if (!response.ok) {
                    throw new Error("Failed to get download queue");
                }
                const info = await response.json();
                setQueueInfo(info);
            }
            catch (error) {
                console.error("Failed to get download queue:", error);
            }
        };
        fetchQueue();
        const interval = setInterval(fetchQueue, 200);
        return () => clearInterval(interval);
    }, []);
    return queueInfo;
}
