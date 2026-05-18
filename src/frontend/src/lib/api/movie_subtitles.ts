import type { SubtitleInfo, SubtitleProviderResult } from '$lib/types/export_types';

export async function loadMovieSubtitles(
	bearerToken: string,
	movieId: string
): Promise<SubtitleInfo[]> {
	try {
		const res = await fetch('/api/v1/movie/subtitles/' + movieId, {
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

export async function searchMovieSubtitles(
	bearerToken: string,
	searchTerm: string,
	lang: string
): Promise<SubtitleProviderResult[]> {
	try {
		const params = new URLSearchParams({
			q: searchTerm,
			lang
		});
		const res = await fetch(`/api/v1/movie/subtitles/search?${params.toString()}`, {
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

export async function downloadMovieSubtitle(
	bearerToken: string,
	movieId: string,
	body: SubtitleProviderResult
) {
	try {
		const res = await fetch('/api/v1/movie/subtitles/' + movieId + '/create', {
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
