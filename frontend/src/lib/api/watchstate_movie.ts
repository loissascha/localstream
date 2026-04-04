import { API_URL } from '$lib/consts';
import type { SaveMovieWatchstateRequest, WatchstateInfo } from '$lib/types/export_types';

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
