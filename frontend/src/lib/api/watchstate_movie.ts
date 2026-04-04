import { API_URL } from '$lib/consts';
import type {
	SaveMovieWatchstateRequest,
	WatchstateInfo,
	WatchstateMovieResponse,
	WatchstateMoviesListResponse,
	WatchstateResponse
} from '$lib/types/export_types';

export async function updateWatchstateMovie(
	bearerToken: string,
	body: SaveMovieWatchstateRequest
): Promise<WatchstateInfo> {
	const response = await fetch(API_URL + '/api/v1/movie/watchstate', {
		method: 'POST',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		},
		body: JSON.stringify(body)
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateInfo;
	return result;
}

export async function listLatestWatchstateByMovie(
	bearerToken: string
): Promise<WatchstateMovieResponse[]> {
	const response = await fetch(API_URL + '/api/v1/watchstate/movie/latest', {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as WatchstateMoviesListResponse;
	return result.watchstates;
}

export async function getWatchstateForMovie(
	bearerToken: string,
	movieId: string
): Promise<WatchstateInfo> {
	const response = await fetch(API_URL + '/api/v1/watchstate/movie/' + movieId, {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('' + response.status);
	}
	const result = (await response.json()) as WatchstateInfo;
	return result;
}
