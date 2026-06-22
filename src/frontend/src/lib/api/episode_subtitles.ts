import type { SubtitleProviderResult } from '$lib/types/export_types';

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
