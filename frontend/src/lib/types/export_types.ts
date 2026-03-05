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
