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
