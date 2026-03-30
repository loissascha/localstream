export interface AuthResponse {
  token: string;
}
export interface AuthUserIsAdminResponse {
  id: number;
  is_admin: boolean;
}
export interface AuthUserResponse {
  id: number;
  username: string;
}
export interface CreateLibraryRequest {
  name: string;
  type: string;
  path: string;
}
export interface CreateLibraryResponse {
  library: LibraryListItem;
}
export interface EpisodeInfo {
  id: string;
  season_id: string;
  number: number;
  watchstate: WatchstateInfo;
}
export interface EpisodeListResponse {
  episodes: EpisodeInfo[];
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
export interface SaveWatchstateRequest {
  show_id: string;
  season_id: string;
  episode_id: string;
  position: number;
  duration: number;
  finished: boolean;
}
export interface SeasonInfo {
  id: string;
  number: number;
}
export interface SeasonListResponse {
  seasons: SeasonInfo[];
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
export interface VideoListItem {
  id: string;
  name: string;
  size: number;
  mimeType: string;
}
export interface VideoListResponse {
  videos: VideoListItem[];
}
export interface WatchstateInfo {
  position: number;
  duration: number;
  percentage: number;
  finished: boolean;
}
export interface WatchstateListResponse {
  watchstates: WatchstateResponse[];
}
export interface WatchstateResponse {
  id: string;
  show_id: string;
  show_info: ShowInfo;
  season_id: string;
  season_info: SeasonInfo;
  episode_id: string;
  episode_info: EpisodeInfo;
  position: number;
  duration: number;
  finished: boolean;
  created_at: string;
  updated_at: string;
  percentage: number;
}
