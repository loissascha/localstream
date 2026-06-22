import type { SubtitleInfo, SubtitleProviderResult } from '$lib/types/export_types';

export async function loadEpisodeSubtitles(
	bearerToken: string,
	episodeId: string
): Promise<SubtitleInfo[]> {
	try {
		const res = await fetch('/api/v1/show/subtitles/' + episodeId, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load metadata: ${res.status}`);
		}

		const data = (await res.json()) as SubtitleInfo[];
		return data;
	} catch (error) {
		throw error;
	}
}

export async function downloadEpisodeSubtitle(
	bearerToken: string,
	episodeId: string,
	body: SubtitleProviderResult
) {
	try {
		const res = await fetch('/api/v1/show/subtitles/' + episodeId + '/create', {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(body)
		});
		if (!res.ok) {
			throw new Error(`Failed to download: ${res.status}`);
		}
	} catch (error) {
		throw error;
	}
}
