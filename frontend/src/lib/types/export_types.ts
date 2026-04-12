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
export interface EpisodeMetadataInfo {
  id: string;
  episode_id: string;
  url: string;
  name: string;
  number: number;
  summary: string;
  medium_image_url: string;
  original_image_url: string;
  fetch_id: number;
  fetch_source: string;
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
export interface MovieInfo {
  id: string;
  name: string;
  year: number;
  description: string;
  fetch_source: string;
}
export interface MovieListResponse {
  movies: MovieInfo[];
}
export interface MovieMetadataInfo {
  id: string;
  movie_id: string;
  name: string;
  url: string;
  description: string;
  medium_image_url: string;
  backdrop_image_url: string;
  fetch_source: string;
}
export interface SaveMovieWatchstateRequest {
  movie_id: string;
  position: number;
  duration: number;
  finished: boolean;
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
export interface SeasonMetadataInfo {
  id: string;
  season_id: string;
  url: string;
  number: number;
  summary: string;
  premiere_date: string;
  medium_image_url: string;
  original_image_url: string;
  fetch_id: number;
  fetch_source: string;
}
export interface ShowInfo {
  id: string;
  name: string;
  year: number;
  fetch_source: string;
}
export interface ShowListResponse {
  shows: ShowInfo[];
}
export interface ShowMetadataInfo {
  id: string;
  show_id: string;
  name: string;
  url: string;
  description: string;
  medium_image_url: string;
  original_image_url: string;
  fetch_source: string;
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
export interface WatchstateMovieResponse {
  id: string;
  movie_id: string;
  position: number;
  duration: number;
  finished: boolean;
  created_at: string;
  updated_at: string;
  percentage: number;
  movie_info: MovieInfo;
}
export interface WatchstateMoviesListResponse {
  watchstates: WatchstateMovieResponse[];
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
