import { useState, useEffect, useRef } from "react";
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
    const intervalRef = useRef<number | null>(null);
    useEffect(() => {
        const pollProgress = async () => {
            try {
                const response = await fetch("/api/download-progress");
                if (!response.ok) {
                    throw new Error("Failed to get download progress");
                }
                const progressInfo = await response.json();
                setProgress(progressInfo);
            }
            catch (error) {
                console.error("Failed to get download progress:", error);
            }
        };
        intervalRef.current = window.setInterval(pollProgress, 200);
        pollProgress();
        return () => {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
            }
        };
    }, []);
    return progress;
}
