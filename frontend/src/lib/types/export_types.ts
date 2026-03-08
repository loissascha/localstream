export interface LibraryListItem {
  id: string;
  name: string;
  path: string;
  library_type: string;
}
export interface LibraryListResponse {
  libraries: LibraryListItem[];
}
export interface ShowInfo {
  id: string;
  name: string;
}
export interface ShowListResponse {
  shows: ShowInfo[];
}
export interface AuthUserResponse {
  id: number;
  username: string;
}
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
