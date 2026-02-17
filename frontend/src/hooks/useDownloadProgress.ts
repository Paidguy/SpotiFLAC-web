import { useState, useEffect } from "react";

export interface DownloadProgressInfo {
    is_downloading: boolean;
    mb_downloaded: number;
    speed_mbps: number;
}

export function useDownloadProgress() {
    const [progress, setProgress] = useState<DownloadProgressInfo>({
        is_downloading: false,
        mb_downloaded: 0,
        speed_mbps: 0,
    });

    useEffect(() => {
        // Setup SSE connection for real-time progress
        const eventSource = new EventSource('/api/events');

        eventSource.addEventListener('download:progress', (event: MessageEvent) => {
            try {
                const data = JSON.parse(event.data);
                // Update progress state from SSE event
                // data shape: {type, item_id, status, percent, speed, message}
                if (data.status === 'downloading') {
                    setProgress({
                        is_downloading: true,
                        mb_downloaded: data.percent || 0,
                        speed_mbps: data.speed || 0,
                    });
                } else if (data.status === 'completed' || data.status === 'failed' || data.status === 'exists') {
                    setProgress({
                        is_downloading: false,
                        mb_downloaded: 0,
                        speed_mbps: 0,
                    });
                }
            } catch (err) {
                console.error('Failed to parse download progress event:', err);
            }
        });

        eventSource.onerror = (error) => {
            console.error('SSE connection error in useDownloadProgress:', error);
            // EventSource will auto-reconnect
        };

        return () => {
            eventSource.close();
        };
    }, []);

    return progress;
}
