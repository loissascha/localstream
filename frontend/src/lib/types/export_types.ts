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
export interface LibraryListItem {
  id: string;
  name: string;
  path: string;
}
export interface LibraryListResponse {
  libraries: LibraryListItem[];
}
export interface AuthUserResponse {
  id: number;
  username: string;
}
