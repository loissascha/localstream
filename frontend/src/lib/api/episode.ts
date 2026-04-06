import { API_URL } from '$lib/consts';
import type { EpisodeInfo } from '$lib/types/export_types';

export async function getNextEpisode(
	bearerToken: string,
	currentEpisodeID: string
): Promise<EpisodeInfo> {
	const response = await fetch(API_URL + '/api/episodes/next/' + currentEpisodeID, {
		method: 'GET',
		headers: {
			Authorization: 'Bearer ' + bearerToken
		}
	});
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as EpisodeInfo;
	return result;
}
