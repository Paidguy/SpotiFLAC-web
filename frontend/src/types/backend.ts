// Backend type definitions for FileManager operations
export namespace backend {
    export interface FileInfo {
        name: string;
        path: string;
        is_dir: boolean;
        size: number;
        children?: FileInfo[];
        expanded?: boolean;
    }

    export interface RenamePreview {
        old_name: string;
        new_name: string;
        error?: string;
    }

    export interface RenameResult {
        old_path: string;
        new_path: string;
        success: boolean;
        error?: string;
    }

    export interface AudioMetadata {
        title: string;
        artist: string;
        album: string;
        album_artist: string;
        track_number: number;
        disc_number: number;
        year: string;
    }
}
