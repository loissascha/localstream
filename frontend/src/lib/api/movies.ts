import { API_URL } from '$lib/consts';
import type { MovieListResponse } from '$lib/types/export_types';

export async function listMovies(bearerToken: string): Promise<MovieListResponse> {
	try {
		const res = await fetch(API_URL + '/api/v1/movies/list', {
			headers: {
				Authorization: 'Bearer ' + bearerToken
			}
		});
		if (!res.ok) {
			const m = res.text();
			throw new Error(`Failed to load movies: ${res.status} ${m}`);
		}

		const data = (await res.json()) as MovieListResponse;
		return data;
	} catch (error) {
		throw error;
	}
}
