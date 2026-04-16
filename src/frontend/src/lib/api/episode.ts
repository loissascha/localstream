import type { EpisodeInfo, EpisodeMetadataInfo } from '$lib/types/export_types';

export async function getEpisodeDetails(
	bearerToken: string,
	episodeID: string
): Promise<EpisodeInfo> {
	const response = await fetch('/api/v1/episodes/' + episodeID, {
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

export async function getNextEpisode(
	bearerToken: string,
	currentEpisodeID: string
): Promise<EpisodeInfo> {
	const response = await fetch('/api/episodes/next/' + currentEpisodeID, {
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

export async function getEpisodeMetadata(
	bearerToken: string,
	showID: string,
	seasonNumber: number,
	episodeNumber: number
): Promise<EpisodeMetadataInfo> {
	const response = await fetch(
		`/api/v1/episode/metadata/${showID}/${seasonNumber}/${episodeNumber}`,
		{
			method: 'GET',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		}
	);
	if (response.status !== 200) {
		console.error(response);
		throw new Error('Error: ' + response.status);
	}
	const result = (await response.json()) as EpisodeMetadataInfo;
	return result;
}
