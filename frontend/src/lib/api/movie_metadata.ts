import type { MovieMetadataInfo, MovieResult } from '$lib/types/export_types';

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

export async function setPrimaryMovieMetadataByFetchID(
	bearerToken: string,
	movieID: string,
	fetchID: number
) {
	try {
		const res = await fetch(`/api/v1/movie/metadata/${movieID}/set/primary/by-fetchid/${fetchID}`, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to set fetch id: ${res.status}`);
		}
	} catch (e) {
		throw e;
	}
}

export async function searchMovieMetadata(bearerToken: string, searchQuery: string) {
	try {
		const res = await fetch('/api/v1/movie/metadata/search?q=' + searchQuery, {
			method: 'POST',
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			throw new Error(`Failed to load metadata: ${res.status}`);
		}

		const data = (await res.json()) as MovieResult[];
		return data;
	} catch (e) {
		throw e;
	}
}
