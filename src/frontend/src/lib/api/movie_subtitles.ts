import type { MovieSubtitleInfo, SubtitleProviderResult } from '$lib/types/export_types';

export async function loadMovieSubtitles(
	bearerToken: string,
	movieId: string
): Promise<MovieSubtitleInfo[]> {
	try {
		const res = await fetch('/api/v1/movie/subtitles/' + movieId, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load metadata: ${res.status}`);
		}

		const data = (await res.json()) as MovieSubtitleInfo[];
		return data;
	} catch (error) {
		throw error;
	}
}

export async function searchMovieSubtitles(
	bearerToken: string,
	searchTerm: string
): Promise<SubtitleProviderResult[]> {
	try {
		const res = await fetch('/api/v1/movie/subtitles/search?q=' + searchTerm, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken,
				'Content-Type': 'application/json'
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load metadata: ${res.status}`);
		}

		const data = (await res.json()) as SubtitleProviderResult[];
		return data;
	} catch (error) {
		throw error;
	}
}
