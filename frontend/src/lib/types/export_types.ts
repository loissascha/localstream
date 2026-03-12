export interface AuthResponse {
  token: string;
}
export interface VideoListResponse {
  videos: VideoListItem[];
}
export interface VideoListItem {
  id: string;
  name: string;
  size: number;
  mimeType: string;
}
export interface ShowInfo {
  id: string;
  name: string;
  year: number;
  description: string;
}
export interface ShowListResponse {
  shows: ShowInfo[];
}
export interface LibraryListItem {
  id: string;
  name: string;
  path: string;
  library_type: string;
}
export interface LibraryListResponse {
  libraries: LibraryListItem[];
}
export interface SeasonInfo {
  id: string;
  number: number;
}
export interface SeasonListResponse {
  seasons: SeasonInfo[];
}
export interface EpisodeInfo {
  id: string;
  number: number;
}
export interface EpisodeListResponse {
  episodes: EpisodeInfo[];
}
export interface SaveWatchstateRequest {
  show_id: string;
  season_id: string;
  episode_id: string;
  position: number;
  duration: number;
  finished: boolean;
}
export interface WatchstateResponse {
  id: string;
  show_id: string;
  season_id: string;
  episode_id: string;
  position: number;
  duration: number;
  finished: boolean;
  created_at: string;
  updated_at: string;
}
export interface WatchstateListResponse {
  watchstates: WatchstateResponse[];
}
export interface AuthUserResponse {
  id: number;
  username: string;
}
