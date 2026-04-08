import type { MovieMetadataInfo } from '$lib/types/export_types';

export async function loadMovieMetadata(
	bearerToken: string,
	movieId: string
): Promise<MovieMetadataInfo[]> {
	try {
		const res = await fetch('/api/v1/movie/metadata/' + movieId, {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load metadata: ${res.status}`);
		}

		const data = (await res.json()) as MovieMetadataInfo[];
		return data;
	} catch (error) {
		throw error;
	}
}
