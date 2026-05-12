import type { MovieSubtitleInfo } from '$lib/types/export_types';

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
