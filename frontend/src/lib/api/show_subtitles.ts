import type { SubtitleProviderResult } from '$lib/types/export_types';

export async function searchEpisodeSubtitles(
	bearerToken: string,
	showName: string,
	seasonNumber: number,
	episodeNumber: number,
	lang: string
): Promise<SubtitleProviderResult[]> {
	try {
		const params = new URLSearchParams({
			show: showName,
			season: seasonNumber.toString(),
			episode: episodeNumber.toString(),
			lang
		});
		const res = await fetch(`/api/v1/show/subtitles/search?${params.toString()}`, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken,
				'Content-Type': 'application/json'
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to search for subtitle: ${res.status}`);
		}

		const data = (await res.json()) as SubtitleProviderResult[];
		return data;
	} catch (error) {
		throw error;
	}
}
